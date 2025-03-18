package service

import (
	"errors"
	"implementation/internal/service/adapters/network"
	"log"
	"time"
)

type IfaceService struct{}

func NewIfaceService() *IfaceService {
	return &IfaceService{}
}

func (i IfaceService) GetIpv6Address(interfaceName string) (string, error) {
	log.Printf("Get ipv6 address of interface: %s", interfaceName)

	for attempt := 0; attempt < 8; attempt++ {
		info, err := network.GetTunnelInterfaceByName(interfaceName)
		if err != nil {
			return "", err
		}

		for _, address := range info.Addresses {
			log.Printf("Scaning %s ...", address.String())

			if address.IP.To4() == nil {
				return address.IP.String(), nil
			}
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("failed to get ipv6 address")
}
