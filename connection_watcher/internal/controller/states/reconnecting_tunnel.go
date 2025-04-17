package states

import (
	"context"
	"fmt"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
)

type ReconnectingTunnel struct {
	connectionService ConnectionProvider
	repo              map[string]*domain.Connection
}

type ConnectionProvider interface {
	Connect(username string) error
}

func NewReconnectingTunnel(service ConnectionProvider, repo map[string]*domain.Connection) *ReconnectingTunnel {
	return &ReconnectingTunnel{
		connectionService: service,
		repo:              repo,
	}
}

func (r *ReconnectingTunnel) Execute(ctx context.Context) domain_consts.State {
	connection, ok := r.repo["data"]
	if !ok {
		fmt.Println("connection data not found")

		return domain_consts.StateStopped
	}

	if err := r.connectionService.Connect(connection.Username); err != nil {
		return domain_consts.StateStopped
	}

	return domain_consts.StateWatching
}
