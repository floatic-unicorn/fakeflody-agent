package router

import (
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/internal/robot/vrobot_msg"
	"github.com/gofiber/fiber/v3"
)

func Route(
	api *fiber.App,
	client agent.FlodyClient,
) {

	v1 := api.Group("/v1")

	// client.GetRobots godoc
	// @Summary		게시판 등록 API
	// @Accept		json
	// @Produce		json
	// @Success		200		{object}	vrobot.VRobotList
	// @Router			/v1/robots [get]
	v1.Get("/robots", func(c fiber.Ctx) error {
		robots := make([]vrobot_msg.GetRobotResult, len(client.GetRobots()))
		for i, vrobot := range client.GetRobots() {
			robots[i] = vrobot.GetInfo()
		}
		return c.JSON(robots)
	})
	v1.Get("/robots/:robotId", func(c fiber.Ctx) error {
		robotId := fiber.Params[int](c, "robotId")

		vrobot := client.GetRobotById(robotId)
		if vrobot == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Not found Robot"))
		}

		return c.JSON(vrobot.GetInfo())
	})

	v1.Post("/robots/boot", func(c fiber.Ctx) error {
		req := new(vrobot_msg.BootRobotRequest)
		if err := c.Bind().Body(req); err != nil {
			return err
		}

		vrobot := client.GetRobotById(req.RobotId)
		if vrobot != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "Already Robot booted"))
		}
		client.AddRobot(req.RobotId, req.Memo)

		robots := make([]vrobot_msg.GetRobotResult, len(client.GetRobots()))
		for i, vrobot := range client.GetRobots() {
			robots[i] = vrobot.GetInfo()
		}
		return c.JSON(robots)
	})

	v1.Delete("/robots/shutdown", func(c fiber.Ctx) error {
		req := new(vrobot_msg.ShutDownRobotRequest)
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
