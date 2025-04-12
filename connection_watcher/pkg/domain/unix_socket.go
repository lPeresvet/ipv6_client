package domain

const (
	WatcherSocketPath = "/tmp/ipv6_watcher"
)

type UnixSocketCommand string

var (
	IfaceUP UnixSocketCommand = "IFACE_UP"
)
