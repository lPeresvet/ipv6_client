package linux

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	l2tpControlPipePath = "/var/run/xl2tpd/l2tp-control"
)

type LinuxAdapter struct {
}

func NewLinuxAdapter() *LinuxAdapter {
	return &LinuxAdapter{}
}

func (l *LinuxAdapter) Connect(username string) error {
	message := fmt.Sprintf("c %s\n", username)

	if err := sendCommand(message); err != nil {
		return err
	}

	return nil
}

func (l *LinuxAdapter) Disconnect(username string) error {
	message := fmt.Sprintf("d %s\n", username)

	if err := sendCommand(message); err != nil {
		return err
	}

	return nil
}

func sendCommand(message string) error {
	pipe, err := os.OpenFile(l2tpControlPipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return fmt.Errorf("failed to open l2tp control pipe: %w", err)
	}
	defer pipe.Close()

	writer := bufio.NewWriter(pipe)

	if _, err = writer.WriteString(message); err != nil {
		log.Fatalf("failed to write '%s' to l2tp control pipe: %v", message, err)
	}

	if err := writer.Flush(); err != nil {
		log.Fatalf("failed to flush buffer to l2tp control pipe: %v", err)
	}

	return nil
}
