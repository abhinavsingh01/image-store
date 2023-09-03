package controllers

import (
	"authservice/models"
	"authservice/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	registerSvc services.Register
	loginSvc    services.LoginService
	logger      *zap.Logger
}

func NewAuthController(registerSvc services.Register,
	loginSvc services.LoginService, logger *zap.Logger) *AuthController {
	return &AuthController{
		registerSvc: registerSvc,
		loginSvc:    loginSvc,
		logger:      logger,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	response := *&models.Response{}
	var userLogin models.UserLogin

	if err := c.BindJSON(&userLogin); err != nil {
		a.logger.Error("error in login payload", zap.String("message", err.Error()))
		response.Error = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := a.loginSvc.GetUser(&userLogin)

	if err != nil {
		a.logger.Error("error in login payload", zap.String("message", err.Error()))
		response.Error = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Message = "Login successfully"
	response.Data = user

	c.JSON(http.StatusOK, response)
	return
}

func (a *AuthController) Register(c *gin.Context) {
	var user models.UserRequest
	if err := c.BindJSON(&user); err != nil {
		a.logger.Error("error in register payload", zap.String("message", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.registerSvc.RegisterUser(&user); err != nil {
		a.logger.Error("error in register payload", zap.String("message", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, a)
}
