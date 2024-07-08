package robot

import (
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/patrickmn/go-cache"
	"strconv"
)

type IFakeRobotService interface {
	Boot(req *BootRobotRequest) error
	Shutdown(req *ShutDownRobotRequest) error
	GetRobots() []*GetRobotResult
	GetRobotById(robotId int) (*GetRobotResult, error)
	Estop(robotId int) error
	ClearEstop(robotId int) error
	RefreshSession(robotId int) error
}

// FakeRobotService interface
type FakeRobotService struct {
	config *config.FakeFlodyConfig
	client agent.FlodyClient
}

func NewFakeRobotService(
	config *config.FakeFlodyConfig,
	client agent.FlodyClient,
) IFakeRobotService {
	return &FakeRobotService{
		config: config,
		client: client,
	}
}

func (svc *FakeRobotService) GetRobots() []*GetRobotResult {
	robots := make([]*GetRobotResult, 0)

	for _, vrobot := range svc.client.GetRobots() {
		robot, _ := svc.GetRobotById(vrobot.GetRobotId())
		robots = append(robots, robot)
	}
	return robots
}

func (svc *FakeRobotService) GetRobotById(robotId int) (*GetRobotResult, error) {

	vrobot := svc.client.GetRobotById(robotId)
	if vrobot == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	robotInfo := vrobot.GetInfo()
	_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(robotId))

	return &GetRobotResult{
		RobotId:        vrobot.GetRobotId(),
		RobotName:      robotInfo.RobotName,
		Memo:           robotInfo.Memo,
		State:          robotInfo.State,
		Estop:          robotInfo.EmergencyStop.Estop,
		Problems:       robotInfo.EmergencyStop.Problems,
		Solutions:      robotInfo.EmergencyStop.Solutions,
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

func (svc *FakeRobotService) Boot(req *BootRobotRequest) error {
	vrobot := svc.client.GetRobotById(req.RobotId)
	if vrobot != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Not found Robot")
	}
	err := svc.client.AddRobot(req.RobotId, req.Memo)
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

func (svc *FakeRobotService) Shutdown(req *ShutDownRobotRequest) error {
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
