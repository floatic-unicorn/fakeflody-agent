package router

import (
	"fakeflody-agent/src/internal/robot"
	"fakeflody-agent/src/message"
	"fakeflody-agent/src/thirdparty"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Route(
	api *fiber.App,
	fakeRobotSvc robot.IFakeRobotService,
	robotInfoSvc thirdparty.RobotInfoService,
) {

	v1 := api.Group("/v1")

	v1.Post("/robots/boot", func(c *fiber.Ctx) error {
		req := new(message.BootRobotRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		err := fakeRobotSvc.Boot(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fakeRobotSvc.GetRobots())
	})

	v1.Delete("/robots/:robotId/shutdown", func(c *fiber.Ctx) error {
		robotId := getParamsToInt(c, "robotId")

		err := fakeRobotSvc.Shutdown(&message.ShutDownRobotRequest{
			RobotId: robotId,
		})
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fakeRobotSvc.GetRobots())
	})

	v1.Get("/robots", func(c *fiber.Ctx) error {
		return c.JSON(fakeRobotSvc.GetRobots())
	})

	v1.Get("/warehouses/:warehouseId/robotInfos", func(c *fiber.Ctx) error {
		warehouseId := getParamsToInt(c, "warehouseId")
		return c.JSON(robotInfoSvc.GetRobotInfosByWarehouse(warehouseId))
	})

	v1.Get("/robots/:robotId", func(c *fiber.Ctx) error {
		robotId := getParamsToInt(c, "robotId")

		vrobot, err := fakeRobotSvc.GetRobotById(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(vrobot)
	})

	v1.Patch("/robots/:robotId/estop", func(c *fiber.Ctx) error {
		robotId := getParamsToInt(c, "robotId")

		err := fakeRobotSvc.Estop(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})

	v1.Patch("/robots/:robotId/estop/clear", func(c *fiber.Ctx) error {
		robotId := getParamsToInt(c, "robotId")

		err := fakeRobotSvc.ClearEstop(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})

	v1.Patch("/robots/:robotId/refresh", func(c *fiber.Ctx) error {
		robotId := getParamsToInt(c, "robotId")

		err := fakeRobotSvc.RefreshSession(robotId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		return c.JSON(fiber.Map{
			"msg": "success",
		})
	})
}

func getParamsToInt(c *fiber.Ctx, name string) int {
	param := c.Params(name)
	numberValue, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return numberValue
}
