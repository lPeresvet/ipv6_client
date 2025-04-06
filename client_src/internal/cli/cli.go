package cli

import (
	"context"
	"implementation/client_src/internal/domain/connections"
)

type CLI struct {
	commonArgs *CommonAgs
}

type ClientController interface {
	Connector
	StatusProvider
}

type UnixSocketListener interface {
	ListenIpUp(ctx context.Context, control chan *connections.IfaceEvent) error
}

func New(controller ClientController, filler ConfigFiller, listener UnixSocketListener) *CLI {
	baseCmd := NewCommonAgs()
	NewConnectCmd(baseCmd.cmd, controller, listener)
	NewStatusCmd(baseCmd.cmd, controller)
	NewConfigurerAgs(baseCmd.cmd, filler)

	return &CLI{
		commonArgs: baseCmd,
	}
}

func (cli *CLI) Execute() error {
	return cli.commonArgs.cmd.Execute()
}
