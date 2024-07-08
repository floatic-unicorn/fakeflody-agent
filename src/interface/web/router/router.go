package router

import (
	"fakeflody-agent/src/internal/robot"
	"github.com/gofiber/fiber/v3"
)

func Route(
	api *fiber.App,
	fakeRobotSvc robot.IFakeRobotService,
) {

	v1 := api.Group("/v1")

	v1.Post("/robots/boot", func(c fiber.Ctx) error {
		req := new(robot.BootRobotRequest)
		if err := c.Bind().Body(req); err != nil {
			return err
		}

		err := fakeRobotSvc.Boot(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fakeRobotSvc.GetRobots())
	})

	v1.Delete("/robots/:robotId/shutdown", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		err := fakeRobotSvc.Shutdown(&robot.ShutDownRobotRequest{
			RobotId: robotId,
		})
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fakeRobotSvc.GetRobots())
	})

	v1.Get("/robots", func(c fiber.Ctx) error {

		return c.JSON(fakeRobotSvc.GetRobots())
	})
	v1.Get("/robots/:robotId", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		vrobot, err := fakeRobotSvc.GetRobotById(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(vrobot)
	})

	v1.Patch("/robots/:robotId/estop", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		err := fakeRobotSvc.Estop(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})

	v1.Patch("/robots/:robotId/estop/clear", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		err := fakeRobotSvc.ClearEstop(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})

	v1.Patch("/robots/:robotId/refresh", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		err := fakeRobotSvc.RefreshSession(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})
}
