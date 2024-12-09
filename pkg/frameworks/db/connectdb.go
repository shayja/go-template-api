// frameworks/db/connectdb.go
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shayja/go-template-api/config"
)

func OpenDBConnection() (*sql.DB) {
    
    // Read the connection propertied from the env variables.
    v := NewDbInfo(
        config.Config("DB_HOST"),
        config.Config("DB_PORT"),
        config.Config("DB_USER"),
        config.Config("DB_PASSWORD"),
        config.Config("DB_NAME"),
        config.Config("SSL_MODE"),
    )

    // Format the connection string to the database
    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", v.Host, v.Port, v.User, v.Password, v.DBName, v.SSLMode)
    
    // Establish a connection to the PostgreSQL database
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        fmt.Print("Error connecting to database:", err)
        panic(err)
    }

    // checks if we are connected to db
    err = db.Ping()
    if err != nil {
        fmt.Print("Ping database error:", err)
        panic(err)
    }

    fmt.Printf("Successfully connected to the database `%s`", v.DBName)
    return db
}