package network

import (
	"fmt"
	"implementation/internal/domain/connections"
	"log"
	"net"
)

func GetTunnelInterfaceByName(ifaceName string) (*connections.InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return &connections.InterfaceInfo{}, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	ifaceInfo := connections.InterfaceInfo{
		Name:      ifaceName,
		Addresses: make([]*net.IPNet, 0),
	}

	for _, iface := range interfaces {
		//TODO подобрать битовую маску
		log.Printf("Check interface %s == %s -> %v", iface.Name, ifaceName, len(ifaceName))
		if iface.Name == ifaceName {
			log.Printf("Found ptp interface %s", iface.Name)

			addrs, err := iface.Addrs()
			if err != nil {
				return &connections.InterfaceInfo{}, fmt.Errorf("failed to get '%s' addresses: %w", ifaceName, err)
			}

			for _, add := range addrs {
				if ip, ok := add.(*net.IPNet); ok {
					ifaceInfo.Addresses = append(ifaceInfo.Addresses, ip)
				}
			}

			break
		}
	}

	return &ifaceInfo, nil
}
