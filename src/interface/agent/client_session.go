package agent

import (
	"context"
	"fakeflody-agent/src/utils"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/fx"
	"strconv"
)

func ClientSessionHandler(
	client FlodyClient,
	lifecycle fx.Lifecycle,
) {
	err := gocron.Every(2).Seconds().Do(func() {
		for _, robot := range client.GetRobots() {
			_, isCached := utils.Cache.Get(strconv.Itoa(robot.GetRobotId()))
			if !isCached {
				client.RemoveRobot(robot.GetRobotId())
			}
		}
	})
	if err != nil {
		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				<-gocron.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
