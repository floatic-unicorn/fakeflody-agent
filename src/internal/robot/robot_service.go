package robot

import (
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/core"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/message"
	"fakeflody-agent/src/thirdparty"
	"fakeflody-agent/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"strconv"
)

type IFakeRobotService interface {
	Boot(req *message.BootRobotRequest) error
	Shutdown(req *message.ShutDownRobotRequest) error
	GetRobots() []*message.GetRobotResult
	GetRobotById(robotId int) (*message.GetRobotResult, error)
	Estop(robotId int) error
	ClearEstop(robotId int) error
	AllEstop() error
	AllClearEstop() error
	RefreshSession(robotId int) error
}

func NewFakeRobotService(
	config *config.FakeFlodyConfig,
	client agent.FlodyClient,
	robotInfoSvc thirdparty.RobotInfoService,
) IFakeRobotService {
	return &FakeRobotService{
		config:       config,
		client:       client,
		robotInfoSvc: robotInfoSvc,
	}
}

// FakeRobotService interface
type FakeRobotService struct {
	config       *config.FakeFlodyConfig
	client       agent.FlodyClient
	robotEvent   core.RobotEventOutput
	robotInfoSvc thirdparty.RobotInfoService
}

func (svc *FakeRobotService) GetRobots() []*message.GetRobotResult {
	robots := make([]*message.GetRobotResult, 0)

	for _, vrobot := range svc.client.GetRobots() {
		robot, _ := svc.GetRobotById(vrobot.GetRobotId())
		robots = append(robots, robot)
	}
	return robots
}

func (svc *FakeRobotService) GetRobotById(robotId int) (*message.GetRobotResult, error) {

	vrobot := svc.client.GetRobotById(robotId)
	if vrobot == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	robotInfo := vrobot.GetInfo()
	_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(robotId))

	return &message.GetRobotResult{
		RobotId:        vrobot.GetRobotId(),
		RobotName:      robotInfo.RobotName,
		Memo:           robotInfo.Memo,
		State:          robotInfo.State,
		Estop:          robotInfo.EmergencyStop.Estop,
		Problems:       robotInfo.EmergencyStop.Problems,
		Solutions:      robotInfo.EmergencyStop.Solutions,
		Interval:       robotInfo.Interval,
		SessionStartAt: robotInfo.SessionStartedAt.UnixMilli(),
		SessionEndAt:   endTime.UnixMilli(),
	}, nil
}

func (svc *FakeRobotService) Estop(robotId int) error {
	vrobot := svc.client.GetRobotById(robotId)
	if vrobot == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	vrobot.Estop()

	return nil
}

func (svc *FakeRobotService) Boot(req *message.BootRobotRequest) error {
	vrobot := svc.client.GetRobotById(req.RobotId)
	if vrobot != nil {
		return fiber.NewError(fiber.StatusBadRequest, "The robot is already running")
	}
	err := svc.client.AddRobot(req.RobotId, req.Memo, req.Interval)
	if err != nil {
		return err
	}

	return nil
}

func (svc *FakeRobotService) ClearEstop(robotId int) error {
	vrobot := svc.client.GetRobotById(robotId)
	if vrobot == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	vrobot.ClearEstop()

	return nil
}

func (svc *FakeRobotService) AllEstop() error {
	for _, vrobot := range svc.client.GetRobots() {
		vrobot.Estop()
	}

	return nil
}

func (svc *FakeRobotService) AllClearEstop() error {
	for _, vrobot := range svc.client.GetRobots() {
		vrobot.ClearEstop()
	}

	return nil
}

func (svc *FakeRobotService) Shutdown(req *message.ShutDownRobotRequest) error {
	vrobot := svc.client.GetRobotById(req.RobotId)
	if vrobot == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	svc.client.RemoveRobot(req.RobotId)

	return nil
}

func (svc *FakeRobotService) RefreshSession(robotId int) error {
	vrobot := svc.client.GetRobotById(robotId)
	if vrobot == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	utils.Cache.Set(strconv.Itoa(robotId), robotId, cache.DefaultExpiration)

	return nil
}
