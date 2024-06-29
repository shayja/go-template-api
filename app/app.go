package app

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/controller"
	"github.com/shayja/go-template-api/db"
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
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	fmt.Println(gin.Version)
	controller := controller.CreateProductController(app.DB)
	baseUrl := fmt.Sprintf("%s/v%d/product", prefix, api_ver)
	r.POST(baseUrl, controller.Create)
	r.GET(baseUrl, controller.GetAll)
	r.GET(fmt.Sprintf("%s/:id", baseUrl), controller.GetSingle)
	r.PUT(fmt.Sprintf("%s/:id", baseUrl), controller.Update)
	r.PATCH(fmt.Sprintf("%s/:id", baseUrl), controller.UpdatePrice)
	r.DELETE(fmt.Sprintf("%s/:id", baseUrl), controller.Delete)
	app.Router = r
}

func (a *App) Run() {
	a.Router.Run(":8080")
}
