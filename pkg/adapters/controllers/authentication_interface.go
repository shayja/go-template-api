// adapters/controllers/authentication_interface.go
package controllers

import "github.com/gin-gonic/gin"

type AuthenticationInterface interface {
	Login(*gin.Context)
	Register(*gin.Context)

}