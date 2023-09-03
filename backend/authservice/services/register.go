package services

import (
	"authservice/clients"
	"authservice/models"

	"golang.org/x/crypto/bcrypt"
)

type Register interface {
	RegisterUser(user *models.UserRequest) error
}

type RegisterImpl struct {
	userClient clients.UserClient
}

func NewResgister(userClient clients.UserClient) Register {
	return &RegisterImpl{
		userClient: userClient,
	}
}

func (r *RegisterImpl) RegisterUser(user *models.UserRequest) error {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(cryptedPassword)
	err = r.userClient.Register(user)
	return err
}
