// internal/adapters/middleware/auth_middleware_test.go
package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/adapters/middleware"
	"github.com/stretchr/testify/assert"
)

func mockValidateJWTSuccess(context *gin.Context) error {
	return nil
}

func mockValidateJWTError(context *gin.Context) error {
	return errors.New("invalid token")
}

func TestAuthRequired_Success(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.AuthRequired(mockValidateJWTSuccess))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
}

func TestAuthRequired_Unauthorized(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.AuthRequired(mockValidateJWTError))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error": "Authentication required"}`, w.Body.String())
}
