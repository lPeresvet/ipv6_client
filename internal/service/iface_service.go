package service

import (
	"errors"
	"fmt"
	"github.com/mdlayher/ndp"
	"implementation/internal/domain/config"
	"implementation/internal/domain/connections"
	"implementation/internal/domain/template"
	"implementation/internal/parsers"
	"implementation/internal/service/adapters/network"
	"net"
	"os"
	"time"
)

type IfaceService struct{}

func NewIfaceService() *IfaceService {
	return &IfaceService{}
}

func (i *IfaceService) GetIpv6Address(interfaceName string) (string, error) {
	for attempt := 0; attempt < 5; attempt++ {
		info, err := network.GetTunnelInterfaceByName(interfaceName)
		if err != nil {
			return "", err
		}

		for _, address := range info.Addresses {
			if address.IP.To4() == nil && !address.IP.IsLinkLocalUnicast() {
				return address.IP.String(), nil
			}
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("failed to get ipv6 address")
}

func (i *IfaceService) PrepareIpUpScript() error {
	cmd := fmt.Sprintf(`echo "%s $INTERFACE" | nc -U %s`, connections.IfaceUpCommand, config.UnixSocketName)

	if !parsers.IsFileExists(connections.IfaceUpScriptPath) {
		if err := createFileWithShebang(template.BashShebang); err != nil {
			return err
		}
	}

	scriptExists, err := parsers.IsContainsLineStartWith(connections.IfaceUpScriptPath, cmd)
	if err != nil {
		return err
	}

	if !scriptExists {
		if err := parsers.AppendToFileByPath(connections.IfaceUpScriptPath, cmd+"\n"); err != nil {
			return err
		}
	}

	return nil
}

func createFileWithShebang(shebang string) error {
	file, err := os.OpenFile(connections.IfaceUpScriptPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(shebang + "\n\n"); err != nil {
		return err
	}

	return nil
}

func (i *IfaceService) StartNDPProcedure(ifaceName string) error {
	ifi, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("failed to get interface: %v", err)
	}

	// Set up an *ndp.Conn, bound to this interface's link-local IPv6 address.
	c, ip, err := ndp.Listen(ifi, ndp.LinkLocal)
	if err != nil {
		return fmt.Errorf("failed to dial NDP connection: %v", err)
	}
	// Clean up after the connection is no longer needed.
	defer c.Close()

	// Use target's solicited-node multicast address to request that the target
	// respond with a neighbor advertisement.
	snm, err := ndp.SolicitedNodeMulticast(ip)
	if err != nil {
		return fmt.Errorf("failed to determine solicited-node multicast address: %v", err)
	}

	// Build a neighbor solicitation message, indicate the target's link-local
	// address, and also specify our source link-layer address.
	m := &ndp.NeighborSolicitation{
		TargetAddress: ip,
		Options: []ndp.Option{
			&ndp.LinkLayerAddress{
				Direction: ndp.Source,
				Addr:      ifi.HardwareAddr,
			},
		},
	}

	// Send the multicast message and wait for a response.
	if err := c.WriteTo(m, nil, snm); err != nil {
		return fmt.Errorf("failed to write neighbor solicitation: %v", err)
	}
	msg, _, from, err := c.ReadFrom()
	if err != nil {
		return fmt.Errorf("failed to read NDP message: %v", err)
	}

	// Expect a neighbor advertisement message with a target link-layer
	// address option.
	na, ok := msg.(*ndp.NeighborAdvertisement)
	if !ok {
		return fmt.Errorf("message is not a neighbor advertisement: %T", msg)
	}
	if len(na.Options) != 1 {
		return fmt.Errorf("expected one option in neighbor advertisement")
	}
	tll, ok := na.Options[0].(*ndp.LinkLayerAddress)
	if !ok {
		return fmt.Errorf("option is not a link-layer address: %T", msg)
	}

	fmt.Printf("ndp: neighbor advertisement from %s:\n", from)
	fmt.Printf("  - solicited: %t\n", na.Solicited)
	fmt.Printf("  - link-layer address: %s\n", tll.Addr)

	return nil
}
