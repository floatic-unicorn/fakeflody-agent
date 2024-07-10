package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/utils/v2"
)

func initializeServer() *fiber.App {
	server := fiber.New()
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
		//AllowCredentials: true,
	}))
	server.Use(requestid.New(requestid.Config{
		Next:      nil,
		Header:    fiber.HeaderXRequestID,
		Generator: utils.UUID,
	}))

	server.Use(recover.New())
	return server
}

func (s FiberApiServer) Listen(addr string) error {
	return s.Server.Listen(addr)
}

func (s FiberApiServer) Shutdown() error {
	return s.Server.Shutdown()
}
