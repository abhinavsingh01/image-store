package server

import (
	"fmt"
	"image-store-api-gateway/middleware"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(auth middleware.Auth) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
	// router.RedirectTrailingSlash = false

	gw := router.Group("gw")
	{
		apiAuth := gw.Group("api", auth.Authorize())
		{
			v1 := apiAuth.Group("v1")
			v1.Any("/user/*path", createReverseProxy("http://localhost:8000"))
			v1.Any("/image/*path", createReverseProxy("http://localhost:8002"))
			v1.Any("/album/*path", createReverseProxy("http://localhost:8001"))
		}
		api := gw.Group("api")
		{
			v1 := api.Group("v1")
			v1.Any("/auth/*path", createReverseProxy("http://localhost:8003"))
		}
	}

	return router

}

func Start() {
	auth := middleware.NewAuth()
	router := InitRoutes(auth)
	err := router.Run(":8888")
	fmt.Println(err)
}

// creating reverse proxy for downstream microservices
func createReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		path := strings.Replace(c.Request.RequestURI, "/gw/api/", "", 1)
		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Path = path

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
