package services

import (
	"errors"
	"imageservice/models"
	"imageservice/repository"
	"imageservice/utils"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageServiceImpl struct {
	repo      repository.ImageRepo
	imageUtil *utils.ImageUtils
	logger    *zap.Logger
}

func NewImageServiceImpl(repo repository.ImageRepo,
	imageUtil *utils.ImageUtils, logger *zap.Logger) ImageService {
	return &ImageServiceImpl{
		repo:      repo,
		imageUtil: imageUtil,
		logger:    logger,
	}
}

func (svc *ImageServiceImpl) FindAllImages(albumId string) (*[]models.ImageView, error) {
	images, err := svc.repo.FindAllImages(albumId)
	if err != nil {
		svc.logger.Error("error fetching image" + err.Error())
	}
	return images, err
}

func (svc *ImageServiceImpl) GetImageMetadata(imageId string) (*models.ImageView, error) {
	image, err := svc.repo.GetImageMetadata(imageId)
	if err != nil {
		svc.logger.Error("error fetching image" + err.Error())
	}
	return image, err
}

func (svc *ImageServiceImpl) Upload(albumId string, files []*multipart.FileHeader, c *gin.Context) ([]string, error) {
	const MaxUploadSize = 25 * 1024 * 1024
	for _, file := range files {
		if file.Size > MaxUploadSize {
			svc.logger.Info("File is too large")
			return nil, errors.New("File is too Large. Maximum allowed is 25 MB")
		}
		ext := filepath.Ext(file.Filename)
		if strings.ToLower(ext) != ".jpg" &&
			strings.ToLower(ext) != ".jpeg" &&
			strings.ToLower(ext) != ".png" {
			svc.logger.Info("File is not image")
			return nil, errors.New("Wrong file type")
		}
	}
	var ids []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		newId := uuid.New()
		filename := newId.String()
		destination := filepath.Join("./uploads", filename)
		if err := c.SaveUploadedFile(file, destination); err != nil {
			return nil, err
		}
		newImage := models.Image{ImageId: filename, ImageName: filename, ImageUrl: filename}
		err := svc.repo.Create(albumId, &newImage)
		// Runnig go routine to resize image to different sizes and save.
		svc.imageUtil.Resize(filename, ext)
		if err != nil {
			return nil, err
		}
		ids = append(ids, filename)
	}
	return ids, nil
}

func (svc *ImageServiceImpl) Delete(imageId string) error {

	svc.imageUtil.RemoveImageFile(imageId)
	_, err := svc.repo.Delete(imageId)
	if err != nil {
		svc.logger.Error("error deleting image" + err.Error())
	}
	return err
}

func (svc *ImageServiceImpl) DeleteAllImages(albumId string) error {
	images, err := svc.repo.FindAllImages(albumId)
	if err != nil {
		svc.logger.Error("error deleting image" + err.Error())
		return err
	}
	var wgm sync.WaitGroup
	wgm.Add(len(*images))
	for _, image := range *images {

		// Runnig go routine to remove all images inside an album.
		go func(imageView models.ImageView) {
			defer wgm.Done()
			svc.imageUtil.RemoveImageFile(imageView.ImageId)
			err = svc.Delete(imageView.ImageId)
		}(image)
	}
	wgm.Wait()
	return nil
}
