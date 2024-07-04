package main

import (
	"fakeflody-agent/config"
	"fakeflody-agent/logger"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	Name    = "fakeflody"
	Usage   = "Flody Robot의 가상환경을 지원해줍니다."
	Version = "1.0.2"
)

func main() {
	app := &cli.App{
		Name:     Name,
		Usage:    Usage,
		Version:  Version,
		Flags:    Flags,
		Commands: []*cli.Command{
			//reportCommand,
		},
		Action: func(ctx *cli.Context) error {
			env := ctx.String(EnvFlag.Name)
			if env == "" {
				logger.Fatal("env 를 입력해 주세요. ( --env )")
				return nil
			}

			robotIds := ctx.IntSlice(RobotIdFlag.Name)

			FakeFlodyConfig := config.NewFakeFlodyConfig(
				env,
				robotIds,
				ctx.Int(ResponseTimeFlag.Name),
				config.InterfaceConfig{
					Web: ctx.Bool(WebFlag.Name),
					Cli: ctx.Bool(CliFlag.Name),
					WebConfig: config.WebConfig{
						Port: ctx.Int(WebPortFlag.Name),
					},
				})

			logger.Infof("FakeFlody 실행합니다")
			logger.Infof("환경: %s", env)
			logger.Infof("로봇 ID: %d", FakeFlodyConfig.RobotIds)

			logLevel := ctx.String(LogLevelFlag.Name)

			logger.NewLogger(logLevel)
			logger.NewFileLogger(logLevel, "fakeflody.log")

			FakeFlody(ctx, FakeFlodyConfig).Run()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
