package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	Development bool
	LogLevel    string
	DBUrl       string
	JWTSecret   string
}

var Env *Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Env = &Config{
		AppPort:     getEnv("APP_PORT", "3000"),
		Development: getEnv("DEV_MODE", "false") == "true",
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DBUrl:       getEnv("DATABASE_URL", "postgres://user:pass@localhost/db"),
		JWTSecret:   getEnv("JWT_SECRET", "your_jwt_secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
