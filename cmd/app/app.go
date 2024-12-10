// cmd/app/app.go
package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/config"
	"github.com/shayja/go-template-api/internal/adapters/controllers"
	"github.com/shayja/go-template-api/internal/adapters/middleware"

	productrepo "github.com/shayja/go-template-api/internal/adapters/repositories/product"
	userrepo "github.com/shayja/go-template-api/internal/adapters/repositories/user"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/shayja/go-template-api/pkg/constants"
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

func (app *App) Routes() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Set the api base url.
	baseUrl := fmt.Sprintf("%s/v%d/", constants.ApiPrefix, constants.ApiVersion)
	
	// Register the User module
	userRepo := &userrepo.UserRepository{Db: app.DB}
	userInteractor := &usecases.UserInteractor{UserRepository: userRepo}
	userController := controllers.UserController{UserInteractor: userInteractor}

	// Configure User Routes
	publicRoutes := router.Group(fmt.Sprintf("%s/auth", baseUrl))
	publicRoutes.POST("/register", userController.RegisterUser)
	publicRoutes.POST("/login", userController.Login)
	//publicRoutes.POST("/verify_otp", userController.VerifyOTP)
	//publicRoutes.POST("/resend_otp", userController.ResendOTP)

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
	app.Router.Run(fmt.Sprintf(`:%s`, config.Config("SERVER_PORT")))
}
