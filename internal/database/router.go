package database

import (
	"gorm.io/gorm"
)

// Router provides intelligent database routing
type Router struct {
	manager *Manager
}

// NewRouter creates a new database router
func NewRouter(manager *Manager) *Router {
	return &Router{
		manager: manager,
	}
}

// DatabaseType represents different database contexts
type DatabaseType string

const (
	// User-related data
	UserDatabase DatabaseType = "user"
	AuthDatabase DatabaseType = "auth"

	// Project-related data
	ProjectDatabase DatabaseType = "project"
	TicketDatabase  DatabaseType = "ticket"

	// Activity-related data
	ActivityDatabase DatabaseType = "activity"
	LogDatabase      DatabaseType = "log"

	// Default database
	DefaultDatabase DatabaseType = "default"
)

// GetDB returns the appropriate database for the given context
func (r *Router) GetDB(dbType DatabaseType) *gorm.DB {
	switch dbType {
	case UserDatabase, AuthDatabase:
		return r.getUserDB()
	case ProjectDatabase, TicketDatabase:
		return r.getProjectDB()
	case ActivityDatabase, LogDatabase:
		return r.getActivityDB()
	default:
		return r.manager.Primary
	}
}

// getUserDB returns user database or fallback to primary
func (r *Router) getUserDB() *gorm.DB {
	if r.manager.UserDB != nil {
		return r.manager.UserDB
	}
	return r.manager.Primary
}

// getProjectDB returns project database or fallback to primary
func (r *Router) getProjectDB() *gorm.DB {
	if r.manager.ProjectDB != nil {
		return r.manager.ProjectDB
	}
	return r.manager.Primary
}

// getActivityDB returns activity database or fallback to primary
func (r *Router) getActivityDB() *gorm.DB {
	if r.manager.ActivityDB != nil {
		return r.manager.ActivityDB
	}
	return r.manager.Primary
}

// GetCache returns Redis client
func (r *Router) GetCache() interface{} { // Using interface{} to avoid Redis import error
	return r.manager.Cache
}

// GetPrimary returns primary database (always available)
func (r *Router) GetPrimary() *gorm.DB {
	return r.manager.Primary
}

// IsUsingDedicatedDB checks if a dedicated database is configured
func (r *Router) IsUsingDedicatedDB(dbType DatabaseType) bool {
	switch dbType {
	case UserDatabase, AuthDatabase:
		return r.manager.UserDB != nil
	case ProjectDatabase, TicketDatabase:
		return r.manager.ProjectDB != nil
	case ActivityDatabase, LogDatabase:
		return r.manager.ActivityDB != nil
	default:
		return false
	}
}

// GetDatabaseInfo returns information about database usage
func (r *Router) GetDatabaseInfo() map[string]interface{} {
	info := map[string]interface{}{
		"primary_db": "active",
		"cache":      r.manager.Cache != nil,
	}

	if r.manager.UserDB != nil {
		info["user_db"] = "active"
	} else {
		info["user_db"] = "using_primary"
	}

	if r.manager.ProjectDB != nil {
		info["project_db"] = "active"
	} else {
		info["project_db"] = "using_primary"
	}

	if r.manager.ActivityDB != nil {
		info["activity_db"] = "active"
	} else {
		info["activity_db"] = "using_primary"
	}

	return info
}

// Helper methods for common patterns

// WithUserContext returns database for user-related operations
func (r *Router) WithUserContext() *gorm.DB {
	return r.GetDB(UserDatabase)
}

// WithProjectContext returns database for project-related operations
func (r *Router) WithProjectContext() *gorm.DB {
	return r.GetDB(ProjectDatabase)
}

// WithActivityContext returns database for activity-related operations
func (r *Router) WithActivityContext() *gorm.DB {
	return r.GetDB(ActivityDatabase)
}

// Transaction helper for cross-database operations
type TransactionFunc func(dbs map[DatabaseType]*gorm.DB) error

// RunTransaction executes a function within database transactions
// Note: This is for single database transactions initially
// In the future, we'll implement distributed transactions
func (r *Router) RunTransaction(dbType DatabaseType, fn func(*gorm.DB) error) error {
	db := r.GetDB(dbType)
	return db.Transaction(fn)
}

// RunMultiTransaction executes a function across multiple databases
// Note: This is a placeholder for future distributed transaction support
func (r *Router) RunMultiTransaction(dbTypes []DatabaseType, fn TransactionFunc) error {
	// For now, we'll use the primary database for multi-DB operations
	// In the future, implement Saga pattern or 2PC
	dbs := make(map[DatabaseType]*gorm.DB)
	for _, dbType := range dbTypes {
		dbs[dbType] = r.GetDB(dbType)
	}

	// Execute within primary database transaction for consistency
	return r.manager.Primary.Transaction(func(tx *gorm.DB) error {
		// Override with transaction for primary DB operations
		for dbType := range dbs {
			if r.GetDB(dbType) == r.manager.Primary {
				dbs[dbType] = tx
			}
		}
		return fn(dbs)
	})
}
