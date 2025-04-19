package states

import (
	"context"
	"fmt"
	"implementation/client_src/pkg/config"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"time"
)

type Waiter interface {
	Wait(ctx context.Context) (*domain.Connection, error)
}

type Waiting struct {
	cfg *config.WatcherConfig

	waitService Waiter
	repo        map[string]*domain.Connection
}

func NewWaiting(cfg *config.WatcherConfig, waiter Waiter, repo map[string]*domain.Connection) *Waiting {
	return &Waiting{
		cfg:         cfg,
		waitService: waiter,
		repo:        repo,
	}
}

func (w *Waiting) Execute(ctx context.Context) domain_consts.State {
	ch := make(chan *domain.Connection)
	go func() {
		fmt.Printf("Waiting for connections to be ready...\n")
		info, err := w.waitService.Wait(ctx)
		if err != nil {
			fmt.Printf("Error in waiting state: %v\n", err)

			ch <- nil

			return
		}

		ch <- info
	}()

	select {
	case info := <-ch:
		if info != nil {
			w.repo["data"] = info

			return domain_consts.StateWatching
		} else {
			return domain_consts.StateStopped
		}
	case <-time.After(time.Duration(w.cfg.Reconnect.WaitingTimeout) * time.Second):
		fmt.Printf("Timeout waiting for connections to be ready...\n")

		return domain_consts.StateStopped
	case <-ctx.Done():
		fmt.Printf("Context done\n")

		return domain_consts.StateStopped
	}
}
