package states

import (
	"context"
	"fmt"
	"implementation/connection_watcher/internal/domain"
)

type Waiter interface {
	Wait(ctx context.Context) (*domain.Connection, error)
}

type Waiting struct {
	waitService Waiter
	repo        map[string]*domain.Connection
}

func NewWaiting(waiter Waiter, repo map[string]*domain.Connection) *Waiting {
	return &Waiting{
		waitService: waiter,
		repo:        repo,
	}
}

func (w *Waiting) Execute(ctx context.Context) domain.State {
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

			return domain.StateWatching
		} else {
			return domain.StateStopped
		}

	case <-ctx.Done():
		fmt.Printf("Context done\n")

		return domain.StateStopped
	}
}
