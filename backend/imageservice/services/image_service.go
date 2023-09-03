package services

import (
	"imageservice/models"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type ImageService interface {
	FindAllImages(albumId string) (*[]models.ImageView, error)
	DeleteAllImages(albumId string) error
	GetImageMetadata(imageId string) (*models.ImageView, error)
	Upload(albumId string, files []*multipart.FileHeader, c *gin.Context) ([]string, error)
	Delete(imageId string) error
	// Upload(file) error
}
