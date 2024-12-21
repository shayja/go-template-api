// cmd/app/main.go
package main

import "github.com/shayja/go-template-api/cmd/app"

// Swagger
//
//  @title                      Go Template API
//  @version                    1.0
//  @description                API documentation for the Go Template API.
//  @contact.name               Shay Jacoby
//  @contact.url                https://github.com/shayja/
//  @contact.email              shayja@gmail.com
//  @license.name               Apache 2.0
//  @license.url                http://www.apache.org/licenses/LICENSE-2.0.html
//  @host                       localhost:8080
//  @BasePath                   /api/v1
//  @schemes                    http https
//  @securityDefinitions.apiKey apiKey
//  @in                         header
//  @name                       Authorization
//  @description                Type "Bearer" followed by a space and JWT token. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	var app app.App

	// load database configuration and connection
	app.ConnectDB()

	// Set the routes
	app.Routes()

	// start the server
	app.Run()
}