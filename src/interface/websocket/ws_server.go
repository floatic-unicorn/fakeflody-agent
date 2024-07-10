package websocket

import (
	"fakeflody-agent/src/core"
	"fakeflody-agent/src/interface/web"
	"fakeflody-agent/src/logger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"sync"
)

// Global map to store active WebSocket connections
var (
	connections = make(map[int]map[*websocket.Conn]bool)
	mu          sync.Mutex
)

type WebSocketServer struct {
	Server           *web.FiberApiServer
	RobotEventOutput core.RobotEventOutput
}

func NewWebSocketServer(
	server *web.FiberApiServer,
	robotEventOutput core.RobotEventOutput,
) {
	cfg := websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if err := recover(); err != nil {
				conn.WriteJSON(fiber.Map{"customError": "error occurred", "error": err})
			}
		},
	}

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

		mu.Lock()
		if connections[robotId] == nil {
			connections[robotId] = make(map[*websocket.Conn]bool)
		}
		connections[robotId][c] = true
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(connections[robotId], c)
			mu.Unlock()
			c.Close()
		}()

		for {
			msg := <-robotEventOutput.GetChannel()
			logger.Infof("메시지 전달 %v, %+v", msg.RobotId, msg)

			mu.Lock()
			conns := connections[msg.RobotId]
			mu.Unlock()

			for conn := range conns {
				if err := conn.WriteJSON(msg); err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}, cfg))
}

func getParamsToInt(c *websocket.Conn, name string) int {
	param := c.Params(name)
	numberValue, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return numberValue
}
