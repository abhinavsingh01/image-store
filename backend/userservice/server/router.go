package server

import (
	"userservice/controllers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userCtrl *controllers.UserController
}

func NewRouter(userCtrl *controllers.UserController) *Router {
	return &Router{
		userCtrl: userCtrl,
	}
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	{
		userGroup := v1.Group("user")
		{
			userGroup.GET("/:id", r.userCtrl.GetUserById)
			userGroup.POST("/find", r.userCtrl.FindUser)
			userGroup.POST("", r.userCtrl.CreateNewUser)
		}
	}
	return router
}
