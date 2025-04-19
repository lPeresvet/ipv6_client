package controller

import (
	"context"
	"fmt"
	"implementation/client_src/pkg/config"
	fsm_states "implementation/connection_watcher/internal/controller/states"
	"implementation/connection_watcher/internal/domain"
	domain_consts "implementation/connection_watcher/pkg/domain"
	"sync"
)

type FSMStates map[domain_consts.State]FSMState

type FSM struct {
	states FSMStates
	cfg    *config.WatcherConfig

	mu           sync.Mutex
	currentState domain_consts.State
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
	Execute(ctx context.Context) domain_consts.State
}

func NewFSM(
	cfg *config.WatcherConfig,
	waiter Waiter,
	statusService StatusProvider,
	connectionProvider ConnectionProvider,
	ipv6Service IPv6Service,
) *FSM {
	connectionInfoRepo := make(map[string]*domain.Connection)

	states := map[domain_consts.State]FSMState{
		domain_consts.StateWaiting:            fsm_states.NewWaiting(cfg, waiter, connectionInfoRepo),
		domain_consts.StateWatching:           fsm_states.NewWatching(cfg, statusService, connectionInfoRepo),
		domain_consts.StateReconnectingTunnel: fsm_states.NewReconnectingTunnel(connectionProvider, connectionInfoRepo),
		domain_consts.StateReconnectingIPv6:   fsm_states.NewReconnectingIPv6(ipv6Service, connectionInfoRepo),
	}

	return &FSM{
		cfg:          cfg,
		currentState: domain_consts.StateWaiting,
		states:       states,
	}
}

func (fsm *FSM) Run(ctx context.Context) {
	for fsm.currentState != domain_consts.StateStopped {
		fmt.Println(fsm.currentState)
		nextState := fsm.states[fsm.currentState].Execute(ctx)

		fmt.Printf("Transit %q -> %q\n", fsm.currentState, nextState)
		fsm.setState(nextState)
	}
}

func (fsm *FSM) GetStatus() domain_consts.State {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	return fsm.currentState
}

func (fsm *FSM) setState(state domain_consts.State) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	fsm.currentState = state
}
