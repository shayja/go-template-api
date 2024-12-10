// internal/adapters/controllers/user_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/shayja/go-template-api/internal/utils"
)

type UserController struct {
    UserInteractor usecases.UserInteractor
}

func (uc *UserController) Login(c *gin.Context) {
	var input *entities.AuthenticationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	AddRequestHeader(c)

	user, err := uc.UserInteractor.GetUserByUsername(input.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	err = uc.UserInteractor.ValidatePassword(user, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	jwt, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
} 

func (uc *UserController) RegisterUser(c *gin.Context) {
	
	AddRequestHeader(c)

	var userReq entities.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

    user, err := uc.UserInteractor.RegisterUser(&userReq)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if user!=nil {
		c.JSON(http.StatusCreated, gin.H{"status": "success", "msg": nil, "id": user.Id})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "insert product failed"})
	}
}