package server

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func IsDev() bool {
	return os.Getenv("IS_DEV") == "true"
}
