package main

import (
	"fakeflody-agent/config"
	"fakeflody-agent/interface/agent"
	"fakeflody-agent/interface/prompt"
	"fakeflody-agent/interface/web"
	"fakeflody-agent/thirdparty"
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
