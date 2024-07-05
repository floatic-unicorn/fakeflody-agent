package vrobot

import (
	config "fakeflody-agent/src/config"
	"fakeflody-agent/src/internal/robot/vrobot_msg"
	"fakeflody-agent/src/logger"
	utils2 "fakeflody-agent/src/utils"
	"fakeflody-agent/src/utils/kafka"
	"strconv"
	"time"
)

type VRobot interface {
	Estop()
	ClearEstop()
	Recover()
	IsReady() bool

	GetRobotId() int
	GetInfo() vrobot_msg.GetRobotResult

	UpdateState(state string, LatestCommandId *string)
	SetCommandId(LatestCommandId *string)

	Boot()
	Shutdown()
}

type VRobotInfo struct {
	RobotId          int                     `json:"robotId"`
	RobotName        string                  `json:"robotName"`
	EmergencyStop    RobotEmergencyStopState `json:"emergencyStop"`
	State            string                  `json:"state"`
	LatestCommandId  *string                 `json:"latestCommandId"`
	Memo             string                  `json:"memo"`
	SessionStartedAt time.Time               `json:"sessionStartedAt"`

	// events
	desiredEvent   DesiredEvent
	operationEvent OperationEvent
	reportedEvent  ReportedEvent
}

// 에러 상태
type RobotEmergencyStopState struct {
	Estop     bool     `json:"estop"`
	Solutions []string `json:"solutions"`
	Problems  []string `json:"problems"`
}

func NewRobot(robotId int, robotName string, memo string, cnf *config.FakeFlodyConfig) VRobot {
	topics := config.NewTopicConfig(cnf.Env, robotId)
	adminClinet, err := config.NewAdmin(cnf)
	if err == nil {
		if err = kafka.CreateTopicIfNotExists(adminClinet, []string{
			topics.DesiredTopic,
			topics.OperationTopic,
			topics.ReportedTopic,
			topics.RobotInfoTopic,
		}); err != nil {
			logger.Errorf(err.Error())
		}
	}

	reportedProducer := NewReportedProducer(cnf, topics.ReportedTopic)
	operationConsumer := NewOperationConsumer(cnf, reportedProducer, topics.OperationTopic)
	desiredConsumer := NewDesiredConsumer(cnf, reportedProducer, topics.DesiredTopic)

	robot := &VRobotInfo{
		RobotId:        robotId,
		RobotName:      robotName,
		State:          "BOOT",
		Memo:           memo,
		desiredEvent:   desiredConsumer,
		operationEvent: operationConsumer,
		reportedEvent:  reportedProducer,
		EmergencyStop: RobotEmergencyStopState{
			Estop:     false,
			Solutions: []string{},
			Problems:  []string{},
		},
		SessionStartedAt: time.Now(),
	}

	desiredConsumer.AddRobot(robot)
	operationConsumer.AddRobot(robot)
	reportedProducer.AddRobot(robot)

	return robot
}

func (r *VRobotInfo) Boot() {
	go r.desiredEvent.Subscribe()
	go r.operationEvent.Subscribe()
}

func (r *VRobotInfo) Shutdown() {
	r.desiredEvent.Close()
	r.operationEvent.Close()
}

func (r *VRobotInfo) UpdateState(state string, LatestCommandId *string) {
	r.State = state
	r.LatestCommandId = LatestCommandId
}

func (r *VRobotInfo) SetCommandId(LatestCommandId *string) {
	r.LatestCommandId = LatestCommandId
}

func (r *VRobotInfo) Estop() {
	r.EmergencyStop.Estop = true
	r.EmergencyStop.Problems = []string{"로봇에 문제가 발생했습니다."}
	r.EmergencyStop.Solutions = []string{"문제를 해결해주세요."}
	r.reportedEvent.EStop(r.RobotId)
}

func (r *VRobotInfo) ClearEstop() {
	r.EmergencyStop.Estop = false
	r.EmergencyStop.Problems = []string{}
	r.EmergencyStop.Solutions = []string{}
}

func (r *VRobotInfo) Recover() {
	r.EmergencyStop.Problems = []string{}
	r.EmergencyStop.Solutions = []string{}

	if r.IsReady() {
		r.reportedEvent.UnPauseSuccess(r.RobotId)
	} else {
		r.reportedEvent.UnPauseFail(r.RobotId)
	}
}

func (r *VRobotInfo) IsReady() bool {
	return r.EmergencyStop.Estop == false
}

func (r *VRobotInfo) GetInfo() vrobot_msg.GetRobotResult {
	_, endTime, _ := utils2.Cache.GetWithExpiration(strconv.Itoa(r.RobotId))
	return vrobot_msg.GetRobotResult{
		RobotId:        r.RobotId,
		RobotName:      r.RobotName,
		Memo:           r.Memo,
		State:          r.State,
		Estop:          r.EmergencyStop.Estop,
		Problems:       r.EmergencyStop.Problems,
		Solutions:      r.EmergencyStop.Solutions,
		SessionStartAt: utils2.TimeToStringDateTime(r.SessionStartedAt),
		SessionEndAt:   utils2.TimeToStringDateTime(endTime),
	}
}

func (r *VRobotInfo) GetRobotId() int {
	return r.RobotId
}

func (r *VRobotInfo) Refresh() {
	r.SessionStartedAt = time.Now()
}