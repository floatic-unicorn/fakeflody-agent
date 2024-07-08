package agent

import (
	"context"
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/core"
	"fakeflody-agent/src/logger"
	"fakeflody-agent/src/thirdparty"
	"fakeflody-agent/src/utils"
	"fmt"
	"github.com/patrickmn/go-cache"
	"go.uber.org/fx"
	"strconv"
)

type FlodyClient interface {
	Run() error
	Stop() error
	AddRobot(robotId int, memo string) error
	RemoveRobot(robotId int)
	GetRobots() core.VRobotList
	GetRobotById(robotId int) core.VRobot
}

type FakeFlodyClient struct {
	cnf              *config.FakeFlodyConfig
	robots           core.VRobotList
	robotInfoService thirdparty.RobotInfoService
}

func NewFakeFlodyClient(
	cnf *config.FakeFlodyConfig,
	lifecycle fx.Lifecycle,
	robotInfoService thirdparty.RobotInfoService,
) FlodyClient {

	client := &FakeFlodyClient{
		cnf:              cnf,
		robots:           core.NewRobots(),
		robotInfoService: robotInfoService,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Run()
		},
		OnStop: func(ctx context.Context) error {
			return client.Stop()
		},
	})

	return client
}

func (c *FakeFlodyClient) Run() error {
	for _, robotId := range c.cnf.RobotIds {
		err := c.AddRobot(robotId, "초기 설정")
		if err != nil {
			logger.WWarnf("가상로봇 추가에 실패하였습니다 error: %v", err)
		}
	}
	for _, robot := range c.robots {
		go robot.Boot()
	}
	return nil
}

func (c *FakeFlodyClient) Stop() error {
	for _, robot := range c.robots {
		robot.Shutdown()
	}
	return nil
}

func (c *FakeFlodyClient) AddRobot(robotId int, memo string) error {
	robotName, err := c.getRobotByRobotName(robotId)
	if err != nil {
		return err
	}

	bootRobot := core.NewRobot(robotId, robotName, memo, c.cnf)
	utils.Cache.Set(strconv.Itoa(robotId), robotId, cache.DefaultExpiration)

	c.robots = append(c.robots, bootRobot)
	go bootRobot.Boot()

	return nil
}

func (c *FakeFlodyClient) RemoveRobot(robotId int) {
	for i, robot := range c.robots {
		if robot.GetRobotId() == robotId {
			c.robots = append(c.robots[:i], c.robots[i+1:]...)
			utils.Cache.Delete(strconv.Itoa(robotId))
			break
		}
	}
}

func (c *FakeFlodyClient) GetRobots() core.VRobotList {
	return c.robots
}

func (c *FakeFlodyClient) GetRobotById(robotId int) core.VRobot {
	return c.robots.GetVRobotById(robotId)
}

func (c *FakeFlodyClient) getRobotByRobotName(robotId int) (robotName string, error error) {
	if c.cnf.Env == "dev" {
		robotInfo := c.robotInfoService.GetRobotInfo(robotId)
		if robotInfo != nil && robotInfo.Name != "" {
			robotName = robotInfo.Name
		} else {
			error = fmt.Errorf("robot not found")
		}
		return robotName, error
	}

	return fmt.Sprintf("LOCAL-ROBOT-%v", robotId), nil
}
