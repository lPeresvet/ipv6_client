package main

import (
	"implementation/internal/cli"
	"implementation/internal/controller"
	"implementation/internal/service"
	"implementation/internal/service/adapters/linux"
	"log"
)

func main() {
	adapter := linux.NewLinuxAdapter()
	connectService := service.NewConnectionService(adapter)
	ctrl := controller.NewConnectionController(connectService)

	//err := ctrl.TunnelConnect("kirill")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//repo := repository.NewFileRepository("config")
	//config, err := repo.GetConfig("config-example.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(config)

	//filler := linux.NewConfigFiller("config/templates")
	//
	//err := filler.FillConfig(&config.Config{Servers: []config.ServerConfig{
	//	{
	//		Address: "45.45.45.45",
	//		Users: []config.UserConfig{
	//			{
	//				Username: "admin",
	//				Password: "qwerty007",
	//			},
	//			{
	//				Username: "admin1",
	//				Password: "qwerty00723232",
	//			},
	//		},
	//	},
	//}})
	//if err != nil {
	//	log.Fatal(err)
	//}

	clientCLI := cli.New(ctrl)
	err := clientCLI.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
