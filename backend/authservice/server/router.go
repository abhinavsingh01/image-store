package server

import (
	"authservice/controllers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ctrl *controllers.AuthController
}

func NewRouter(ctrl *controllers.AuthController) *Router {
	return &Router{
		ctrl: ctrl,
	}
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	{
		userGroup := v1.Group("auth")
		{
			userGroup.POST("/register", r.ctrl.Register)
			userGroup.POST("/login", r.ctrl.Login)
		}
	}
	return router
}
