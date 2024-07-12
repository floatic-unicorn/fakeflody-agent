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

func (svc *RobotWebSocketChannel) Notify(vrobot *core.VRobotInfo) {
	_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(vrobot.RobotId))
	svc.Channel <- message.GetRobotResult{
		RobotId:        vrobot.RobotId,
		RobotName:      vrobot.RobotName,
		Memo:           vrobot.Memo,
		State:          vrobot.State,
		Estop:          vrobot.EmergencyStop.Estop,
		Problems:       vrobot.EmergencyStop.Problems,
		Solutions:      vrobot.EmergencyStop.Solutions,
		Interval:       vrobot.Interval,
		SessionStartAt: vrobot.SessionStartedAt.UnixMilli(),
		SessionEndAt:   endTime.UnixMilli(),
	}
}

func (svc *RobotWebSocketChannel) GetChannel() <-chan message.GetRobotResult {
	return svc.Channel
}

func (svc *RobotWebSocketChannel) Close() {
	close(svc.Channel)
}
