package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"implementation/client_src/internal/domain/connections"
	"log"
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
	TunnelStatus() connections.ConnectionStatus
}

func NewConnectCmd(baseCmd *cobra.Command, connector Connector, listener UnixSocketListener) *ConnectionsCmd {
	var username string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to ipv6 prefix provider.",
		RunE:  getConnectHandler(listener, connector, &username),
	}

	connectCmd.Flags().StringVarP(&username, "username", "u", "", "Account username to use creds")
	connectCmd.MarkFlagRequired("username")

	disconnectCmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from ipv6 prefix provider.",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch connector.TunnelStatus() {
			case connections.DOWN:
				fmt.Println("Tunnel is already Disconnected")

				return nil
			}

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

func getConnectHandler(listener UnixSocketListener, connector Connector, username *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		switch connector.TunnelStatus() {
		case connections.UP:
			fmt.Println("Tunnel is already Connected")

			return nil
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		ch := make(chan *connections.IfaceEvent)
		go func() {
			if err := listener.ListenIpUp(ctx, ch, *username); err != nil {
				close(ch)
				connector.TunnelDisconnect(*username)

				log.Fatal(err)
			}
		}()
		if err := connector.TunnelConnect(*username); err != nil {
			return err
		}

		select {
		case event, ok := <-ch:
			if ok && event.Type == connections.IfaceUpEvent {
				log.Printf("Tunnel connected. Your ipv6 address: %s", event.Data)
			}
		case <-time.After(10 * time.Second):
			log.Printf("Tunnel connection failed. Timeout")
			log.Printf("Disconnecting...")

			connector.TunnelDisconnect(*username)
		}

		return nil
	}
}
