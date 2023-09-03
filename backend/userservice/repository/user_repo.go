package repository

import model "userservice/models"

type UserRepo interface {
	FindUserById(id int) (*model.UserView, error)
	FindUserByEmail(email string) (*model.UserView, error)
	FindUserByUsername(username string) (*model.UserView, error)
	FindUserByUsernameAndPassword(username string) (*model.UserLoginResponse, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
}
