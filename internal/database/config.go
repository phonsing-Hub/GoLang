package database

import (
	"fmt"
	"time"
)

// DatabaseConfig represents configuration for a single database
type DatabaseConfig struct {
	Host         string        `env:"DB_HOST" default:"localhost"`
	Port         int           `env:"DB_PORT" default:"5432"`
	User         string        `env:"DB_USER" default:"postgres"`
	Password     string        `env:"DB_PASSWORD" default:"password"`
	Database     string        `env:"DB_NAME" default:"ticket_system"`
	SSLMode      string        `env:"DB_SSLMODE" default:"disable"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" default:"10"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" default:"100"`
	MaxLifetime  time.Duration `env:"DB_MAX_LIFETIME" default:"1h"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string `env:"REDIS_HOST" default:"localhost"`
	Port     int    `env:"REDIS_PORT" default:"6379"`
	Password string `env:"REDIS_PASSWORD" default:""`
	DB       int    `env:"REDIS_DB" default:"0"`
}

// ServiceConfig for future microservices integration
type ServiceConfig struct {
	UseUserService     bool   `env:"USE_USER_SERVICE" default:"false"`
	UseProjectService  bool   `env:"USE_PROJECT_SERVICE" default:"false"`
	UseActivityService bool   `env:"USE_ACTIVITY_SERVICE" default:"false"`
	UserServiceURL     string `env:"USER_SERVICE_URL" default:""`
	ProjectServiceURL  string `env:"PROJECT_SERVICE_URL" default:""`
	ActivityServiceURL string `env:"ACTIVITY_SERVICE_URL" default:""`
}

// Config represents complete database configuration
type Config struct {
	// Primary database (used initially for everything)
	Primary DatabaseConfig `json:"primary"`

	// Future dedicated databases (optional)
	User     *DatabaseConfig `json:"user,omitempty"`
	Project  *DatabaseConfig `json:"project,omitempty"`
	Activity *DatabaseConfig `json:"activity,omitempty"`

	// Cache and services
	Redis    RedisConfig   `json:"redis"`
	Services ServiceConfig `json:"services"`

	// Migration settings
	AutoMigrate bool `env:"DB_AUTO_MIGRATE" default:"true"`
	LogLevel    int  `env:"DB_LOG_LEVEL" default:"1"` // 1=Info, 2=Warn, 3=Error
}

// DSN generates database connection string
func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dc.Host, dc.User, dc.Password, dc.Database, dc.Port, dc.SSLMode)
}

// RedisAddr generates Redis connection address
func (rc *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", rc.Host, rc.Port)
}
