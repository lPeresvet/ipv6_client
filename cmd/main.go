package main

import (
	"implementation/internal/controller"
	"implementation/internal/service"
	"implementation/internal/service/adapters"
	"log"
)

func main() {
	adapter := adapters.NewLinuxAdapter()
	connectService := service.NewConnectionService(adapter)
	ctrl := controller.NewConnectionController(connectService)

	err := ctrl.TunnelConnect("kirill")
	if err != nil {
		log.Fatal(err)
	}
}
