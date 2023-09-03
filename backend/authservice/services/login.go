package services

import (
	"authservice/clients"
	"authservice/models"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	GetUser(user *models.UserLogin) (*models.UserResponse, error)
}

type LoginServiceImpl struct {
	userClient clients.UserClient
	authSvc    AuthService
	logger     *zap.Logger
}

func NewLoginService(userClient clients.UserClient, authSvc AuthService, logger *zap.Logger) LoginService {
	return &LoginServiceImpl{
		userClient: userClient,
		authSvc:    authSvc,
		logger:     logger,
	}
}

// Get user from user service for login
func (l *LoginServiceImpl) GetUser(user *models.UserLogin) (*models.UserResponse, error) {
	loginRequest := &models.UserLoginRequest{
		Username: user.Username,
	}
	userDetails, err := l.userClient.GetUser(loginRequest)
	if err != nil {
		l.logger.Error("error occurred", zap.String("message", err.Error()))
		return nil, err
	}
	hash := fmt.Sprintf("%v", userDetails["password"])
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	if err != nil {
		l.logger.Error("error occurred", zap.String("message", err.Error()))
		return nil, errors.New("Wrong user name or password")
	}
	userId, _ := strconv.Atoi(fmt.Sprintf("%v", userDetails["Id"]))
	usernamr := fmt.Sprintf("%v", userDetails["username"])
	token, err := l.authSvc.CreateJwt(userId, usernamr)
	userResponse := &models.UserResponse{
		Token: token,
	}
	return userResponse, err
}
