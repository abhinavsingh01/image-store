package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Authorize() gin.HandlerFunc
}

type AuthImpl struct {
}

func NewAuth() Auth {
	return &AuthImpl{}
}

func (auth *AuthImpl) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}
		parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("image-store-secret-key"), nil
		})
		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, _ := parsedToken.Claims.(jwt.MapClaims)
		userId := fmt.Sprintf("%v", claims["user_id"])
		c.Request.Header.Add("user-id", userId)

		c.Next()
	}
}
