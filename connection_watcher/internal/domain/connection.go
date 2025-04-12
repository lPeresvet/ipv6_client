package domain

type Connection struct {
	InterfaceName string
	Username      string
}

type ConnectionStatus string

const (
	IPv6UP       ConnectionStatus = "ipv6_up"
	TunnelUP     ConnectionStatus = "tunnel_up"
	Disconnected ConnectionStatus = "disconnected"
)
