package controllers

import (
	"albumservice/models"
	"albumservice/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlbumController struct {
	albumService services.AlbumService
	logger       *zap.Logger
}

func NewAlbumController(albumService services.AlbumService, logger *zap.Logger) *AlbumController {
	return &AlbumController{
		albumService: albumService,
		logger:       logger,
	}
}

func (ac *AlbumController) GetAlbumsForUser(c *gin.Context) {
	response := *&models.Response{}
	userId := c.GetHeader("user-id")
	userIdInt, _ := strconv.Atoi(userId)
	albums, err := ac.albumService.FindAlbumsByUserId(userIdInt)
	ac.logger.Info("Getting all albums for user: " + userId)
	if err == nil {
		response.Data = albums
		if len(*albums) == 0 {
			response.Message = "No album found"
		} else {
			response.Message = "Albums found"
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		ac.logger.Error("Error Occurred: " + err.Error())
		response.Error = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
}

func (ac *AlbumController) GetAlbumById(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("id") != "" {
		userId := c.GetHeader("user-id")
		userIdInt, _ := strconv.Atoi(userId)
		albumId := c.Param("id")
		album, err := ac.albumService.FindAlbumsById(userIdInt, albumId)
		if err == nil {
			response.Data = album
			response.Message = "Albums found"
			c.JSON(http.StatusOK, response)
			return
		} else {
			ac.logger.Error("Error Occurred: " + err.Error())
			response.Error = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	response.Error = "Bad request"
	c.JSON(http.StatusBadRequest, response)
	c.Abort()
	return
}

func (ac *AlbumController) CreateNewAlbum(c *gin.Context) {
	response := *&models.Response{}
	userId := c.GetHeader("user-id")
	userIdInt, _ := strconv.Atoi(userId)

	ac.logger.Info("Creating new album for user: " + userId)

	var album models.AlbumRequest
	if err := c.BindJSON(&album); err != nil {
		ac.logger.Error("Error Occurred: " + err.Error())
		response.Error = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	albumId, err := ac.albumService.CreateNewAlbumForUser(userIdInt, &album)
	if err != nil {
		ac.logger.Error("Error Occurred: " + err.Error())
		response.Error = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Message = "Album created successfully"
	response.Data = gin.H{"albumId": albumId}
	c.JSON(http.StatusOK, response)
}

func (ac *AlbumController) DeleteAlbum(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("id") != "" {
		userId := c.GetHeader("user-id")
		userIdInt, _ := strconv.Atoi(userId)
		albumId := c.Param("id")

		ac.logger.Info("Deleting album: " + albumId)

		err := ac.albumService.DeleteAlbum(userIdInt, albumId)
		if err == nil {
			response.Message = "Album deleted successfully"
			c.JSON(http.StatusOK, response)
			return
		} else {
			ac.logger.Error("Error Occurred: " + err.Error())
			response.Error = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	response.Error = "Bad request"
	c.JSON(http.StatusBadRequest, response)
	c.Abort()
	return
}
