package controller

import (
	"fmt"
	"golang.org/x/net/context"
	"implementation/internal/domain/connections"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const unixSocketName = "/var/run/ipv6-client"

type InterfaceService interface {
	GetIpv6Address(interfaceName string) (string, error)
}

type UnixSocketListener struct {
	InterfaceService InterfaceService
}

func NewUnixSocketListener(ifaceService InterfaceService) *UnixSocketListener {
	return &UnixSocketListener{InterfaceService: ifaceService}
}

func (l *UnixSocketListener) ListenIpUp(ctx context.Context, control chan *connections.IfaceEvent) error {
	os.Remove(unixSocketName)
	listener, err := net.Listen("unix", unixSocketName)
	if err != nil {
		log.Fatalf("Unable to listen on socket %s: %s", unixSocketName, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error on accept: %s", err)
		}

		log.Printf("New connection from %s", conn.RemoteAddr())

		if err := l.HandleConnection(ctx, control, conn); err != nil {
			return err
		}
	}
}

func (l *UnixSocketListener) HandleConnection(ctx context.Context, control chan *connections.IfaceEvent, c net.Conn) error {
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

		if err := l.proceedIncomingUnixMessage(control, string(received)); err != nil {
			return err
		}
	}

	return nil
}

func (l *UnixSocketListener) proceedIncomingUnixMessage(control chan *connections.IfaceEvent, message string) error {
	log.Printf("Proceed incoming unix message: %s", message)

	command := strings.Split(message, " ")

	if len(command) < 2 {
		return fmt.Errorf("invalid command: %s", message)
	}

	if command[0] == connections.IfaceUpCommand {
		log.Printf("Received %s event", connections.IfaceUpCommand)

		control <- &connections.IfaceEvent{
			Type: connections.IfaceUpEvent,
			Data: command[1],
		}
	} else {
		log.Printf("Received unknown event: %s", command[0])

		return fmt.Errorf("invalid command: %s", message)
	}

	return nil
}
