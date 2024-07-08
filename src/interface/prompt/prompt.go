package prompt

import (
	"context"
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/logger"
	"fakeflody-agent/src/utils"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ReportCommandKeyword string

func (r ReportCommandKeyword) String() string {
	return string(r)
}

const (
	ActionEstop    ReportCommandKeyword = "estop"
	ActionALLEstop ReportCommandKeyword = "all-estop"
	ActionALLClear ReportCommandKeyword = "all-clear"
	ActionClear    ReportCommandKeyword = "clear"
)

type Prompt struct {
	config    *config.FakeFlodyConfig
	client    agent.FlodyClient
	lifecycle fx.Lifecycle
}

func NewPrompt(
	conf *config.FakeFlodyConfig,
	client agent.FlodyClient,
	lifecycle fx.Lifecycle,
) *Prompt {

	if conf.InterfaceConfig.Cli == false {
		return nil
	}

	logger.Info("Starting CLI Prompt")

	prompt := &Prompt{
		config: conf,
		client: client,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go prompt.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return prompt.Stop()
		},
	})

	return prompt
}

func (p *Prompt) Run() error {

	//reader := bufio.NewReader(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	time.Sleep(1 * time.Second)

	go func() {
		for {
			chioose := []choose.Choice{
				{Text: "info", Note: "전체 로봇 정보"},
				{Text: ActionEstop.String(), Note: "e-stop 버튼 누름"},
				{Text: ActionClear.String(), Note: "e-stop 버튼 해제"},
			}

			if len(p.client.GetRobots()) > 0 {
				chioose = append(chioose, choose.Choice{Text: ActionALLEstop.String(), Note: "모든 로봇 e-stop"})
				chioose = append(chioose, choose.Choice{Text: ActionALLClear.String(), Note: "모든 로봇 e-stop 해제"})
			}
			chioose = append(chioose, choose.Choice{Text: "exit", Note: "종료"})

			command, err := prompt.New().Ask("🕹️명령어를 입력하세요").
				AdvancedChoose(chioose)

			if err != nil {
				logger.Errorf("명령을 읽는데 오류가 발생했습니다: %s", err)
				continue
			}

			command = strings.TrimSpace(command)

			switch command {
			case "exit":
				os.Exit(0)
			case ActionEstop.String():
				p.chooseAndActionRobot(func(robotId int) {
					p.client.GetRobots().GetVRobotById(robotId).Estop()
					logger.Infof("🕹️%v e-stop을 눌렀습니다.", robotId)
				})

			case ActionClear.String():
				p.chooseAndActionRobot(func(robotId int) {
					p.client.GetRobots().GetVRobotById(robotId).ClearEstop()
					logger.Infof("🕹️비상정지가 해제되었습니다. \n")
				})
			case ActionALLEstop.String():
				robots := p.client.GetRobots()
				for _, robot := range robots {
					robot.Estop()
				}
				logger.Infof("🕹️모든 로봇 e-stop을 눌렀습니다.")
			case ActionALLClear.String():
				robots := p.client.GetRobots()
				for _, robot := range robots {
					robot.ClearEstop()
				}
				logger.Infof("🕹️모든 로봇️의 비상정지가 해제했습니다.")

			case "info":
				robots := p.client.GetRobots()
				for _, robot := range robots {
					robotInfo := robot.GetInfo()
					logger.Infof("[로봇 정보 - ID: %d]", robotInfo.RobotId)
					logger.Infof("🕹️로봇 이름: %s", robotInfo.RobotName)
					logger.Infof("🕹️로봇 상태: %s", robotInfo.State)

					var estopText string
					if robotInfo.EmergencyStop.Estop {
						estopText = "🔴(E-Stop 해제필요)"
					} else {
						estopText = "🔵"
					}
					_, endTime, _ := utils.Cache.GetWithExpiration(strconv.Itoa(robotInfo.RobotId))
					logger.Infof("🕹️로봇 가용상태: %s", estopText)
					logger.Infof("🕹️로봇 메모: %s", robotInfo.Memo)
					logger.Infof("🕹️남은 세션: %s", endTime)
					logger.Info("========================================")
				}
			default:
				logger.Infof("알 수 없는 명령입니다: %s", command)
			}
		}
	}()

	<-sigChan
	return nil
}

func (p *Prompt) chooseAndActionRobot(action func(robotId int)) {
	robots := p.client.GetRobots()
	ids := robots.GetRobotIds()
	ids = append(ids, "뒤로")
	robotIdString, _ := prompt.New().Ask("로봇 선택:").Choose(
		ids,
		choose.WithTheme(choose.ThemeLine),
		choose.WithKeyMap(choose.HorizontalKeyMap),
	)

	robotId, error := strconv.Atoi(robotIdString)
	if error == nil {
		action(robotId)
	}

}

func (p *Prompt) Stop() error {
	return nil
}
