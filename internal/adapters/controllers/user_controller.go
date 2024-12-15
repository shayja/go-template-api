// internal/adapters/controllers/user_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/utils"
)

type UserInteractor interface {
	GetUserById(id string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByMobile(mobile string) (*entities.User, error)
	ValidatePassword(user *entities.User, password string) error
	RegisterUser(request *entities.UserRequest) (*entities.User, error)
	GenerateAndSendOTP(mobile string) error
	VerifyOTP(mobile string, otp string) (*entities.User, error)
	ResendOTP(mobile string) error 
}


type UserController struct {
    UserInteractor UserInteractor
}

func (uc *UserController) Login(c *gin.Context) {
	var input *entities.AuthenticationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	AddRequestHeader(c)

	if len(input.Username) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Username is required"})
		return
	}

	if len(input.Password) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Password is required"})
		return
	}


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

	if len(userReq.Username) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Username is required"})
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


// RequestOTP handles generating and sending an OTP
func (uc *UserController) RequestOTP(c *gin.Context) {
	
	var inputReq entities.OtpRequest
	if err := c.ShouldBindJSON(&inputReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserInteractor.GenerateAndSendOTP(inputReq.Mobile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP handles OTP verification and login
func (uc *UserController) VerifyOTP(c *gin.Context) {
	
	var inputReq entities.VerifyOtpRequest
	if err := c.ShouldBindJSON(&inputReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := uc.UserInteractor.VerifyOTP(inputReq.Mobile, inputReq.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


func (uc *UserController) ResendOTP(c *gin.Context) {
	
	var inputReq entities.OtpRequest
	if err := c.ShouldBindJSON(&inputReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := uc.UserInteractor.ResendOTP(inputReq.Mobile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}