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

func NewConnectCmd(baseCmd *cobra.Command, connector Connector, listener UnixSocketListener) *ConnectionsCmd {
	var username string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to ipv6 prefix provider.",
		RunE:  getConnectHandler(listener, connector, username),
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

func getConnectHandler(listener UnixSocketListener, connector Connector, username string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//ch := make(chan *connections.IfaceEvent)
		//go func() {
		//	if err := listener.ListenIpUp(ctx, ch); err != nil {
		//		close(ch)
		//		log.Fatal(err)
		//	}
		//}()

		if err := connector.TunnelConnect(username); err != nil {
			return err
		}

		//select {
		//case event := <-ch:
		//	if event.Type == connections.IfaceUpEvent {
		//		log.Printf("Tunnel connected. Your ipv6 address: %s", event.Data)
		//	}
		//case <-time.After(5 * time.Second):
		//	log.Printf("Tunnel connection failed. Timeout")
		//
		//	//connector.TunnelDisconnect(username)
		//
		//	ctx.Done()
		//}

		return nil
	}
}
