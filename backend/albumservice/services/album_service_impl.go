package services

import (
	"albumservice/client"
	"albumservice/models"
	"albumservice/repository"
	"errors"

	"go.uber.org/zap"
)

type AlbumServiceImpl struct {
	repo        repository.AlbumRepo
	imageClient client.ImageClient
	logger      *zap.Logger
}

func NewAlbumServiceImpl(repo repository.AlbumRepo,
	imageClient client.ImageClient, logger *zap.Logger) AlbumService {
	return &AlbumServiceImpl{
		repo:        repo,
		imageClient: imageClient,
		logger:      logger,
	}
}

func (albumSvc *AlbumServiceImpl) CreateNewAlbumForUser(userId int, album *models.AlbumRequest) (string, error) {
	newAlbum := models.Album{UserId: userId, AlbumName: album.AlbumName, AlbumDescription: album.AlbumDesc}
	albumId, err := albumSvc.repo.Create(userId, &newAlbum)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return "", errors.New("Not able to create album")
	}
	return albumId, nil
}

func (albumSvc *AlbumServiceImpl) DeleteAlbum(userId int, albumId string) error {
	album, err := albumSvc.repo.FinAlbumByUserAndAlbumId(userId, albumId)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return errors.New("Not able to find album")
	}
	err = albumSvc.imageClient.DeleteImage(userId, album.AlbumId)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return errors.New("Not able to delete images of album")
	}
	err = albumSvc.repo.Delete(userId, albumId)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return errors.New("Not able to delete album")
	}
	return nil
}

func (albumSvc *AlbumServiceImpl) FindAlbumsByUserId(userId int) (*[]models.AlbumView, error) {
	albums, err := albumSvc.repo.FindAllAlbumByUserId(userId)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return nil, errors.New("Not able to find album")
	}
	return albums, nil
}

func (albumSvc *AlbumServiceImpl) FindAlbumsById(userId int, albumId string) (*models.AlbumView, error) {
	album, err := albumSvc.repo.FinAlbumByUserAndAlbumId(userId, albumId)
	if err != nil {
		albumSvc.logger.Error("Error occurred " + err.Error())
		return nil, errors.New("Not able to find album")
	}
	return album, nil
}
