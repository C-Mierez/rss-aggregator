package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Create a map of all environment variables
type EnvVariables string

// @dev Add all environment variables here and update the LoadAndCheck function
const (
	PORT         EnvVariables = "PORT"
	DATABASE_URL EnvVariables = "DATABASE_URL"
)

func LoadAndCheck(shouldLog bool) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found. Error: %s\n", err.Error())
	}

	// Check if all environment variables are set
	for _, key := range []EnvVariables{PORT, DATABASE_URL} {
		_, ok := os.LookupEnv(string(key))
		if !ok {
			log.Fatalf("Environment variable %s is not set.\n", key)
		} else {
			if shouldLog {
				log.Printf("Environment variable %s is set to %s\n", key, os.Getenv(string(key)))
			}
		}
	}
}

func Get(key EnvVariables) string {
	return os.Getenv(string(key))
}
