package states

import (
	"context"
	"fmt"
	"implementation/client_src/pkg/config"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"log"
	"time"
)

type Watching struct {
	cfg *config.WatcherConfig

	statusService StatusProvider
	repo          map[string]*domain.Connection
}

type StatusProvider interface {
	GetStatus(interfaceName string) (domain.ConnectionStatus, error)
}

func NewWatching(cfg *config.WatcherConfig, service StatusProvider, repo map[string]*domain.Connection) *Watching {
	return &Watching{
		cfg:           cfg,
		statusService: service,
		repo:          repo,
	}
}

func (w *Watching) Execute(ctx context.Context) domain_consts.State {
	connection, ok := w.repo["data"]
	if !ok {
		fmt.Println("data not found")

		return domain_consts.StateStopped
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Context done\n")

			return domain_consts.StateStopped
		case <-time.After(time.Duration(w.cfg.Reconnect.WatchingPeriod) * time.Second):
			status, err := w.statusService.GetStatus(connection.InterfaceName)
			if err != nil {
				fmt.Printf("failed to get status for %s: %v", connection.InterfaceName, err)

				return domain_consts.StateStopped
			}

			log.Printf("%s is %s", connection.InterfaceName, status)

			switch status {
			case domain.Disconnected:
				return domain_consts.StateReconnectingTunnel
			case domain.TunnelUP:
				return domain_consts.StateReconnectingIPv6
			}
		}
	}
}
