// cmd/app/app.go
package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/config"
	"github.com/shayja/go-template-api/docs"
	"github.com/shayja/go-template-api/internal/adapters/controllers"
	"github.com/shayja/go-template-api/internal/adapters/middleware"
	"github.com/shayja/go-template-api/internal/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	repositories "github.com/shayja/go-template-api/internal/adapters/repositories/order"
	productrepo "github.com/shayja/go-template-api/internal/adapters/repositories/product"
	userrepo "github.com/shayja/go-template-api/internal/adapters/repositories/user"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/shayja/go-template-api/pkg/constants"
	sql_postgres "github.com/shayja/go-template-api/pkg/drivers/sql"
)

type App struct {
	DB *sql.DB
	Router *gin.Engine
}

func (app *App) ConnectDB(){
	db := sql_postgres.OpenDBConnection()
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
	publicRoutes.POST("/send_otp", userController.RequestOTP)
	publicRoutes.POST("/verify_otp", userController.VerifyOTP)
	publicRoutes.POST("/resend_otp", userController.ResendOTP)

	// Register the Product module
	productRepo := &productrepo.ProductRepository{Db: app.DB}
	productInteractor := usecases.ProductInteractor{ProductRepository: productRepo}
	productController := controllers.ProductController{ProductInteractor: productInteractor}

	// Configure Product Routes
	protectedRoutes := router.Group(fmt.Sprintf("%s/product", baseUrl))
	// Set Auth for the module routes
	protectedRoutes.Use(middleware.AuthRequired(utils.ValidateJWT))

	// Set the product module routes.
	protectedRoutes.POST("", productController.Create)
	protectedRoutes.GET("", productController.GetAll)
	protectedRoutes.GET(":id", productController.GetById)
	protectedRoutes.PUT(":id", productController.Update)
	protectedRoutes.PATCH(":id", productController.UpdatePrice)
	protectedRoutes.POST("/image/:id", productController.UpdateImage)
	protectedRoutes.DELETE(":id", productController.Delete)


	// Register the Order module
	orderRepo := &repositories.OrderRepository{Db: app.DB}
	orderUsecase := &usecases.OrderUsecase{OrderRepo: orderRepo}
	orderController := &controllers.OrderController{OrderUsecase: orderUsecase}

	// Configure Order Routes
	orderRoutes := router.Group(fmt.Sprintf("%s/order", baseUrl))
	orderRoutes.Use(middleware.AuthRequired(utils.ValidateJWT))

	// Set the order module routes.
	orderRoutes.POST("", orderController.Create)
	orderRoutes.GET("", orderController.GetOrders)
	orderRoutes.GET(":id", orderController.GetById)
	orderRoutes.PUT(":id/status", orderController.UpdateStatus)



	// Swagger setup
	docs.SwaggerInfo.Title = "Go Template API"
	docs.SwaggerInfo.Description = "API documentation for the Go Template API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	// Register Routes
	app.Router = router
}


func (app *App) Run() {
	app.Router.Run(fmt.Sprintf(`:%s`, config.Config("SERVER_PORT")))
}
