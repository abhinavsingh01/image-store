package controllers

import (
	"imageservice/models"
	"imageservice/services"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Image struct {
	imageService services.ImageService
	logger       *zap.Logger
}

func NewImage(imageService services.ImageService, logger *zap.Logger) *Image {
	return &Image{
		imageService: imageService,
		logger:       logger,
	}
}

// Getting all images in an album
func (ctrl *Image) GetAllImagesOfAlbum(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("albumId") != "" {
		albumId := c.Param("albumId")

		ctrl.logger.Info("Getting all images from album for " + albumId)

		images, err := ctrl.imageService.FindAllImages(albumId)
		if err == nil {
			response.Data = images
			if len(*images) == 0 {
				response.Message = "No images found"
			} else {
				response.Message = "Images found"
			}
			c.JSON(http.StatusOK, response)
			return
		} else {
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			response.Error = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}
	} else {
		ctrl.logger.Error("Error occurred", zap.String("message", "AlbumId is missing"))
		response.Error = "AlbumId is missing"
		c.JSON(http.StatusBadRequest, response)
		return
	}
}

// Get image by id
func (ctrl *Image) GetAllImagesById(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("id") != "" {
		imageId := c.Param("id")
		ctrl.logger.Info("Getting image for id: " + imageId)
		image, err := ctrl.imageService.GetImageMetadata(imageId)
		if err == nil {
			response.Message = "Image found"
			response.Data = image
			c.JSON(http.StatusOK, response)
			return
		} else {
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			response.Error = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	response.Error = "bad request"
	c.JSON(http.StatusBadRequest, response)
	c.Abort()
	return
}

// Delete image by id
func (ctrl *Image) DeleteImage(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("imageId") != "" {
		imageId := c.Param("imageId")
		ctrl.logger.Info("Deleting image for id: " + imageId)
		err := ctrl.imageService.Delete(imageId)
		if err == nil {
			response.Message = "Image deleted"
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Error = "Not able to deelete image"
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	response.Error = "bad request"
	c.JSON(http.StatusBadRequest, response)
	c.Abort()
	return
}

// Upload image in an album
func (ctrl *Image) UploadImage(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("albumId") != "" {
		albumId := c.Param("albumId")
		ctrl.logger.Info("Uploading image in album: " + albumId)
		form, _ := c.MultipartForm()
		files := form.File["files"]
		if len(files) == 0 {
			response.Error = "No images to upload"
			c.JSON(http.StatusBadRequest, response)
			return
		}
		ids, err := ctrl.imageService.Upload(albumId, files, c)
		if err != nil {
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			response.Error = "Failed to save image"
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response.Message = "Images uploaded successfully"
		response.Data = ids
		c.JSON(http.StatusOK, response)
	} else {
		response.Error = "Album not found"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
}

// Download image by id
func (ctrl *Image) Download(c *gin.Context) {
	if c.Param("imageId") != "" {
		filename := c.Param("imageId")
		ctrl.logger.Info("Getting image: " + filename)
		filePath := filepath.Join("./uploads", filename)
		_, err := os.Stat(filePath)
		if err != nil {
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		file, err := os.Open(filePath)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=file")
		http.ServeContent(c.Writer, c.Request, "file", time.Now(), file)
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
	return
}

// Delete all images in an album
func (ctrl *Image) DeleteAllImages(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("albumId") != "" {
		albumId := c.Param("albumId")
		ctrl.logger.Info("Deleting all images in album: " + albumId)
		err := ctrl.imageService.DeleteAllImages(albumId)
		if err == nil {
			response.Message = "Images deleted successfully"
			c.JSON(http.StatusOK, response)
			return
		} else {
			ctrl.logger.Error("Error occurred", zap.String("message", err.Error()))
			response.Error = "failed to delete images"
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
}
