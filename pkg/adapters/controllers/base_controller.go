// adapters/controllers/base_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddRequestHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}
