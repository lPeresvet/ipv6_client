package main

import (
	"golang.org/x/net/context"
	"implementation/client_src/pkg/adapter"
	ipv6service "implementation/client_src/pkg/service"
	"implementation/connection_watcher/internal/controller"
	"implementation/connection_watcher/internal/service"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	waiter := controller.NewWaitingController()
	statusService := service.NewStatusService()
	ipv6Service := ipv6service.NewIfaceService()
	connectionProvider := adapter.NewLinuxAdapter()

	fsm := controller.NewFSM(waiter, statusService, connectionProvider, ipv6Service)

	exitStatus := make(chan string)

	watcherController := controller.NewWatcherController(fsm, exitStatus)
	watcherController.Start(ctx)

	select {
	case msg := <-exitStatus:
		log.Println("[INFO] context canceled: ", msg)
	}
}
