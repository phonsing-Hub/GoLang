package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	Development        bool
	LogLevel           string
	DBUrl              string
	JWTSecret          string
	CORSAllowOrigins   string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	GoogleTokenURL     string
	GoogleGrantType    string
}

var Env *Config

func LoadEnv() {

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env: %v", err)
	}

	Env = &Config{
		AppPort:            getEnv("APP_PORT", "3000"),
		Development:        getEnv("DEV_MODE", "false") == "true",
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DBUrl:              getEnv("DATABASE_URL", "postgres://user:pass@localhost/db"),
		JWTSecret:          getEnv("JWT_SECRET", "your_jwt_secret"),
		CORSAllowOrigins:   getEnv("CORS_ALLOW_ORIGINS", "*"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", "your_google_client_id"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "your_google_client_secret"),
		GoogleRedirectURI:  getEnv("GOOGLE_REDIRECT_URI", "http://localhost:3000/api/v1/auth/google/callback"),
		GoogleTokenURL:     getEnv("GOOGLE_TOKEN_URL", "https://oauth2.googleapis.com/token"),
		GoogleGrantType:    getEnv("GOOGLE_GRANT_TYPE", "authorization_code"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
