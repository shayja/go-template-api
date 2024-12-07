// adapters/controllers/authentication_controller.go
package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/helper"
	"github.com/shayja/go-template-api/internal/utils"
	repositories "github.com/shayja/go-template-api/pkg/adapters/repositories/user"
	"github.com/shayja/go-template-api/pkg/entities"
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

	var userReq entities.User
	if err := c.ShouldBindJSON(&userReq); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}
	
	repositories := repositories.NewUserRepository(DB)

	insertedId, err := repositories.Create(userReq)
	
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
	var input entities.AuthenticationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	AddRequestHeader(c)
	DB := m.Db

	repositories := repositories.NewUserRepository(DB)

	user, err := repositories.GetByUsername(input.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	err = repositories.ValidatePassword(user, input.Password)

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