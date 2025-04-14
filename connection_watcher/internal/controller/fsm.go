package controller

import (
	"context"
	"fmt"
	fsm_states "implementation/connection_watcher/internal/controller/states"
	"implementation/connection_watcher/internal/domain"
	"sync"
)

type FSMStates map[domain.State]FSMState

type FSM struct {
	states FSMStates

	mu           sync.Mutex
	currentState domain.State
}

type Waiter interface {
	Wait(ctx context.Context) (*domain.Connection, error)
}

type ConnectionProvider interface {
	Connect(username string) error
}

type StatusProvider interface {
	GetStatus(interfaceName string) (domain.ConnectionStatus, error)
}

type IPv6Service interface {
	StartNDPProcedure(ifaceName string) error
}

type FSMState interface {
	Execute(ctx context.Context) domain.State
}

func NewFSM(
	waiter Waiter,
	statusService StatusProvider,
	connectionProvider ConnectionProvider,
	ipv6Service IPv6Service,
) *FSM {
	connectionInfoRepo := make(map[string]*domain.Connection)

	states := map[domain.State]FSMState{
		domain.StateWaiting:            fsm_states.NewWaiting(waiter, connectionInfoRepo),
		domain.StateWatching:           fsm_states.NewWatching(statusService, connectionInfoRepo),
		domain.StateReconnectingTunnel: fsm_states.NewReconnectingTunnel(connectionProvider, connectionInfoRepo),
		domain.StateReconnectingIPv6:   fsm_states.NewReconnectingIPv6(ipv6Service, connectionInfoRepo),
	}

	return &FSM{
		currentState: domain.StateWaiting,
		states:       states,
	}
}

func (fsm *FSM) Run(ctx context.Context) {
	for fsm.currentState != domain.StateStopped {
		fmt.Println(fsm.currentState)
		nextState := fsm.states[fsm.currentState].Execute(ctx)

		fmt.Printf("Transit %q -> %q\n", fsm.currentState, nextState)
		fsm.setState(nextState)
	}
}

func (fsm *FSM) GetStatus() domain.State {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	return fsm.currentState
}

func (fsm *FSM) setState(state domain.State) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	fsm.currentState = state
}
