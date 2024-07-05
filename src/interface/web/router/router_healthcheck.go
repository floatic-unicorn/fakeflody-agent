package router

import (
	_ "fakeflody-agent/docs"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func HealthCheckRoute(
	router fiber.Router,
) {
	router.Get(healthcheck.DefaultLivenessEndpoint, healthcheck.NewHealthChecker())
	router.Get(healthcheck.DefaultReadinessEndpoint, healthcheck.NewHealthChecker())
}
