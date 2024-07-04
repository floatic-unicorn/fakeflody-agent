package web

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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
	return s.server.Listen(addr)
}

func (s FiberApiServer) Shutdown() error {
	return s.server.Shutdown()
}
