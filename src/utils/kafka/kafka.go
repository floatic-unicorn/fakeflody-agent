package kafka

import (
	"encoding/json"
	"fakeflody-agent/src/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"os/signal"
	"syscall"
)

func Subscribe[T any](topic string, consumer *kafka.Consumer, process func(message *T)) {
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		logger.Fatalf(err.Error())
		return
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			logger.Infof("Caught signal %v: [%s] kafka consumer close", sig, topic)
			consumer.Close()
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				var desiredEvent T
				err := json.Unmarshal(e.Value, &desiredEvent)
				if err != nil {
					logger.Errorf(err.Error())
				}
				process(&desiredEvent)

			case kafka.Error:
				logger.Infof("Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				logger.Debugf("Ignored %v", e)
			}
		}
	}
}
