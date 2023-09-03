package middleware

import (
	"albumservice/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	repo repository.AlbumRepo
}

func NewAuth(repo repository.AlbumRepo) *Auth {
	return &Auth{
		repo: repo,
	}
}

func (auth *Auth) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Param("id") != "" {
			userId, _ := strconv.Atoi(c.GetHeader("user-id"))
			albumId := c.Param("id")
			_, err := auth.repo.FinAlbumByUserAndAlbumId(userId, albumId)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}
