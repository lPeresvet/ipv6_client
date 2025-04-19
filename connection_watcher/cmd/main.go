package main

import (
	"fmt"
	"golang.org/x/net/context"
	"implementation/client_src/pkg/adapter"
	"implementation/client_src/pkg/repository"
	ipv6service "implementation/client_src/pkg/service"
	"implementation/connection_watcher/internal/config"
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

	repo := repository.NewFileRepository("")

	loader := config.NewLoader(repo)

	cfg, err := loader.Load("config/config-example.yaml")

	fmt.Println(cfg)

	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	fsm := controller.NewFSM(cfg, waiter, statusService, connectionProvider, ipv6Service)

	exitStatus := make(chan string)

	watcherController := controller.NewWatcherController(fsm, exitStatus)
	watcherController.Start(ctx)

	select {
	case msg := <-exitStatus:
		log.Println("[INFO] context canceled: ", msg)
	}
}
