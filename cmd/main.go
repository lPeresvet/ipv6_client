package main

import (
	"implementation/internal/cli"
	"implementation/internal/controller"
	"implementation/internal/repository"
	"implementation/internal/service"
	"implementation/internal/service/adapters/linux"
	"log"
)

func main() {
	adapter := linux.NewLinuxAdapter()
	demonProvider := linux.NewSystemdProvider()

	connectService := service.NewConnectionService(adapter, demonProvider)
	ctrl := controller.NewConnectionController(connectService)

	repo := repository.NewFileRepository("")

	filler := linux.NewConfigFiller("config/templates")
	configService := service.NewConfigService(repo, filler, demonProvider)

	ifaceService := service.NewIfaceService()
	listener := controller.NewUnixSocketListener(ifaceService)

	clientCLI := cli.New(ctrl, configService, listener)
	if err := clientCLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
