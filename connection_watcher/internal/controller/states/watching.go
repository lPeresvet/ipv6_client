package states

import (
	"context"
	"fmt"
	"implementation/connection_watcher/internal/domain"
	"log"
	"time"
)

type Watching struct {
	statusService StatusProvider
	repo          map[string]*domain.Connection
}

type StatusProvider interface {
	GetStatus(interfaceName string) (domain.ConnectionStatus, error)
}

func NewWatching(service StatusProvider, repo map[string]*domain.Connection) *Watching {
	return &Watching{
		statusService: service,
		repo:          repo,
	}
}

func (w *Watching) Execute(ctx context.Context) domain.State {
	connection, ok := w.repo["data"]
	if !ok {
		fmt.Println("data not found")

		return domain.StateStopped
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Context done\n")

			return domain.StateStopped
		case <-time.After(5 * time.Second):
			status, err := w.statusService.GetStatus(connection.InterfaceName)
			if err != nil {
				fmt.Printf("failed to get status for %s: %v", connection.InterfaceName, err)

				return domain.StateStopped
			}

			log.Printf("%s is %s", connection.InterfaceName, status)

			switch status {
			case domain.Disconnected:
				return domain.StateReconnectingTunnel
			case domain.TunnelUP:
				return domain.StateReconnectingIPv6
			}
		}
	}
}
