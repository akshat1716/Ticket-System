package controllers

import (
	"net/http"

	"ticket-system/services"
	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input services.RegisterInput
	if !utils.BindJSON(c, &input) {
		return
	}

	user, err := ctrl.authService.Register(input)
	if err != nil {
		if err.Error() == "Email already registered" {
			utils.JSONError(c, http.StatusConflict, err.Error())
			return
		}
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input services.LoginInput
	if !utils.BindJSON(c, &input) {
		return
	}

	token, err := ctrl.authService.Login(input)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Invalid email or password" {
			status = http.StatusUnauthorized
		}
		utils.JSONError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
