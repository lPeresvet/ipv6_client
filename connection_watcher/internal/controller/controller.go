package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	domain_consts "implementation/connection_watcher/pkg/domain"
)

type FSMInterface interface {
	Run(ctx context.Context)
	GetStatus() domain_consts.State
}

var (
	ErrControllerNotStarted = errors.New("controller not started")
)

type WatcherController struct {
	fsm     FSMInterface
	stopFSM func()

	ch chan string
}

func NewWatcherController(fsm FSMInterface, ch chan string) *WatcherController {
	return &WatcherController{
		fsm: fsm,
		ch:  ch,
	}
}

func (c *WatcherController) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.stopFSM = cancel

	go c.fsm.Run(ctx)

	go func() {
		if err := c.listenSocket(ctx); err != nil {
			log.Printf("error listening socket: %v", err)

			if err := c.Stop("defer stopped watcher"); err != nil {
				c.ch <- err.Error()
			}
		}
	}()

	return nil
}

func (c *WatcherController) Stop(msg string) error {
	if c.stopFSM != nil {
		c.stopFSM()
		c.ch <- msg

		return nil
	}

	return fmt.Errorf("%s: failed to stop controller: %w", msg, ErrControllerNotStarted)
}

func (c *WatcherController) listenSocket(ctx context.Context) error {
	listenPath := domain_consts.StatusSocketPath
	os.Remove(listenPath)

	listener, err := net.Listen("unix", listenPath)
	if err != nil {
		log.Fatalf("Unable to listen on socket %s: %s", listenPath, err)
	}
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping controller socket: %s", listenPath)
			return nil
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error on accept: %s", err)
		}

		go c.handleConnection(ctx, conn)
	}
}

func (c *WatcherController) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}

		receivedData := string(buf[:n])
		receivedData = strings.ReplaceAll(receivedData, "\n", "")

		var response string

		switch receivedData {
		case string(domain_consts.GetStatus):
			response = string(c.fsm.GetStatus())
		case string(domain_consts.TurnOff):

			if err := c.Stop("stopped watcher, get message"); err != nil {
				c.ch <- err.Error()

				return
			}

			response = domain_consts.OK
		default:
			response = domain_consts.ErrorMessage
		}

		_, err = conn.Write([]byte(response))

		return
	}
}
