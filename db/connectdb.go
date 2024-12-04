package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbInfo struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func OpenDBConnection() (*sql.DB) {

    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error getting env, not comming through %v", err)
    }
    
    v := DbInfo{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  os.Getenv("SSL_MODE"),
    }

    // Format the connection string to the database
    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", v.Host, v.Port, v.User, v.Password, v.DBName, v.SSLMode)
    
    // Establish a connection to the PostgreSQL database
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
        panic(err)
    }

    // checks if we are connected to db
    err = db.Ping()
    if err != nil {
        log.Fatal("Ping database error:", err)
        panic(err)
    }

    fmt.Printf("Successfully connected to the database `%s`", v.DBName)
    return db
}