package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadAndCheckENV(shouldLog bool) {
	// Load .env file
	godotenv.Load()

	// Create a map of all environment variables
	envVars := map[string]string{
		"PORT": os.Getenv("PORT"),
	}

	// Check if any of the environment variables are empty
	for key, value := range envVars {
		if value == "" {
			log.Fatalf("Environment variable %s is not set", key)
		}
		if shouldLog {
			log.Printf("ENV Variable loaded: %s=%s\n", key, value)
		}
	}
}
