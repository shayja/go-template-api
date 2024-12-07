// cmd/app/main.go
package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/middleware"
	"github.com/shayja/go-template-api/pkg/adapters/controllers"
	"github.com/shayja/go-template-api/pkg/frameworks/db"
)

type App struct {
	DB *sql.DB
	Router *gin.Engine
}

func (app *App) ConnectDB(){
	db := db.OpenDBConnection()
	app.DB = db
}

const (
    prefix string = "/api"
    api_ver uint8 = 1
  )

func (app *App) Routes() {
	router := gin.Default()
	fmt.Println(gin.Version)
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Set the api base url.
	baseUrl := fmt.Sprintf("%s/v%d/", prefix, api_ver)
	
	// Register user module
	userController := controllers.CreateUserController(app.DB)
	publicRoutes := router.Group(fmt.Sprintf("%s/auth", baseUrl))
	publicRoutes.POST("/register", userController.Register)
	publicRoutes.POST("/login", userController.Login)

	// Register product module
	productController := controllers.CreateProductController(app.DB)
	protectedRoutes := router.Group(fmt.Sprintf("%s/product", baseUrl))
	// Set Auth for the module routes
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

	// Set the product module routes.
	protectedRoutes.POST("", productController.Create)
	protectedRoutes.GET("", productController.GetAll)
	protectedRoutes.GET(":id", productController.GetSingle)
	protectedRoutes.PUT(":id", productController.Update)
	protectedRoutes.PATCH(":id", productController.UpdatePrice)
	protectedRoutes.POST("/image/:id", productController.UpdateImage)
	protectedRoutes.DELETE(":id", productController.Delete)

	// Register Routes
	app.Router = router
}


func (app *App) Run() {
	app.Router.Run(":8080")
}