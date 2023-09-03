package main

import (
	"fmt"
	config "imageservice/configs"
	"imageservice/configs/logger"
	"imageservice/controllers"
	"imageservice/db"
	"imageservice/middleware"
	"imageservice/repository"
	"imageservice/server"
	"imageservice/services"
	"imageservice/utils"

	"go.uber.org/dig"
)

func main() {

	container := dig.New()

	container.Provide(config.LoadConfig)

	container.Provide(logger.NewLogger)

	container.Provide(db.Init)

	container.Provide(utils.NewImageUtils)

	container.Provide(repository.NewImageRepoImpl)

	container.Provide(repository.NewAuthRepoImpl)

	container.Provide(services.NewImageServiceImpl)

	container.Provide(controllers.NewImage)

	container.Provide(middleware.NewAuth)

	container.Provide(server.NewRouter)

	container.Provide(server.NewServer)

	err := container.Invoke(func(srv *server.Server) {
		srv.Init()
	})

	fmt.Println(err)

}
