package connections

type ConnectionStatus bool

const (
	UP   ConnectionStatus = true
	DOWN ConnectionStatus = false
)
