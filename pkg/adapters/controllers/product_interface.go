// adapters/controllers/product_interface.go
package controllers

import "github.com/gin-gonic/gin"

type ProductControllerInterface interface {
	GetAll(*gin.Context)
	GetSingle(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	UpdatePrice(c *gin.Context)
	UpdateImage(c *gin.Context)
	Delete(*gin.Context)
}