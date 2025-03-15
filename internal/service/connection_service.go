package service

import (
	"fmt"
	"implementation/internal/domain/connections"
)

type ConnectionProvider interface {
	Connect(username string) error
	Disconnect(username string) error
}

type DemonProvider interface {
	StartDemon(demonName string) error
	StopDemon(demonName string) error
	DemonStatus(demonName string) (*connections.DemonInfo, error)
}
type ConnectionService struct {
	status             connections.ConnectionStatus
	connectionProvider ConnectionProvider
	demonProvider      DemonProvider
}

func NewConnectionService(provider ConnectionProvider, demonProvider DemonProvider) *ConnectionService {
	return &ConnectionService{
		status:             connections.DOWN,
		connectionProvider: provider,
		demonProvider:      demonProvider,
	}
}

func (service *ConnectionService) Status() connections.ConnectionStatus {
	return service.status
}

func (service *ConnectionService) StartConnection(username string) error {
	if err := service.connectionProvider.Connect(username); err != nil {
		service.status = connections.DOWN

		return fmt.Errorf("failed to start connection: %w", err)
	}

	service.status = connections.UP

	return nil
}

func (service *ConnectionService) TerminateConnection(username string) error {
	if err := service.connectionProvider.Disconnect(username); err != nil {
		service.status = connections.UP

		return fmt.Errorf("failed to terminate connection: %w", err)
	}

	service.status = connections.DOWN

	return nil
}

func (service *ConnectionService) InitDemon(name string) error {
	return service.demonProvider.StartDemon(name)
}

func (service *ConnectionService) StopDemon(name string) error {
	return service.demonProvider.StopDemon(name)
}

func (service *ConnectionService) GetDemonInfo() (*connections.DemonInfo, error) {
	return service.demonProvider.DemonStatus(demonName)
}
