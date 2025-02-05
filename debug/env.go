package debug

import (
	"os"

	"github.com/joho/godotenv"
)

// get HEADLESS environment variable
var IsHeadless = GetEnvVar("HEADLESS") == "true"

var Username = GetEnvVar("USERNAME")

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		Log("error", "Error loading .env file")
	}
	return os.Getenv(key)
}
