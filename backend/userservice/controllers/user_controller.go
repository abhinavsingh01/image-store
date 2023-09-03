package controllers

import (
	"net/http"
	"strconv"
	"userservice/models"
	"userservice/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userService services.UserService
	logger      *zap.Logger
}

func NewUserController(userService services.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (u *UserController) GetUserById(c *gin.Context) {
	response := *&models.Response{}
	if c.Param("id") != "" {
		id, _ := strconv.Atoi(c.Param("id"))
		user, err := u.userService.FindUserById(id)
		if err == nil {
			response.Message = "User found"
			response.Data = user
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Error = "User not found"
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	response.Error = "Bad Request"
	c.JSON(http.StatusBadRequest, response)
	c.Abort()
	return
}

func (u *UserController) GetUserDetails(c *gin.Context) {
	response := *&models.Response{}

	id, _ := strconv.Atoi(c.GetHeader("user-id"))
	user, err := u.userService.FindUserById(id)
	if err == nil {
		response.Message = "User found"
		response.Data = user
		c.JSON(http.StatusOK, response)
		return
	} else {
		response.Error = "User not found"
		c.JSON(http.StatusBadRequest, response)
		return
	}

}

func (u *UserController) CreateNewUser(c *gin.Context) {
	var user models.User
	response := *&models.Response{}
	if err := c.BindJSON(&user); err != nil {
		response.Error = "Bad Request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := u.userService.CreateNewUser(&user); err != nil {
		response.Error = "User already exist"
		u.logger.Error("User already exist")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (u *UserController) FindUser(c *gin.Context) {
	var user models.UserLogin
	response := *&models.Response{}
	if err := c.BindJSON(&user); err != nil {
		response.Error = "Bad Request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userView, err := u.userService.FindUser(user.Username)

	if err != nil {
		response.Error = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Data = userView
	c.JSON(http.StatusOK, response)
}
