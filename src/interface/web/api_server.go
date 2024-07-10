package web

import (
	"context"
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/interface/web/router"
	"fakeflody-agent/src/internal/robot"
	"fakeflody-agent/src/logger"
	"fakeflody-agent/src/thirdparty"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type FiberApiServer struct {
	Server          *fiber.App
	FakeFlodyClient agent.FlodyClient
}

func NewFiberApiServer(client agent.FlodyClient) *FiberApiServer {
	server := initializeServer()
	return &FiberApiServer{
		Server:          server,
		FakeFlodyClient: client,
	}
}

func Server(
	lc fx.Lifecycle,
	server *FiberApiServer,
	conf *config.FakeFlodyConfig,
	fakeRobotSvc robot.IFakeRobotService,
	robotInfoSvc thirdparty.RobotInfoService,
) {

	if conf.InterfaceConfig.Web == false {
		return
	}

	logger.Info("Starting Web Server")

	api := server

	router.HealthCheckRoute(api.Server)
	router.Route(api.Server, fakeRobotSvc, robotInfoSvc)

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
