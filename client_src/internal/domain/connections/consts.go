package connections

type ConnectionStatus string

const (
	UP   ConnectionStatus = "UP"
	DOWN ConnectionStatus = "DOWN"
)

type DemonInfo struct {
	Status DemonStatus
}

type DemonStatus string

const (
	DemonActive   DemonStatus = "active"
	DemonInactive DemonStatus = "inactive"
)
