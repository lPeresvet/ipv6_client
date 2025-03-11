package cli

import (
	"github.com/spf13/cobra"
	"implementation/internal/domain/connections"
	"log"
)

type StatusProvider interface {
	TunnelStatus() connections.ConnectionStatus
}

type StatusCmd struct {
	cmd *cobra.Command
}

func NewStatusCmd(baseCmd *cobra.Command, provider StatusProvider) *StatusCmd {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show status of ipv6 client.",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("Connection status: %v", provider.TunnelStatus())

			return nil
		},
	}
	baseCmd.AddCommand(statusCmd)

	return &StatusCmd{
		cmd: statusCmd,
	}
}
