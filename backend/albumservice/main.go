package main

import (
	"albumservice/client"
	config "albumservice/configs"
	"albumservice/configs/logger"
	"albumservice/controllers"
	"albumservice/db"
	"albumservice/middleware"
	"albumservice/repository"
	"albumservice/server"
	"albumservice/services"
	"fmt"

	"go.uber.org/dig"
)

func main() {

	container := dig.New()

	container.Provide(logger.NewLogger)

	container.Provide(config.LoadConfig)

	container.Provide(db.Init)

	container.Provide(client.NewImageClient)

	container.Provide(repository.NewAlbumRepoImpl)

	container.Provide(services.NewAlbumServiceImpl)

	container.Provide(controllers.NewAlbumController)

	container.Provide(middleware.NewAuth)

	container.Provide(server.NewRouter)

	container.Provide(server.NewServer)

	err := container.Invoke(func(srv *server.Server) {
		srv.Init()
	})

	fmt.Println(err)

}
