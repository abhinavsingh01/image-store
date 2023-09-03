package server

import (
	"albumservice/controllers"
	"albumservice/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	albumCtrl *controllers.AlbumController
	auth      *middleware.Auth
}

func NewRouter(albumCtrl *controllers.AlbumController, auth *middleware.Auth) *Router {
	return &Router{
		albumCtrl: albumCtrl,
		auth:      auth,
	}
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.RedirectTrailingSlash = false
	v1 := router.Group("v1")
	{
		userGroup := v1.Group("album")
		{
			userGroup.GET("/:id", r.albumCtrl.GetAlbumById)
			userGroup.GET("/all", r.albumCtrl.GetAlbumsForUser)
			userGroup.POST("/new", r.albumCtrl.CreateNewAlbum)
			userGroup.DELETE("/:id", r.albumCtrl.DeleteAlbum)
			// userGroup.PUT("/:id", r.albumCtrl.UpdateAlbum)
		}
	}
	return router
}
