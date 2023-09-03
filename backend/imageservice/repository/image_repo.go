package repository

import "imageservice/models"

type ImageRepo interface {
	FindAllImages(albumId string) (*[]models.ImageView, error)
	GetImageMetadata(imageId string) (*models.ImageView, error)
	Create(albumId string, image *models.Image) error
	Delete(imageId string) (*models.Image, error)
}
