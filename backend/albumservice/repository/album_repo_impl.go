package repository

import (
	model "albumservice/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumRepoImpl struct {
	db *gorm.DB
}

func NewAlbumRepoImpl(db *gorm.DB) AlbumRepo {
	return &AlbumRepoImpl{
		db: db,
	}
}

func (repo *AlbumRepoImpl) FindAllAlbumByUserId(userId int) (*[]model.AlbumView, error) {
	var albumView []model.AlbumView
	var album model.Album
	result := repo.db.Model(&album).Select("AlbumId", "AlbumName", "AlbumDescription").Where("user_id = ?", userId).Find(&albumView)
	return &albumView, result.Error
}

func (repo *AlbumRepoImpl) Create(userId int, album *model.Album) (string, error) {
	albumId := uuid.New().String()
	album.AlbumId = albumId
	result := repo.db.Create(album)
	if result.Error != nil {
		return "", result.Error
	}
	return albumId, nil
}

func (repo *AlbumRepoImpl) Delete(userId int, albumId string) error {
	result := repo.db.Where("user_id = ?", userId).Where("album_id = ?", albumId).Delete(&model.Album{})
	if result.RowsAffected == 0 {
		return errors.New("Album not found")
	}
	return result.Error
}

func (repo *AlbumRepoImpl) FinAlbumByUserAndAlbumId(userId int, albumId string) (*model.AlbumView, error) {
	var albumView model.AlbumView
	var album model.Album
	result := repo.db.Model(&album).Select("AlbumId", "AlbumName", "AlbumDescription").
		Where("user_id = ?", userId).Where("album_id = ?", albumId).First(&albumView)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("Album not found")
	}
	return &albumView, nil
}
