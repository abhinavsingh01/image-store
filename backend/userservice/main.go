package main

import (
	"fmt"
	config "userservice/configs"
	"userservice/configs/logger"
	"userservice/controllers"
	"userservice/db"
	"userservice/repository"
	"userservice/server"
	"userservice/services"

	"go.uber.org/dig"
)

func main() {
	// environment := flag.String("e", "development", "")
	// flag.Usage = func() {
	// 	os.Exit(1)
	// }
	// flag.Parse()
	// config.LoadConfig()
	// db.Init()
	// server.Init()

	container := dig.New()

	container.Provide(logger.NewLogger)

	container.Provide(config.LoadConfig)

	container.Provide(db.Init)

	container.Provide(repository.NewUserRepoImpl)

	container.Provide(services.NewUserServiceImpl)

	container.Provide(controllers.NewUserController)

	container.Provide(server.NewRouter)

	container.Provide(server.NewServer)

	err := container.Invoke(func(srv *server.Server) {
		srv.Init()
	})

	fmt.Println(err)

}
