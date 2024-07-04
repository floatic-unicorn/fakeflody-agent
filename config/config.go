package config

import (
	"fakeflody-agent/utils/hashids"
	"fmt"
)

type FakeFlodyConfig struct {
	RobotIds        []int  `json:"robot_id"`
	ResponseTime    int    `json:"response_time"`
	Env             string `json:"env"`
	InterfaceConfig InterfaceConfig
	BelugaConfig    BelugaConfig
}

func NewFakeFlodyConfig(env string, robotIds []int, responseTime int, interfaces InterfaceConfig) *FakeFlodyConfig {
	return &FakeFlodyConfig{
		RobotIds:        robotIds,
		ResponseTime:    responseTime,
		Env:             env,
		InterfaceConfig: interfaces,
		BelugaConfig:    NewBelugaConfig(),
	}
}

type TopicConfig struct {
	DesiredTopic   string
	ReportedTopic  string
	OperationTopic string
	RobotInfoTopic string
}

func NewTopicConfig(env string, robotId int) *TopicConfig {
	return &TopicConfig{
		DesiredTopic:   fmt.Sprintf("%s.fleet.%s.desired.json", env, hashids.ToUid(robotId)),
		ReportedTopic:  fmt.Sprintf("%s.fleet.%s.reported.json", env, hashids.ToUid(robotId)),
		OperationTopic: fmt.Sprintf("%s.fleet.%s.robot-operation.json", env, hashids.ToUid(robotId)),
		RobotInfoTopic: fmt.Sprintf("%s.fleet.%s.robot-info.json", env, hashids.ToUid(robotId)),
	}
}

type InterfaceConfig struct {
	Web       bool
	WebConfig WebConfig
	Cli       bool
}
