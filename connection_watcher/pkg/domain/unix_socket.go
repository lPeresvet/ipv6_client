package domain

const (
	WatcherSocketPath = "/var/run/ipv6-watcher"
	StatusSocketPath  = "/var/run/ipv6-status"
)

type UnixSocketCommand string

var (
	IfaceUP UnixSocketCommand = "IFACE_UP"
)

type StatusSocketCommand string

var (
	GetStatus StatusSocketCommand = "GET_STATUS"
)

var (
	ErrorMessage = "undefined command"
)
