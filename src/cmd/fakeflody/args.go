package main

import "github.com/urfave/cli/v2"

var Flags = []cli.Flag{
	EnvFlag,
	RobotIdFlag,
	IntervalFlag,
	LogLevelFlag,
	WebFlag,
	WebPortFlag,
	CliFlag,
}

var EnvFlag = &cli.StringFlag{
	Name:        "env",
	Aliases:     []string{"e"},
	DefaultText: "dev",
	Value:       "dev",
	Usage:       "환경 설정",
}

var LogLevelFlag = &cli.StringFlag{
	Name:        "logLevel",
	Aliases:     []string{"log"},
	Usage:       "Log Level 설정",
	DefaultText: "info",
	Value:       "info",
}

var RobotIdFlag = &cli.IntSliceFlag{
	Name:    "robotId",
	Aliases: []string{"rid"},
	Usage:   "robotId 설정",
}

var IntervalFlag = &cli.IntFlag{
	Name:        "interval",
	Aliases:     []string{"t"},
	DefaultText: "2",
	Value:       2,
	Usage:       "응답 시간(second)을 설정합니다.",
}

var WebFlag = &cli.BoolFlag{
	Name:        "web",
	DefaultText: "false",
	Value:       false,
	Usage:       "web server를 실행합니다.",
}

var WebPortFlag = &cli.IntFlag{
	Name:        "webport",
	DefaultText: "8080",
	Value:       8080,
	Usage:       "web server port를 설정합니다.",
}

var CliFlag = &cli.BoolFlag{
	Name:        "cli",
	DefaultText: "false",
	Value:       false,
	Usage:       "CLI를 실행합니다.",
}
