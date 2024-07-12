package core

import (
	"errors"
	config "fakeflody-agent/src/config"
	"fakeflody-agent/src/logger"
	kafkautil "fakeflody-agent/src/utils/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type IDesiredEvent interface {
	AddRobot(robot *VRobotInfo)
	Subscribe()
	Close() error
}

type DesiredConsumer struct {
	config           *config.FakeFlodyConfig
	desiredConsumer  *kafka.Consumer
	reportedProducer IReportedEvent
	topic            string

	robot *VRobotInfo
}

func NewDesiredConsumer(cnf *config.FakeFlodyConfig, reportedService IReportedEvent, topic string) IDesiredEvent {
	consumer, err := config.NewConsumer(cnf)
	if err != nil {
		logger.WErrorf(err.Error())
	}
	return &DesiredConsumer{
		config:           cnf,
		desiredConsumer:  consumer,
		reportedProducer: reportedService,
		topic:            topic,
	}
}

func (c *DesiredConsumer) AddRobot(robot *VRobotInfo) {
	c.robot = robot
}

func (c *DesiredConsumer) Subscribe() {
	logger.Infof("[%s] 해당 토픽을 구독합니다.", c.topic)
	kafkautil.Subscribe[DesiredEvent](c.topic, c.desiredConsumer, func(msg *DesiredEvent) {
		state, ok := msg.Payload["state"]
		if !ok {
			return
		}
		robotState := RobotCommand(state.(string))
		nextStates := NextReports(robotState)

		for i, state := range nextStates {

			validRobotStatusErr := c.validRobotStatus(i, state)
			if validRobotStatusErr != nil {
				return
			}

			msg.Header.TimeStamp = time.Now().Unix()
			msg.Header.Type = RESPONSE.String()
			msg.Payload["state"] = state.String()

			c.reportedProducer.SendReport(&ReportedEvent{
				Header:  msg.Header,
				Payload: msg.Payload,
			})
			c.robot.UpdateState(state.String(), msg.Header.CommandId)
			go c.robot.robotEventOutput.Notify(c.robot)
		}
	})
}

func (c *DesiredConsumer) Close() error {
	return c.desiredConsumer.Close()
}

func (c *DesiredConsumer) validRobotStatus(seq int, state RobotReportState) error {
	if seq == 0 { // start
		time.Sleep(1 * time.Second)

		if !c.robot.IsReady() {
			logger.WWarnf("🤖[%v] 로봇의 estop 해제가 필요합니다 - %v", c.robot.RobotId, state.String())
			return errors.New("🤖 로봇의 estop 해제가 필요합니다")
		}
		logger.WInfof("🤖[%v] 로봇이 출발합니다 - %v", c.robot.RobotId, state.String())
	} else { // end
		time.Sleep(time.Duration(c.robot.Interval) * time.Second)

		if !c.robot.IsReady() {
			logger.WWarnf("🤖[%v] 로봇의 estop 해제가 필요합니다 - %v", c.robot.RobotId, state.String())
			return errors.New("🤖 로봇의 estop 해제가 필요합니다")
		}
		logger.WInfof("🤖[%v] 로봇이 도착했습니다 - %v", c.robot.RobotId, state.String())
	}

	return nil
}
