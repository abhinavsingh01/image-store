package main

import (
	"authservice/clients"
	config "authservice/configs"
	"authservice/configs/logger"
	"authservice/controllers"
	"authservice/server"
	"authservice/services"
	"fmt"

	"go.uber.org/dig"
)

func main() {

	// providing all deps to container for DI

	container := dig.New()

	container.Provide(logger.NewLogger)

	container.Provide(config.LoadConfig)

	container.Provide(clients.NewUserClient)

	container.Provide(services.NewAuthService)

	container.Provide(services.NewLoginService)

	container.Provide(services.NewResgister)

	container.Provide(controllers.NewAuthController)

	container.Provide(server.NewRouter)

	container.Provide(server.NewServer)

	err := container.Invoke(func(srv *server.Server) {
		srv.Init()
	})

	fmt.Println(err)
}
