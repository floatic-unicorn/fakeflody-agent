package config

import (
	"fakeflody-agent/src/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var consumerLocalConfig = kafka.ConfigMap{
	"bootstrap.servers":  "localhost:9092",
	"group.id":           "fakeflody",
	"auto.offset.reset":  "latest",
	"enable.auto.commit": true,
	"session.timeout.ms": 6000,
}

var productLocalConfig = kafka.ConfigMap{
	"bootstrap.servers": "localhost:9092",
	"acks":              "all",
}

var consumerDevConfig = kafka.ConfigMap{
	"bootstrap.servers":  "pkc-e82om.ap-northeast-2.aws.confluent.cloud:9092",
	"group.id":           "fakeflody",
	"auto.offset.reset":  "latest",
	"enable.auto.commit": true,
	"session.timeout.ms": 6000,

	"sasl.mechanism":    "PLAIN",
	"sasl.username":     "ATKZAMUD2MADXMNN",
	"sasl.password":     "OLPfOeb1C4MhGQBzWNRuowkhhXc+cDSTOVau3IzlqfA9q98OM1yNoYU9HaRVPb+q",
	"security.protocol": "SASL_SSL",
}

var producerDevConfig = kafka.ConfigMap{
	"bootstrap.servers": "pkc-e82om.ap-northeast-2.aws.confluent.cloud:9092",
	"acks":              "all",

	"sasl.username":     "ATKZAMUD2MADXMNN",
	"sasl.password":     "OLPfOeb1C4MhGQBzWNRuowkhhXc+cDSTOVau3IzlqfA9q98OM1yNoYU9HaRVPb+q",
	"sasl.mechanism":    "PLAIN",
	"security.protocol": "SASL_SSL",
}

func NewConsumer(config *FakeFlodyConfig) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(loadConsumerConfigByEnv(config.Env))
	if err != nil {
		logger.Fatalf(err.Error())
	}
	return c, err
}

func NewProducer(config *FakeFlodyConfig) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(loadProducerConfigByEnv(config.Env))
	if err != nil {
		logger.Fatalf(err.Error())
	}
	return p, err
}

func NewAdmin(config *FakeFlodyConfig) (*kafka.AdminClient, error) {
	c, err := kafka.NewAdminClient(loadProducerConfigByEnv(config.Env))
	if err != nil {
		logger.Fatalf(err.Error())
	}
	return c, err
}

func loadConsumerConfigByEnv(env string) *kafka.ConfigMap {
	switch env {
	case "local":
		return &consumerLocalConfig
	case "dev":
		return &consumerDevConfig
	default:
		return &consumerLocalConfig
	}
}

func loadProducerConfigByEnv(env string) *kafka.ConfigMap {
	switch env {
	case "local":
		return &productLocalConfig
	case "dev":
		return &producerDevConfig
	default:
		return &productLocalConfig
	}
}
