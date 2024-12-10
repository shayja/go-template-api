// adapters/middleware/auth_middleware.go
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := utils.ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			fmt.Println(err)
			context.Abort()
			return
		}
		context.Next()
	}
}