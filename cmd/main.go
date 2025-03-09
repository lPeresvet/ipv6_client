package main

import (
	"implementation/internal/controller"
	"implementation/internal/service"
	"implementation/internal/service/adapters/linux"
	"log"
)

func main() {
	adapter := linux.NewLinuxAdapter()
	connectService := service.NewConnectionService(adapter)
	ctrl := controller.NewConnectionController(connectService)

	err := ctrl.TunnelConnect("kirill")
	if err != nil {
		log.Fatal(err)
	}
}
