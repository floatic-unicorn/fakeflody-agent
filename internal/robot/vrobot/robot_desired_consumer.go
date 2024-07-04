package vrobot

import (
	"fakeflody-agent/config"
	"fakeflody-agent/internal/robot/message"
	"fakeflody-agent/logger"
	kafkautil "fakeflody-agent/utils/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type DesiredEvent interface {
	AddRobot(robot *VRobotInfo)
	Subscribe()
	Close() error
}

type DesiredConsumer struct {
	config           *config.FakeFlodyConfig
	desiredConsumer  *kafka.Consumer
	reportedProducer *ReportedProducer
	topic            string

	robot *VRobotInfo
}

func NewDesiredConsumer(cnf *config.FakeFlodyConfig, reportedService *ReportedProducer, topic string) *DesiredConsumer {
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
	kafkautil.Subscribe(c.topic, c.desiredConsumer, func(msg *message.DesiredEvent) {
		state, ok := msg.Payload["state"]
		if !ok {
			return
		}
		robotState := message.RobotCommand(state.(string))
		nextStates := message.NextReports(robotState)

		for i, state := range nextStates {

			time.Sleep(time.Duration(c.config.ResponseTime) * time.Second)

			if !c.robot.IsReady() {
				logger.WWarnf("🤖[%v] 로봇의 estop 해제가 필요합니다 - %v", c.robot.RobotId, robotState.String())
				return
			}

			if i == 0 {
				logger.WInfof("🤖[%v] 로봇이 출발합니다 - %v", c.robot.RobotId, robotState.String())
			} else {
				logger.WInfof("🤖[%v] 로봇이 도착했습니다 - %v", c.robot.RobotId, robotState.String())
			}

			msg.Header.TimeStamp = time.Now().Unix()
			msg.Header.Type = message.RESPONSE.String()
			msg.Payload["state"] = state.String()

			c.reportedProducer.SendReport(&message.ReportedEvent{
				Header:  msg.Header,
				Payload: msg.Payload,
			})
			c.robot.UpdateState(state.String(), msg.Header.CommandId)
		}
	})
}

func (c *DesiredConsumer) Close() error {
	return c.desiredConsumer.Close()
}
