package linux

import (
	"fmt"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"net"
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
	if err := exec.Command(watcherName).Run(); err != nil {
		return fmt.Errorf("failed to start %s demon: %w", watcherName, err)
	}

	return nil
}

func (w *WatcherProvider) Stop() error {
	addr := &net.UnixAddr{Name: domain_consts.WatcherSocketPath, Net: "unix"}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial unix socket: %w", err)
	}
	defer conn.Close()

	message := domain_consts.TurnOff

	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write to socket: %w", err)
	}

	return nil
}
