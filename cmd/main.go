// cmd/app/main.go
package main

import "github.com/shayja/go-template-api/cmd/app"


func main() {
	var app app.App

	// load database configuration and connection
	app.ConnectDB()

	// Set the routes
	app.Routes()

	// start the server
	app.Run()
}