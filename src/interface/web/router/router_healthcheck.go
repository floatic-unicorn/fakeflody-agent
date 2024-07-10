package router

import (
	_ "fakeflody-agent/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func HealthCheckRoute(
	router fiber.Router,
) {
	router.Get(healthcheck.DefaultLivenessEndpoint, func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		return nil
	})
	router.Get(healthcheck.DefaultReadinessEndpoint, func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		return nil
	})
}
