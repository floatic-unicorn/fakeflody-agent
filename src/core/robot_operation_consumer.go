package core

import (
	config "fakeflody-agent/src/config"
	"fakeflody-agent/src/logger"
	kafkautil "fakeflody-agent/src/utils/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type IOperationEvent interface {
	AddRobot(robot *VRobotInfo)
	Subscribe()
	Close() error
}

type OperationConsumer struct {
	config          *config.FakeFlodyConfig
	desiredConsumer *kafka.Consumer
	reportedService IReportedEvent
	topic           string

	robot *VRobotInfo
}

func NewOperationConsumer(cnf *config.FakeFlodyConfig, reportedService IReportedEvent, topic string) IOperationEvent {
	consumer, err := config.NewConsumer(cnf)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &OperationConsumer{
		config:          cnf,
		desiredConsumer: consumer,
		reportedService: reportedService,
		topic:           topic,
	}
}

func (c *OperationConsumer) AddRobot(robot *VRobotInfo) {
	c.robot = robot
}

func (c *OperationConsumer) Subscribe() {
	logger.Infof("[%s] 해당 토픽을 구독합니다.", c.topic)
	kafkautil.Subscribe[DesiredEvent](c.topic, c.desiredConsumer, func(msg *DesiredEvent) {

		state, ok := msg.Payload["state"]
		if !ok {
			return
		}
		robotState := RobotOperation(state.(string))

		if robotState == UNPAUSED {
			time.Sleep(1 * time.Second)

			logger.WInfof("🤖[%v] 명령을 처리합니다 - %v", c.robot.RobotId, state)

			msg.Header.TimeStamp = time.Now().Unix()
			msg.Header.Type = REPORT.String()

			c.robot.SetCommandId(msg.Header.CommandId)
			c.robot.Recover()
			c.robot.robotEventOutput.Notify(c.robot)
		}
	})
}

func (c *OperationConsumer) Close() error {
	return c.desiredConsumer.Close()
}
