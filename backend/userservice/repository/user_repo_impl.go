package repository

import (
	"errors"
	model "userservice/models"

	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (repo *UserRepoImpl) FindUserById(id int) (*model.UserView, error) {
	var userView model.UserView
	var user model.User
	result := repo.db.Model(&user).Select("Id", "name", "username", "email").First(&userView, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &userView, result.Error
}

func (repo *UserRepoImpl) FindUserByEmail(email string) (*model.UserView, error) {
	var user model.UserView
	result := repo.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (repo *UserRepoImpl) FindUserByUsername(username string) (*model.UserView, error) {
	var user model.UserView
	result := repo.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (repo *UserRepoImpl) FindUserByUsernameAndPassword(username string) (*model.UserLoginResponse, error) {
	var userLoginResponse model.UserLoginResponse
	var user model.User
	result := repo.db.Model(&user).Where("username = ?", username).First(&userLoginResponse)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &userLoginResponse, result.Error
}

func (repo *UserRepoImpl) Create(user *model.User) error {
	result := repo.db.Create(user)
	return result.Error
}

func (repo *UserRepoImpl) Update(user *model.User) error {
	result := repo.db.Save(user)
	return result.Error
}

func (repo *UserRepoImpl) Delete(id int) error {
	result := repo.db.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}
