package models

import (
	"gorm.io/gorm"
)

type UserCredentials struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt fields

	UserID uint `gorm:"not null;uniqueIndex:idx_user_credential_type_unique" json:"user_id"` // Foreign Key to Users table

	// CredentialType defines the type of credential (ee.g., "password", "oauth_google", "api_key")
	// Combined with UserID, this can ensure unique credentials of a certain type per user.
	CredentialType string `gorm:"not null;uniqueIndex:idx_user_credential_type_unique" json:"credential_type"`

	// For password-based authentication
	Username     string `gorm:"uniqueIndex;default:null" json:"username"` // Optional: If separate username is used for login
	PasswordHash string `gorm:"default:null" json:"password_hash"`        // Hashed password
	PasswordSalt string `gorm:"default:null" json:"password_salt"`        // Salt used for hashing the password

	// For OAuth-based authentication (if applicable)
	OAuthProvider   string `gorm:"default:null" json:"oauth_provider"`    // e.g., "google", "facebook"
	OAuthProviderID string `gorm:"default:null" json:"oauth_provider_id"` // ID from the OAuth provider

	// Add other credential-specific fields here, e.g.:
	// LastLoginAt time.Time
	// FailedLoginAttempts int
	// PasswordResetToken string

	// Define the Belongs To relationship with Users
	// This tells GORM that a UserCredential belongs to a User.
	User *Users `gorm:"foreignKey:UserID" json:"User,omitempty"`
}
