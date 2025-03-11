package cli

type CLI struct {
	commonArgs *CommonAgs
}

type ClientController interface {
	Connector
	StatusProvider
}

func New(controller ClientController) *CLI {
	baseCmd := NewCommonAgs()
	NewConnectCmd(baseCmd.cmd, controller)
	NewStatusCmd(baseCmd.cmd, controller)

	return &CLI{
		commonArgs: baseCmd,
	}
}

func (cli *CLI) Execute() error {
	return cli.commonArgs.cmd.Execute()
}
