package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/helper"
	"github.com/shayja/go-template-api/model"
	repository "github.com/shayja/go-template-api/repository/user"
	"github.com/shayja/go-template-api/utils"
)

type AuthenticationController struct {
	Db *sql.DB
}

func CreateUserController(db *sql.DB) AuthenticationInterface {
	return &AuthenticationController{Db: db}
}

func (m *AuthenticationController) Register(c *gin.Context) {
	
	AddRequestHeader(c)
	DB := m.Db

	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	
	repository := repository.NewUserRepository(DB)

	insertedId, err := repository.Create(userReq)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if utils.IsValidUUID(insertedId) {
		c.JSON(http.StatusCreated, gin.H{"status": "success", "msg": nil, "id": insertedId})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": "insert product failed"})
	}
}

func (m *AuthenticationController) Login(c *gin.Context) {
	var input model.AuthenticateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	AddRequestHeader(c)
	DB := m.Db

	repository := repository.NewUserRepository(DB)

	user, err := repository.GetByUsername(input.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	err = repository.ValidatePassword(user, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
} 