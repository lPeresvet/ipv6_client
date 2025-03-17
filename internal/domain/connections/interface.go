package connections

import "net"

type InterfaceInfo struct {
	Name      string
	Addresses []*net.IPNet
}
