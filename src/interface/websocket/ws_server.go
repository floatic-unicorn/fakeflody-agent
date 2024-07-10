package websocket

import (
	"fakeflody-agent/src/core"
	"fakeflody-agent/src/interface/web"
	"fakeflody-agent/src/logger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type WebSocketServer struct {
	Server           *web.FiberApiServer
	RobotEventOutput core.RobotEventOutput
}

func NewWebSocketServer(
	server *web.FiberApiServer,
	robotEventOutput core.RobotEventOutput,
) {
	server.Server.Use("/v1/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	server.Server.Get("/v1/ws/robots/:robotId", websocket.New(func(c *websocket.Conn) {
		robotId := getParamsToInt(c, "robotId")
		for {
			msg := <-robotEventOutput.GetChannel(robotId)
			logger.Infof("메시지 전달 %v, %+v", robotId, msg)
			if msg.RobotId == robotId {
				if err := c.WriteJSON(msg); err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}))
}

func getParamsToInt(c *websocket.Conn, name string) int {
	param := c.Params(name)
	numberValue, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return numberValue
}
