package controller

import (
	"fmt"
	"golang.org/x/net/context"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type WaitingController struct{}

func NewWaitingController() *WaitingController {
	return &WaitingController{}
}

func (controller *WaitingController) Wait(ctx context.Context) (*domain.Connection, error) {
	os.Remove(domain_consts.WatcherSocketPath) //TODO проверить, а надо ли это вообще
	listener, err := net.Listen("unix", domain_consts.WatcherSocketPath)
	if err != nil {
		log.Fatalf("Unable to listen on socket %s: %s", domain_consts.WatcherSocketPath, err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("error on accept: %s", err)
	}
	defer conn.Close()

	return listen(ctx, conn)
}

func listen(ctx context.Context, connection net.Conn) (*domain.Connection, error) {
	fmt.Printf("Listening on %s\n", connection.RemoteAddr())

	if err := connection.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %s", err)
	}

	received := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		count, err := connection.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("error on read: %s", err)
			}

			break
		}

		received = append(received, buf[:count]...)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return proceedIncomingUnixMessage(string(received))
		}
	}

	return nil, nil
}

func proceedIncomingUnixMessage(message string) (*domain.Connection, error) {
	log.Printf("Proceed incoming unix message: %s", message)

	command := strings.Split(message, " ")

	if len(command) < 3 {
		return nil, fmt.Errorf("invalid command format: %s", message)
	}

	if command[0] == string(domain_consts.IfaceUP) {
		log.Printf("Received %s event", domain_consts.IfaceUP)

		return &domain.Connection{
			InterfaceName: command[1],
			Username:      strings.ReplaceAll(command[2], "\n", ""),
		}, nil
	}

	return nil, fmt.Errorf("invalid command: %s", command[0])
}
