package core

import (
	"encoding/json"
	config "fakeflody-agent/src/config"
	"fakeflody-agent/src/logger"
	"fakeflody-agent/src/utils/hashids"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type IReportedEvent interface {
	AddRobot(robot *VRobotInfo)
	EStop(robotId int)
	UnPauseSuccess(robotId int)
	UnPauseFail(robotId int)
	SendReport(event *ReportedEvent)
	Close()
}

type ReportedProducer struct {
	config        *config.FakeFlodyConfig
	producer      *kafka.Producer
	reportedTopic string

	robot *VRobotInfo
}

func NewReportedProducer(cnf *config.FakeFlodyConfig, topic string) IReportedEvent {
	producer, err := config.NewProducer(cnf)
	if err != nil {
		logger.Errorf(err.Error())
		return nil
	}
	return &ReportedProducer{
		config:        cnf,
		producer:      producer,
		reportedTopic: topic,
	}
}

func (c *ReportedProducer) AddRobot(robot *VRobotInfo) {
	c.robot = robot
}

func (c *ReportedProducer) EStop(robotId int) {

	receiveState := EMERGENCY_STOPPED

	eventMessage := &ReportedEvent{
		Header: Header{
			TimeStamp: time.Now().UnixMilli(),
			Type:      REPORT.String(),
			RobotId:   hashids.ToUid(robotId),
		},
		Payload: map[string]interface{}{
			"state":     receiveState,
			"problems":  []string{"로봇에 문제가 발생했습니다."},
			"solutions": []string{"문제를 해결해주세요."},
		},
	}
	c.SendReport(eventMessage)
	c.robot.UpdateState(receiveState.String(), eventMessage.Header.CommandId)
}

func (c *ReportedProducer) UnPauseSuccess(robotId int) {

	receiveState := RECOVERED_FROM_EMERGENCY_STOP

	eventMessage := &ReportedEvent{
		Header: Header{
			CommandId: c.robot.LatestCommandId,
			TimeStamp: time.Now().UnixMilli(),
			Type:      REPORT.String(),
			RobotId:   hashids.ToUid(robotId),
		},
		Payload: map[string]interface{}{
			"state": receiveState,
		},
	}
	c.SendReport(eventMessage)
	c.robot.UpdateState(receiveState.String(), eventMessage.Header.CommandId)

	logger.WInfof("🤖[%v] 복구가 완료되었습니다", robotId)
}

func (c *ReportedProducer) UnPauseFail(robotId int) {

	receiveState := FAILED_TO_UNPAUSE

	eventMessage := &ReportedEvent{
		Header: Header{
			CommandId: c.robot.LatestCommandId,
			TimeStamp: time.Now().UnixMilli(),
			Type:      REPORT.String(),
			RobotId:   hashids.ToUid(robotId),
		},
		Payload: map[string]interface{}{
			"state":  receiveState,
			"reason": "COLLISION IS NOT FREE",
		},
	}
	c.SendReport(eventMessage)
	c.robot.UpdateState(receiveState.String(), eventMessage.Header.CommandId)

	logger.WWarnf("🤖[%v] estop 해제가 필요합니다", robotId)
}

func (c *ReportedProducer) SendReport(event *ReportedEvent) {
	value, _ := json.Marshal(&event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &c.reportedTopic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	prettyJSON, _ := json.MarshalIndent(event, "", "    ")
	logger.WDebugf("🔖 [%s] 메시지 전송:\n %s", c.reportedTopic, prettyJSON)
	err := c.producer.Produce(msg, nil) // delivery channel
	if err != nil {
		logger.WError(err.Error())
	}
}

func (c *ReportedProducer) Close() {
	c.producer.Close()
}
