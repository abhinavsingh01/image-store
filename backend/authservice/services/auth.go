package services

import (
	config "authservice/configs"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	CreateJwt(userId int, username string) (string, error)
}

type AuthServiceImpl struct {
}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}

// generating jwt
func (authSvc *AuthServiceImpl) CreateJwt(userId int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
	})
	appConfig := config.GetConfig()
	return token.SignedString([]byte(appConfig.Secret))
}
