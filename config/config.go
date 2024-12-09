package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	if err := godotenv.Load(); err != nil {
        fmt.Printf("Error getting env, not comming through %v", err)
    }
	return os.Getenv(key)
}