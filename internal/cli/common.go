package cli

import (
	"github.com/spf13/cobra"
)

type CommonAgs struct {
	cmd *cobra.Command
}

func NewCommonAgs() *CommonAgs {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Provide ipv6 connection over ipv4",
	}

	return &CommonAgs{cmd: cmd}
}
