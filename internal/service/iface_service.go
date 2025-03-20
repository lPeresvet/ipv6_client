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

	fmt.Println("ndp: bound to address:", ip)
	return nil
}
