package connections

import "net"

const (
	IfaceUpCommand = "iface-up"
)

type EventType string

var (
	IfaceUpEvent EventType = "ip-up"
)

type InterfaceInfo struct {
	Name      string
	Addresses []*net.IPNet
}

type IfaceEvent struct {
	Type EventType
	Data string
}
