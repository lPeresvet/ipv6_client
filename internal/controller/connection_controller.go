package controller

import (
	"implementation/internal/domain/connections"
)

type ConnectionService interface {
	StartConnection(username string) error
	TerminateConnection(username string) error
	Status() connections.ConnectionStatus
}

type ConnectionController struct {
	service ConnectionService
}

func NewConnectionController(service ConnectionService) *ConnectionController {
	return &ConnectionController{
		service: service,
	}
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
