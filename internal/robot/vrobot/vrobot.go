package vrobot

import (
	"fakeflody-agent/config"
	"fakeflody-agent/internal/robot/message"
	"fakeflody-agent/logger"
	"fakeflody-agent/utils"
	"fakeflody-agent/utils/kafka"
	"strconv"
	"time"
)

type VRobot interface {
	Estop()
	ClearEstop()
	Recover()
	IsReady() bool

	GetRobotId() int
	GetInfo() message.GetRobotResult

	UpdateState(state string, LatestCommandId *string)
	SetCommandId(LatestCommandId *string)

	Boot()
	Shutdown()
}

type VRobotInfo struct {
	RobotId          int                     `json:"robotId"`
	RobotName        string                  `json:"robotName"`
	EmrgencyStop     RobotEmargencyStopState `json:"emrgencyStop"`
	State            string                  `json:"state"`
	LatestCommandId  *string                 `json:"latestCommandId;omitempty"`
	Memo             string                  `json:"memo"`
	SessionStartedAt time.Time               `json:"sessionStartedAt"`

	// events
	desiredEvent   DesiredEvent
	operationEvent OperationEvent
	reportedEvent  ReportedEvent
}

// 에러 상태
type RobotEmargencyStopState struct {
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
		EmrgencyStop: RobotEmargencyStopState{
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
	r.EmrgencyStop.Estop = true
	r.EmrgencyStop.Problems = []string{"로봇에 문제가 발생했습니다."}
	r.EmrgencyStop.Solutions = []string{"문제를 해결해주세요."}
	r.reportedEvent.EStop(r.RobotId)
}

func (r *VRobotInfo) ClearEstop() {
	r.EmrgencyStop.Estop = false
	r.EmrgencyStop.Problems = []string{}
	r.EmrgencyStop.Solutions = []string{}
}

func (r *VRobotInfo) Recover() {
	r.EmrgencyStop.Problems = []string{}
	r.EmrgencyStop.Solutions = []string{}

	if r.IsReady() {
		r.reportedEvent.UnPauseSuccess(r.RobotId)
	} else {
		r.reportedEvent.UnPauseFail(r.RobotId)
	}
}

func (r *VRobotInfo) IsReady() bool {
	return r.EmrgencyStop.Estop == false
}

func (r *VRobotInfo) GetInfo() message.GetRobotResult {
	_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(r.RobotId))
	return message.GetRobotResult{
		RobotId:        r.RobotId,
		RobotName:      r.RobotName,
		Memo:           r.Memo,
		State:          r.State,
		Estop:          r.EmrgencyStop.Estop,
		Problems:       r.EmrgencyStop.Problems,
		Solutions:      r.EmrgencyStop.Solutions,
		SessionStartAt: utils.TimeToStringDateTime(r.SessionStartedAt),
		SessionEndAt:   utils.TimeToStringDateTime(endTime),
	}
}

func (r *VRobotInfo) GetRobotId() int {
	return r.RobotId
}

func (r *VRobotInfo) Refresh() {
	r.SessionStartedAt = time.Now()
}
