package repository

import model "albumservice/models"

type AlbumRepo interface {
	//FindAlbumById(userId int, id int) (*model.AlbumView, error)
	Create(userId int, album *model.Album) (string, error)
	//Update(albumId int, user *model.Album) error
	Delete(userId int, albumId string) error
	FindAllAlbumByUserId(userId int) (*[]model.AlbumView, error)
	FinAlbumByUserAndAlbumId(userId int, albumId string) (*model.AlbumView, error)
}
