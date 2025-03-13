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
	connectService := service.NewConnectionService(adapter)
	ctrl := controller.NewConnectionController(connectService)

	repo := repository.NewFileRepository("")

	filler := linux.NewConfigFiller("config/templates")
	configService := service.NewConfigService(repo, filler)

	clientCLI := cli.New(ctrl, configService)
	if err := clientCLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
