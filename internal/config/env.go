package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	Development      bool
	LogLevel         string
	DBUrl            string
	JWTSecret        string
	CORSAllowOrigins string
}

var Env *Config

func LoadEnv() {

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env: %v", err)
	}

	Env = &Config{
		AppPort:          getEnv("APP_PORT", "3000"),
		Development:      getEnv("DEV_MODE", "false") == "true",
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		DBUrl:            getEnv("DATABASE_URL", "postgres://user:pass@localhost/db"),
		JWTSecret:        getEnv("JWT_SECRET", "your_jwt_secret"),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
