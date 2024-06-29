package controller

import "github.com/gin-gonic/gin"

type ProductControllerInterface interface {
	GetAll(*gin.Context)
	GetSingle(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}