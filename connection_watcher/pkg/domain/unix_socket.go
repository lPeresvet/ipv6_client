package domain

const (
	WatcherSocketPath = "/var/run/ipv6-watcher"
)

type UnixSocketCommand string

var (
	IfaceUP UnixSocketCommand = "IFACE_UP"
)
