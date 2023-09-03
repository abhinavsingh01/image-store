package services

import model "userservice/models"

type UserService interface {
	CreateNewUser(*model.User) error
	FindUserById(id int) (*model.UserView, error)
	FindUser(username string) (*model.UserLoginResponse, error)
}
