package configs

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

var (
	GIN_MODE = os.Getenv("GIN_MODE")
	BE_HOST  = os.Getenv("BE_HOST")
	BE_PORT  = os.Getenv("BE_PORT")
)
