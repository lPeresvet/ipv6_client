package controller

import (
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"os"
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

func (l *UnixSocketListener) ListenIpUp(ctx context.Context, control chan string) error {
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

		if err := l.HandleConnection(ctx, control, conn); err != nil {
			return err
		}
	}
}

func (l *UnixSocketListener) HandleConnection(ctx context.Context, control chan string, c net.Conn) error {
	received := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		count, err := c.Read(buf)
		received = append(received, buf[:count]...)
		if err != nil {
			fmt.Printf("%s", string(received))
			if err != io.EOF {
				log.Printf("Error on read: %s", err)
			}
			break
		}
	}

	return nil
}
