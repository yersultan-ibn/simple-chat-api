package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}