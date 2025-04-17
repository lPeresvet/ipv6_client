package linux

import (
	"fmt"
	"implementation/client_src/pkg/adapter"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"os/exec"
)

const (
	watcherName = "./watcher"
)

type WatcherProvider struct{}

func NewWatcherProvider() *WatcherProvider {
	return &WatcherProvider{}
}

func (w *WatcherProvider) Start() error {
	if err := exec.Command(watcherName).Start(); err != nil {
		return fmt.Errorf("failed to start %s demon: %w", watcherName, err)
	}

	return nil
}

func (w *WatcherProvider) Stop() error {
	return adapter.SendMessageToSocket(domain_consts.StatusSocketPath, string(domain_consts.TurnOff))
}
