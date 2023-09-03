package services

import (
	"fmt"
	model "userservice/models"
	"userservice/repository"
)

type UserServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserServiceImpl(userRepo repository.UserRepo) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (usvc *UserServiceImpl) CreateNewUser(user *model.User) error {
	err := usvc.userRepo.Create(user)
	if err != nil {
		fmt.Println("error creating new user")
	}
	return err
}

func (usvc *UserServiceImpl) FindUserById(id int) (*model.UserView, error) {
	result, err := usvc.userRepo.FindUserById(id)
	if err != nil {
		fmt.Println("error getting user")
		return nil, err
	}
	return result, nil
}

func (usvc *UserServiceImpl) FindUser(username string) (*model.UserLoginResponse, error) {
	result, err := usvc.userRepo.FindUserByUsernameAndPassword(username)
	if err != nil {
		fmt.Println("error getting user")
		return nil, err
	}
	return result, nil
}
