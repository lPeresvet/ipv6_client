package service

import (
	"implementation/client_src/pkg/adapter"
	"implementation/connection_watcher/internal/domain"
)

type InterfaceAdapter interface {
}

type StatusService struct {
}

func NewStatusService() *StatusService {
	return &StatusService{}
}

func (s *StatusService) GetStatus(interfaceName string) (domain.ConnectionStatus, error) {
	ifaceInfo, err := adapter.GetTunnelInterfaceByName(interfaceName)
	if err != nil {
		return domain.Disconnected, err
	}

	status := domain.TunnelUP

	for _, address := range ifaceInfo.Addresses {
		if address.IP.To4() == nil && !address.IP.IsLinkLocalUnicast() {
			status = domain.IPv6UP

			break
		}
	}

	return status, nil
}
