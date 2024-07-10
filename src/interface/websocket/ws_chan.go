package websocket

import (
	"fakeflody-agent/src/core"
	"fakeflody-agent/src/message"
	"fakeflody-agent/src/utils"
	"strconv"
)

type RobotWebSocketChannel struct {
	Channel chan message.GetRobotResult
}

func NewWebSocketQueue() core.RobotEventOutput {
	return &RobotWebSocketChannel{
		Channel: make(chan message.GetRobotResult),
	}
}

func (svc *RobotWebSocketChannel) Notify(msg *core.VRobotInfo) {
	robotInfo := msg
	_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(msg.RobotId))
	svc.Channel <- message.GetRobotResult{
		RobotId:        robotInfo.RobotId,
		RobotName:      robotInfo.RobotName,
		Memo:           robotInfo.Memo,
		State:          robotInfo.State,
		Estop:          robotInfo.EmergencyStop.Estop,
		Problems:       robotInfo.EmergencyStop.Problems,
		Solutions:      robotInfo.EmergencyStop.Solutions,
		SessionStartAt: robotInfo.SessionStartedAt.UnixMilli(),
		SessionEndAt:   endTime.UnixMilli(),
	}
}

func (svc *RobotWebSocketChannel) GetChannel(robotId int) <-chan message.GetRobotResult {
	return svc.Channel
}

func (svc *RobotWebSocketChannel) Close() {
	close(svc.Channel)
}
