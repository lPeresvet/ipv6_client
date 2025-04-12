package domain

type State string

const (
	// StateWaiting waiting for iface up event
	StateWaiting State = "waiting"
	// StateWatching watching for connection status
	StateWatching = "watching"
	// StateReconnectingTunnel trying to reestablish connection
	StateReconnectingTunnel = "reconnecting_tunnel"
	// StateReconnectingIPv6 trying to reestablish connection
	StateReconnectingIPv6 = "reconnecting_ipv6"
	// StateStopped app is stopped by user
	StateStopped = "stopped"
)
