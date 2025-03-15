package connections

type ConnectionStatus bool

const (
	UP   ConnectionStatus = true
	DOWN ConnectionStatus = false
)

type DemonInfo struct {
	Status DemonStatus
}

type DemonStatus string

const (
	DemonActive   DemonStatus = "active"
	DemonInactive DemonStatus = "inactive"
)
