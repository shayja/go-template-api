// internal/adapters/middleware/auth_middleware.go
package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTValidator func(context *gin.Context) error

func AuthRequired(validateJWT JWTValidator) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := validateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			fmt.Println(err)
			context.Abort()
			return
		}
		context.Next()
	}
}
