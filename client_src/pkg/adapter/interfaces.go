package adapter

import (
	"fmt"
	"implementation/client_src/internal/domain/connections"
	"net"
)

func GetTunnelInterfaceByName(ifaceName string) (*connections.InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return &connections.InterfaceInfo{}, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	addresses := make([]*net.IPNet, 0)

	for _, iface := range interfaces {
		//TODO подобрать битовую маску
		if iface.Name == ifaceName {
			addrs, err := iface.Addrs()
			if err != nil {
				return &connections.InterfaceInfo{}, fmt.Errorf("failed to get '%s' addresses: %w", ifaceName, err)
			}

			for _, add := range addrs {
				if ip, ok := add.(*net.IPNet); ok {
					addresses = append(addresses, ip)
				}
			}

			break
		}
	}

	if len(addresses) == 0 {
		return &connections.InterfaceInfo{}, fmt.Errorf("failed to get '%s' interface", ifaceName)
	}

	return &connections.InterfaceInfo{
		Name:      ifaceName,
		Addresses: addresses,
	}, nil
}
