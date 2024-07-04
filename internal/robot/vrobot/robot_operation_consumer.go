package vrobot

import (
	"fakeflody-agent/config"
	"fakeflody-agent/internal/robot/message"
	"fakeflody-agent/logger"
	kafkautil "fakeflody-agent/utils/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type OperationEvent interface {
	AddRobot(robot *VRobotInfo)
	Subscribe()
	Close() error
}

type OperationConsumer struct {
	config          *config.FakeFlodyConfig
	desiredConsumer *kafka.Consumer
	reportedService *ReportedProducer
	topic           string

	robot *VRobotInfo
}

func NewOperationConsumer(cnf *config.FakeFlodyConfig, reportedService *ReportedProducer, topic string) *OperationConsumer {
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
	kafkautil.Subscribe(c.topic, c.desiredConsumer, func(msg *message.DesiredEvent) {

		state, ok := msg.Payload["state"]
		if !ok {
			return
		}
		robotState := message.RobotOperation(state.(string))

		if robotState == message.UNPAUSED {
			time.Sleep(1 * time.Second)

			logger.WInfof("🤖[%v] 명령을 처리합니다 - %v", c.robot.RobotId, state)

			msg.Header.TimeStamp = time.Now().Unix()
			msg.Header.Type = message.REPORT.String()

			c.robot.SetCommandId(msg.Header.CommandId)
			c.robot.Recover()
		}
	})
}

func (c *OperationConsumer) Close() error {
	return c.desiredConsumer.Close()
}
