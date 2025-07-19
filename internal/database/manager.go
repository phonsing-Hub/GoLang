package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager handles all database connections and routing
type Manager struct {
	// Current databases
	Primary *gorm.DB
	Cache   *redis.Client

	// Future databases (nil initially, activated when needed)
	UserDB     *gorm.DB
	ProjectDB  *gorm.DB
	ActivityDB *gorm.DB

	// Configuration
	config *Config
	ctx    context.Context
}

// NewManager creates a new database manager
func NewManager(ctx context.Context, config *Config) (*Manager, error) {
	manager := &Manager{
		config: config,
		ctx:    ctx,
	}

	// Initialize primary database (always required)
	if err := manager.initPrimaryDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize primary database: %w", err)
	}

	// Initialize Redis cache
	if err := manager.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize additional databases if configured
	if err := manager.initAdditionalDBs(); err != nil {
		return nil, fmt.Errorf("failed to initialize additional databases: %w", err)
	}

	log.Println("Database Manager initialized successfully")
	return manager, nil
}

// initPrimaryDB initializes the primary PostgreSQL database
func (m *Manager) initPrimaryDB() error {
	db, err := m.openPostgreSQL(&m.config.Primary)
	if err != nil {
		return fmt.Errorf("failed to open primary database: %w", err)
	}

	m.Primary = db
	log.Println("Primary database connected")
	return nil
}

// initRedis initializes Redis connection
func (m *Manager) initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     m.config.Redis.Addr(),
		Password: m.config.Redis.Password,
		DB:       m.config.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	m.Cache = rdb
	log.Println("Redis connected")
	return nil
}

// initAdditionalDBs initializes dedicated databases if configured
func (m *Manager) initAdditionalDBs() error {
	// Initialize User DB if configured
	if m.config.User != nil {
		db, err := m.openPostgreSQL(m.config.User)
		if err != nil {
			return fmt.Errorf("failed to open user database: %w", err)
		}
		m.UserDB = db
		log.Println("User database connected")
	}

	// Initialize Project DB if configured
	if m.config.Project != nil {
		db, err := m.openPostgreSQL(m.config.Project)
		if err != nil {
			return fmt.Errorf("failed to open project database: %w", err)
		}
		m.ProjectDB = db
		log.Println("Project database connected")
	}

	// Initialize Activity DB if configured
	if m.config.Activity != nil {
		db, err := m.openPostgreSQL(m.config.Activity)
		if err != nil {
			return fmt.Errorf("failed to open activity database: %w", err)
		}
		m.ActivityDB = db
		log.Println("Activity database connected")
	}

	return nil
}

// openPostgreSQL opens a PostgreSQL connection with proper configuration
func (m *Manager) openPostgreSQL(config *DatabaseConfig) (*gorm.DB, error) {
	// Configure GORM logger
	logLevel := logger.LogLevel(m.config.LogLevel)
	gormLogger := logger.Default.LogMode(logLevel)

	// Open database connection
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close closes all database connections
func (m *Manager) Close() error {
	var errors []error

	// Close primary database
	if m.Primary != nil {
		if sqlDB, err := m.Primary.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close primary DB: %w", err))
			}
		}
	}

	// Close additional databases
	for name, db := range map[string]*gorm.DB{
		"UserDB":     m.UserDB,
		"ProjectDB":  m.ProjectDB,
		"ActivityDB": m.ActivityDB,
	} {
		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				if err := sqlDB.Close(); err != nil {
					errors = append(errors, fmt.Errorf("failed to close %s: %w", name, err))
				}
			}
		}
	}

	// Close Redis
	if m.Cache != nil {
		if err := m.Cache.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Redis: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing databases: %v", errors)
	}

	log.Println("All database connections closed")
	return nil
}

// Health checks the health of all database connections
func (m *Manager) Health(ctx context.Context) map[string]error {
	health := make(map[string]error)

	// Check primary database
	if m.Primary != nil {
		if sqlDB, err := m.Primary.DB(); err == nil {
			health["primary"] = sqlDB.PingContext(ctx)
		} else {
			health["primary"] = err
		}
	}

	// Check additional databases
	for name, db := range map[string]*gorm.DB{
		"user":     m.UserDB,
		"project":  m.ProjectDB,
		"activity": m.ActivityDB,
	} {
		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				health[name] = sqlDB.PingContext(ctx)
			} else {
				health[name] = err
			}
		}
	}

	// Check Redis
	if m.Cache != nil {
		health["redis"] = m.Cache.Ping(ctx).Err()
	}

	return health
}
