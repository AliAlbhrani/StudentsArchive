package env

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT string

	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string

	REDIS_PASSWORD string
	REDIS_HOST     string

	DEV bool

	STORAGE_PASSWORD string
	STORAGE_APP_NAME string
	STORAGE_BASE_URL string
	STORAGE_PATH     string

	ISSUER     string
	SECRET_KEY string
)

var _ = func() bool {
	godotenv.Load()

	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")

	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	REDIS_HOST = os.Getenv("REDIS_HOST")

	DEV = os.Getenv("DEV") == "true"

	STORAGE_PASSWORD = os.Getenv("STORAGE_PASSWORD")
	STORAGE_APP_NAME = os.Getenv("STORAGE_APP_NAME")
	STORAGE_BASE_URL = os.Getenv("STORAGE_BASE_URL")
	STORAGE_PATH = os.Getenv("STORAGE_PATH")

	ISSUER = os.Getenv("ISSUER")
	SECRET_KEY = os.Getenv("SECRET_KEY")

	return true
}()
