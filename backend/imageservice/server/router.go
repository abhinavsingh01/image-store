package server

import (
	"imageservice/controllers"
	"imageservice/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	imageCtrl *controllers.Image
	auth      *middleware.Auth
}

func NewRouter(imageCtrl *controllers.Image, auth *middleware.Auth) *Router {
	return &Router{
		imageCtrl: imageCtrl,
		auth:      auth,
	}
}

func (r *Router) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	{
		albumgroup := v1.Group("image/album/:albumId", r.auth.AuthorizeAlbum())
		{
			albumgroup.GET("/images", r.imageCtrl.GetAllImagesOfAlbum)
			albumgroup.DELETE("/images", r.imageCtrl.DeleteAllImages)
			albumgroup.POST("/image/upload", r.imageCtrl.UploadImage)
		}
		imagegroup := v1.Group("image", r.auth.AuthorizeImage())
		{
			imagegroup.GET("/:imageId/metadata", r.imageCtrl.GetAllImagesById)
			imagegroup.DELETE("/:imageId", r.imageCtrl.DeleteImage)
			imagegroup.GET("/:imageId/download", r.imageCtrl.Download)

		}
	}
	return router
}
