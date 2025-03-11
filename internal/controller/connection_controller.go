package controller

import (
	"errors"
	"implementation/internal/domain/connections"
)

type ConnectionService interface {
	StartConnection(username string) error
	TerminateConnection(username string) error
	Status() connections.ConnectionStatus
}

type connectionsHandler func(name string) error

type ConnectionController struct {
	service           ConnectionService
	handlersWithNames map[string]connectionsHandler
}

func NewConnectionController(service ConnectionService) *ConnectionController {
	controller := &ConnectionController{
		service: service,
	}

	handlersWithNames := map[string]connectionsHandler{
		"connect":    controller.TunnelConnect,
		"disconnect": controller.TunnelDisconnect,
	}

	controller.handlersWithNames = handlersWithNames

	return controller
}

func (c *ConnectionController) TunnelConnect(username string) error {
	return c.service.StartConnection(username)
}

func (c *ConnectionController) TunnelDisconnect(username string) error {
	return c.service.TerminateConnection(username)
}

func (c *ConnectionController) TunnelStatus() connections.ConnectionStatus {
	return c.service.Status()
}

func (c *ConnectionController) Proceed(args []string) error {
	if len(args) <= 1 {
		return errors.New("invalid arguments")
	}

	return c.handlersWithNames[args[0]](args[1])
}
