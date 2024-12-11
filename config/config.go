package config

import (
	"fmt"
	"os"
)

// Config func to get env value
// func Config(key string) string {
// 	// load .env file
// 	if err := godotenv.Load(); err != nil {
//         fmt.Printf("Error getting env, not comming through %v", err)
//     }
// 	return os.Getenv(key)
// }

func Config(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Print()
		return value
	}
	return ""
}