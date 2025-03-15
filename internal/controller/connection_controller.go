package controller

import (
	"errors"
	"implementation/internal/domain/connections"
)

type ConnectionService interface {
	StartConnection(username string) error
	TerminateConnection(username string) error
	Status() connections.ConnectionStatus
	GetDemonInfo() (*connections.DemonInfo, error)
}

type ConnectionController struct {
	service ConnectionService
}

func NewConnectionController(service ConnectionService) *ConnectionController {
	controller := &ConnectionController{
		service: service,
	}

	return controller
}

func (c *ConnectionController) TunnelConnect(username string) error {
	info, err := c.service.GetDemonInfo()
	if err != nil {
		return err
	}

	if info.Status != connections.DemonActive {
		return errors.New("connection is not active")
	}

	return c.service.StartConnection(username)
}

func (c *ConnectionController) TunnelDisconnect(username string) error {
	return c.service.TerminateConnection(username)
}

func (c *ConnectionController) TunnelStatus() connections.ConnectionStatus {
	return c.service.Status()
}
