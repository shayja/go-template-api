package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/controller"
	"github.com/shayja/go-template-api/db"
	"github.com/shayja/go-template-api/middleware"
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

	baseUrl := fmt.Sprintf("%s/v%d/", prefix, api_ver)
	productController := controller.CreateProductController(app.DB)
	userController := controller.CreateUserController(app.DB)


	publicRoutes := router.Group(baseUrl+"/auth")
	publicRoutes.POST("/register", userController.Register)
	publicRoutes.POST("/login", userController.Login)

	protectedRoutes := router.Group(baseUrl+"/product")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

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
