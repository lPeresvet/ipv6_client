package main

import (
	"implementation/client_src/internal/cli"
	controller "implementation/client_src/internal/controller"
	service "implementation/client_src/internal/service"
	linux "implementation/client_src/internal/service/adapters/linux"
	linux_adapter "implementation/client_src/pkg/adapter"
	"implementation/client_src/pkg/repository"
	ipv6service "implementation/client_src/pkg/service"
	"log"
)

func main() {
	adapter := linux_adapter.NewLinuxAdapter()
	demonProvider := linux.NewSystemdProvider()

	connectService := service.NewConnectionService(adapter, demonProvider)
	ctrl := controller.NewConnectionController(connectService)

	repo := repository.NewFileRepository("")

	filler := linux.NewConfigFiller("config/templates")
	configService := service.NewConfigService(repo, filler, demonProvider)

	ifaceService := ipv6service.NewIfaceService()
	listener := controller.NewUnixSocketListener(ifaceService)

	clientCLI := cli.New(ctrl, configService, listener)
	if err := clientCLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
