// cmd/app/main.go
package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/middleware"
	"github.com/shayja/go-template-api/pkg/adapters/controllers"
	productrepo "github.com/shayja/go-template-api/pkg/adapters/repositories/product"
	userrepo "github.com/shayja/go-template-api/pkg/adapters/repositories/user"
	"github.com/shayja/go-template-api/pkg/frameworks/db"
	"github.com/shayja/go-template-api/pkg/usecases"
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
	
	// Register the User module
	userRepo := &userrepo.UserRepository{Db: app.DB}
    userInteractor := usecases.UserInteractor{UserRepository: userRepo}
    userController := controllers.UserController{UserInteractor: userInteractor}
	// Configure User Routes
	publicRoutes := router.Group(fmt.Sprintf("%s/auth", baseUrl))
	publicRoutes.POST("/register", userController.RegisterUser)
	publicRoutes.POST("/login", userController.Login)
	//publicRoutes.POST("/verify_otp", userController.)
	//publicRoutes.POST("/resend_otp", userController.)
	//publicRoutes.POST("/me", userController.)

	// Register the Product module
	productRepo := &productrepo.ProductRepository{Db: app.DB}
	productInteractor := usecases.ProductInteractor{ProductRepository: productRepo}
	productController := controllers.ProductController{ProductInteractor: productInteractor}

	// Configure Product Routes
	protectedRoutes := router.Group(fmt.Sprintf("%s/product", baseUrl))
	// Set Auth for the module routes
	protectedRoutes.Use(middleware.AuthRequired())

	// Set the product module routes.
	protectedRoutes.POST("", productController.Create)
	protectedRoutes.GET("", productController.GetAll)
	protectedRoutes.GET(":id", productController.GetById)
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
