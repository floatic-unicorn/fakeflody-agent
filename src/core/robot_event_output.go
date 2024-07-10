package core

import (
	"fakeflody-agent/src/message"
)

type RobotEventOutput interface {
	Notify(message *VRobotInfo)
	GetChannel(robotId int) <-chan message.GetRobotResult
}
