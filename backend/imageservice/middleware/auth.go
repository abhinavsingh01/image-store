package middleware

import (
	"fmt"
	"imageservice/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	repo repository.AuthRepo
}

func NewAuth(repo repository.AuthRepo) *Auth {
	return &Auth{
		repo: repo,
	}
}

func (auth *Auth) AuthorizeImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := ""
		if c.Param("imageId") != "" {
			filename = c.Param("imageId")
		} else {
			filename = c.DefaultQuery("image", "")
		}
		if filename != "" && c.GetHeader("user-id") != "" {
			userId, _ := strconv.Atoi(c.GetHeader("user-id"))
			slice := strings.Split(filename, "_")
			filename = slice[0]
			exist, err := auth.repo.CheckImageAuth(userId, filename)
			if err == nil && exist {
				c.Next()
			} else {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user or image id"})
			c.Abort()
		}
	}
}

func (auth *Auth) AuthorizeAlbum() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("albumId") != "" && c.GetHeader("user-id") != "" {
			userId, _ := strconv.Atoi(c.GetHeader("user-id"))
			albumId := c.Param("albumId")
			exist, err := auth.repo.CheckAlbumAuth(userId, albumId)
			if err == nil && exist {
				c.Next()
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user or album id"})
			c.Abort()
		}
	}
}
