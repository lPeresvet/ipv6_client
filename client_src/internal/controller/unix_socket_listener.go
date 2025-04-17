package controller

import (
	"fmt"
	"golang.org/x/net/context"
	"implementation/client_src/internal/domain/config"
	"implementation/client_src/internal/domain/connections"
	"implementation/client_src/pkg/adapter"
	"implementation/connection_watcher/pkg/domain"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type InterfaceService interface {
	GetIpv6Address(interfaceName string) (string, error)
	PrepareIpUpScript() error
	StartNDPProcedure(ifaceName string) error
}

type UnixSocketListener struct {
	InterfaceService InterfaceService
}

func NewUnixSocketListener(ifaceService InterfaceService) *UnixSocketListener {
	return &UnixSocketListener{InterfaceService: ifaceService}
}

func (l *UnixSocketListener) ListenIpUp(ctx context.Context, control chan *connections.IfaceEvent, username string) error {
	if err := l.InterfaceService.PrepareIpUpScript(); err != nil {
		return fmt.Errorf("failed to prepare ip up script: %w", err)
	}

	os.Remove(config.UnixSocketName) //TODO проверить, а надо ли это вообще
	listener, err := net.Listen("unix", config.UnixSocketName)
	if err != nil {
		log.Fatalf("Unable to listen on socket %s: %s", config.UnixSocketName, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error on accept: %s", err)
		}

		if err := l.HandleConnection(ctx, control, conn, username); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (l *UnixSocketListener) HandleConnection(_ context.Context, control chan *connections.IfaceEvent, c net.Conn, username string) error {
	received := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		count, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error on read: %s", err)
			}
			break
		}

		received = append(received, buf[:count]...)

		if err := l.proceedIncomingUnixMessage(control, string(received), username); err != nil {
			return err
		}
	}

	return nil
}

func (l *UnixSocketListener) proceedIncomingUnixMessage(control chan *connections.IfaceEvent, message, username string) error {
	log.Printf("Proceed incoming unix message: %s", message)

	command := strings.Split(message, " ")

	if len(command) < 2 {
		return fmt.Errorf("invalid command: %s", message)
	}

	if command[0] == connections.IfaceUpCommand {
		log.Printf("Received %s event", connections.IfaceUpCommand)

		outMsg := fmt.Sprintf("%s %s %s", domain.IfaceUP, command[1], username)

		if err := adapter.SendMessageToSocket(domain.WatcherSocketPath, outMsg); err != nil {
			log.Printf("Error sending message to socket: %s", err)
		}

		if err := l.InterfaceService.StartNDPProcedure(strings.Trim(command[1], "\n")); err != nil {
			return fmt.Errorf("failed to start NDP procedure: %w", err)
		}

		ipv6address, err := l.InterfaceService.GetIpv6Address(strings.Trim(command[1], "\n"))
		if err != nil {
			return err
		}

		control <- &connections.IfaceEvent{
			Type: connections.IfaceUpEvent,
			Data: ipv6address,
		}
	} else {
		log.Printf("Received unknown event: %s", command[0])

		return fmt.Errorf("invalid command: %s", message)
	}

	return nil
}
