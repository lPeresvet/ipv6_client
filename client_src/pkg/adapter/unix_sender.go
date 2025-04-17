package adapter

import (
	"fmt"
	"net"
)

func SendMessageToSocket(unixName, message string) error {
	addr := &net.UnixAddr{Name: unixName, Net: "unix"}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial unix socket: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write to socket: %w", err)
	}

	return nil
}
