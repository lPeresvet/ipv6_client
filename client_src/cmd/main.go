package main

import (
	"implementation/client_src/internal/cli"
	controller2 "implementation/client_src/internal/controller"
	"implementation/client_src/internal/repository"
	service2 "implementation/client_src/internal/service"
	linux2 "implementation/client_src/internal/service/adapters/linux"
	"log"
)

func main() {
	adapter := linux2.NewLinuxAdapter()
	demonProvider := linux2.NewSystemdProvider()

	connectService := service2.NewConnectionService(adapter, demonProvider)
	ctrl := controller2.NewConnectionController(connectService)

	repo := repository.NewFileRepository("")

	filler := linux2.NewConfigFiller("config/templates")
	configService := service2.NewConfigService(repo, filler, demonProvider)

	ifaceService := service2.NewIfaceService()
	listener := controller2.NewUnixSocketListener(ifaceService)

	clientCLI := cli.New(ctrl, configService, listener)
	if err := clientCLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
