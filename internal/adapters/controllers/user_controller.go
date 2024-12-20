package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/common/helpers"
	"github.com/shayja/go-template-api/internal/entities"
	appErrors "github.com/shayja/go-template-api/internal/errors"
	"github.com/shayja/go-template-api/internal/utils"
)

type UserInteractor interface {
	GetUserById(id string) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByMobile(mobile string) (*entities.User, error)
	ValidatePassword(passwordHash string, plainPassword string) error
	RegisterUser(request *entities.UserRequest) (*entities.User, error)
	GenerateAndSendOTP(mobile string) error
	VerifyOTP(mobile string, otp string) (*entities.User, error)
	ResendOTP(mobile string) error 
}

type UserController struct {
    UserInteractor UserInteractor
}

// @Summary Login User
// @Description Authenticate a user with username and password
// @Tags Users
// @Accept json
// @Produce json
// @Param input body entities.AuthenticationInput true "Authentication Input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	AddRequestHeader(c)

	var input *entities.AuthenticationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

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

	err = uc.UserInteractor.ValidatePassword(user.Password, input.Password)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": err})
		return
	}

	jwt, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}

// @Summary Register User
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body entities.UserRequest true "User Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/register [post]
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

// @Summary Request OTP
// @Description Generate and send an OTP to a user's mobile number
// @Tags Users
// @Accept json
// @Produce json
// @Param input body entities.OtpRequest true "OTP Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/request-otp [post]
func (uc *UserController) RequestOTP(c *gin.Context) {
	var inputReq entities.OtpRequest
	if err := c.ShouldBindJSON(&inputReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mobile number"})
		return
	}

	mobile, errBadRequest :=  helpers.ConvertToMobile(inputReq.Mobile)
	if errBadRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": mobile})
		return
	}
	inputReq.Mobile = mobile

	err := uc.UserInteractor.GenerateAndSendOTP(inputReq.Mobile)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": appErrors.ErrUserNotFound.Message })
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// @Summary Verify OTP
// @Description Verify the OTP and authenticate the user
// @Tags Users
// @Accept json
// @Produce json
// @Param input body entities.VerifyOtpRequest true "Verify OTP Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/verify-otp [post]
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

// @Summary Resend OTP
// @Description Resend the OTP to a user's mobile number
// @Tags Users
// @Accept json
// @Produce json
// @Param input body entities.OtpRequest true "OTP Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/resend-otp [post]
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
