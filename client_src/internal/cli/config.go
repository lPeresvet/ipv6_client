package cli

import "github.com/spf13/cobra"

type ConfigFiller interface {
	InitConfig(path string) error
}

type ConfigurerAgs struct {
	cmd *cobra.Command
}

func NewConfigurerAgs(baseCmd *cobra.Command, filler ConfigFiller) *ConfigurerAgs {
	var configPath string

	cmd := &cobra.Command{
		Use:   "configure [flags]",
		Short: "Write users and servers from yaml to config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return filler.InitConfig(configPath)
		},
	}

	cmd.Flags().StringVarP(&configPath, "path", "p", "config/config-example.yaml", "Path to config file with users and peers")

	baseCmd.AddCommand(cmd)

	return &ConfigurerAgs{cmd: cmd}
}
