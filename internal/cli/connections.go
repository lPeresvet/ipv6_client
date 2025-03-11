package cli

import (
	"github.com/spf13/cobra"
)

type ConnectionsCmd struct {
	connectCmd    *cobra.Command
	disconnectCmd *cobra.Command

	username *string
}

type Connector interface {
	TunnelConnect(username string) error
	TunnelDisconnect(username string) error
}

func NewConnectCmd(baseCmd *cobra.Command, connector Connector) *ConnectionsCmd {
	var username string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to ipv6 prefix provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return connector.TunnelConnect(username)
		},
	}

	connectCmd.Flags().StringVarP(&username, "username", "u", "", "Account username to use creds")

	disconnectCmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from ipv6 prefix provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return connector.TunnelDisconnect(username)
		},
	}

	disconnectCmd.Flags().StringVarP(&username, "username", "u", "", "Account username to use creds")
	baseCmd.AddCommand(connectCmd)
	baseCmd.AddCommand(disconnectCmd)

	return &ConnectionsCmd{
		username: &username,

		connectCmd:    connectCmd,
		disconnectCmd: disconnectCmd,
	}
}
