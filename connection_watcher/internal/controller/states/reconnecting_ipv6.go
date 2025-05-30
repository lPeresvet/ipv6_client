package states

import (
	"context"
	"fmt"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
)

type ReconnectingIPv6 struct {
	ipv6Service IPv6Service
	repo        map[string]*domain.Connection
}

type IPv6Service interface {
	StartNDPProcedure(ifaceName string) error
}

func NewReconnectingIPv6(service IPv6Service, repo map[string]*domain.Connection) *ReconnectingIPv6 {
	return &ReconnectingIPv6{
		ipv6Service: service,
		repo:        repo,
	}
}

func (r *ReconnectingIPv6) Execute(ctx context.Context) domain_consts.State {
	connection, ok := r.repo["data"]
	if !ok {
		fmt.Println("connection data not found")

		return domain_consts.StateStopped
	}

	if err := r.ipv6Service.StartNDPProcedure(connection.Username); err != nil {
		return domain_consts.StateReconnectingTunnel
	}

	return domain_consts.StateWatching
}
