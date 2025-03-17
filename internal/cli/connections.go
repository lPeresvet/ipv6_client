package cli

import (
	"github.com/spf13/cobra"
	"time"
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

func NewConnectCmd(baseCmd *cobra.Command, connector Connector, listener UnixSocketListener) *ConnectionsCmd {
	var username string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to ipv6 prefix provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ch := make(chan string)
			go listener.ListenIpUp(ctx, ch)

			if err := connector.TunnelConnect(username); err != nil {
				return err
			}

			time.Sleep(5 * time.Second)
			return nil
		},
	}

	connectCmd.Flags().StringVarP(&username, "username", "u", "", "Account username to use creds")
	connectCmd.MarkFlagRequired("username")

	disconnectCmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from ipv6 prefix provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return connector.TunnelDisconnect(username)
		},
	}

	disconnectCmd.Flags().StringVarP(&username, "username", "u", "", "Account username to use creds")
	disconnectCmd.MarkFlagRequired("username")

	baseCmd.AddCommand(connectCmd)
	baseCmd.AddCommand(disconnectCmd)

	return &ConnectionsCmd{
		username: &username,

		connectCmd:    connectCmd,
		disconnectCmd: disconnectCmd,
	}
}
