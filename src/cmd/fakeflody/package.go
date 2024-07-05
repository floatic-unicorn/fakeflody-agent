package main

import (
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/interface/agent"
	"fakeflody-agent/src/interface/prompt"
	"fakeflody-agent/src/interface/web"
	"fakeflody-agent/src/thirdparty"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func FakeFlody(ctx *cli.Context, conf *config.FakeFlodyConfig) *fx.App {
	return fx.New(
		fx.NopLogger,
		fx.Provide(
			append(
				providers(),
				func() *config.FakeFlodyConfig { return conf },
				func() *cli.Context { return ctx },
			)...,
		),
		fx.Invoke(
			invokers()...,
		),
	)
}

func providers() []interface{} {
	return []interface{}{
		agent.NewFakeFlodyClient,
		thirdparty.NewRobotInfoService,
	}
}

func invokers() []interface{} {
	return []interface{}{
		prompt.NewPrompt,
		agent.ClientSessionHandler,
		web.Server,
	}
}
