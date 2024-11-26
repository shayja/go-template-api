package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func OpenDBConnection() *sql.DB {

    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
    

    host := os.Getenv("DB_HOST")
    if host == "" {
        log.Fatal("DB_HOST not set")
    }
    log.Println("DB_HOST is set to", host)

    port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
  
    // Format the connecttion string to the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    
    // Establish a connection to the PostgreSQL database
    db, err := sql.Open("postgres", psqlInfo)
    
     // checks if we are connected to db
    if err != nil {
        log.Fatal("Error connecting to database :", err)
        panic(err)
    }

    fmt.Printf("The database `%s` is successfully connected!", dbname)
    return db
}