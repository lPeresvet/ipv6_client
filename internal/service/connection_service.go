package service

import (
	"fmt"
	"implementation/internal/domain/connections"
)

type ConnectionProvider interface {
	Connect(username string) error
	Disconnect(username string) error
}
type ConnectionService struct {
	status             connections.ConnectionStatus
	connectionProvider ConnectionProvider
}

func NewConnectionService(provider ConnectionProvider) *ConnectionService {
	return &ConnectionService{
		status:             connections.DOWN,
		connectionProvider: provider,
	}
}

func (service *ConnectionService) Status() connections.ConnectionStatus {
	return service.status
}

func (service *ConnectionService) StartConnection(username string) error {
	if service.status == connections.UP {
		return fmt.Errorf("connection already started")
	}

	if err := service.connectionProvider.Connect(username); err != nil {
		service.status = connections.DOWN

		return fmt.Errorf("failed to start connection: %w", err)
	}

	service.status = connections.UP

	return nil
}

func (service *ConnectionService) TerminateConnection(username string) error {
	if service.status == connections.DOWN {
		return fmt.Errorf("connection already terminated")
	}

	if err := service.connectionProvider.Disconnect(username); err != nil {
		service.status = connections.UP

		return fmt.Errorf("failed to terminate connection: %w", err)
	}

	service.status = connections.DOWN

	return nil
}
