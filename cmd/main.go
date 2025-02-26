package main

import (
	"implementation/internal/controller"
	"implementation/internal/service"
	"implementation/internal/service/adapters/linux"
	"log"
	"os"
)

func main() {
	adapter := linux.NewLinuxAdapter()
	connectService := service.NewConnectionService(adapter)
	ctrl := controller.NewConnectionController(connectService)

	args := os.Args

	err := ctrl.Proceed(args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
