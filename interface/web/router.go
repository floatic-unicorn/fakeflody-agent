package web

import (
	"fakeflody-agent/interface/agent"
	"fakeflody-agent/internal/robot/message"
	"github.com/gofiber/fiber/v3"
)

func Route(
	api *fiber.App,
	client agent.FlodyClient,
) {

	v1 := api.Group("/fakeflody/v1")

	v1.Get("/robots", func(c fiber.Ctx) error {
		return c.JSON(client.GetRobots())
	})
	v1.Get("/robots/:robotId", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		vrobot := client.GetRobotById(robotId)
		if vrobot == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Not found Robot"))
		}

		return c.JSON(vrobot)
	})

	v1.Post("/robots/boot", func(c fiber.Ctx) error {
		req := new(message.BootRobotRequest)
		if err := c.Bind().Body(req); err != nil {
			return err
		}

		vrobot := client.GetRobotById(req.RobotId)
		if vrobot != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Already Robot booted"))
		}
		client.AddRobot(req.RobotId, req.Memo)

		return c.JSON(client.GetRobots())
	})

	v1.Delete("/robots/shutdown", func(c fiber.Ctx) error {
		req := new(message.ShutDownRobotRequest)
		if err := c.Bind().Body(req); err != nil {
			return err
		}

		vrobot := client.GetRobotById(req.RobotId)
		if vrobot == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Not found Robot"))
		}
		client.RemoveRobot(req.RobotId)

		return c.JSON(client.GetRobots())
	})
}
