package web

import (
	"context"
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/interface/web/router"
	"fakeflody-agent/src/internal/robot"
	"fakeflody-agent/src/logger"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type FiberApiServer struct {
	server          *fiber.App
	fakeFlodyClient agent.FlodyClient
}

func NewFiberApiServer(client agent.FlodyClient) *FiberApiServer {
	server := initializeServer()
	return &FiberApiServer{
		server:          server,
		fakeFlodyClient: client,
	}
}

func Server(
	lc fx.Lifecycle,
	conf *config.FakeFlodyConfig,
	client agent.FlodyClient,
	fakeRobotSvc robot.IFakeRobotService,
) {

	if conf.InterfaceConfig.Web == false {
		return
	}

	logger.Info("Starting Web Server")

	api := NewFiberApiServer(client)

	router.HealthCheckRoute(api.server)
	router.Route(api.server, fakeRobotSvc)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := api.Listen(":" + fmt.Sprint(conf.InterfaceConfig.WebConfig.Port))
				if err != nil {
					logger.Error(err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return api.Shutdown()
		},
	})
}
