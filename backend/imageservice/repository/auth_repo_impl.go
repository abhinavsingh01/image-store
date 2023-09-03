package repository

import (
	"imageservice/models"

	"gorm.io/gorm"
)

type AuthRepoImpl struct {
	db *gorm.DB
}

func NewAuthRepoImpl(db *gorm.DB) AuthRepo {
	return &AuthRepoImpl{
		db: db,
	}
}

func (aRepo AuthRepoImpl) CheckAlbumAuth(userId int, albumId string) (bool, error) {
	var count int64
	err := aRepo.db.Model(&models.User{}).
		Joins("INNER JOIN albums ON users.id = albums.user_id").
		Where("users.id = ? AND albums.album_id = ?", userId, albumId).
		Count(&count).Error

	return count > 0, err
}

func (aRepo AuthRepoImpl) CheckImageAuth(userId int, imageId string) (bool, error) {
	var count int64
	err := aRepo.db.Model(&models.User{}).
		Joins("INNER JOIN albums ON users.id = albums.user_id").
		Joins("INNER JOIN images ON images.album_id = albums.id").
		Where("users.id = ? AND images.image_id = ?", userId, imageId).
		Count(&count).Error

	return count > 0, err
}
