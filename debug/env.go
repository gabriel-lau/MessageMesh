package debug

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		Log("error", "Error loading .env file")
	}
	return os.Getenv(key)
}
