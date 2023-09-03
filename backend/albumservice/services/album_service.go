package services

import "albumservice/models"

type AlbumService interface {
	CreateNewAlbumForUser(userId int, album *models.AlbumRequest) (string, error)
	FindAlbumsByUserId(userId int) (*[]models.AlbumView, error)
	FindAlbumsById(userId int, albumId string) (*models.AlbumView, error)
	// UpdateAlbum(albumId int, album *models.AlbumRequest) error
	DeleteAlbum(userId int, albumId string) error
}
