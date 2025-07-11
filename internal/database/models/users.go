package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model

	Email           string `gorm:"uniqueIndex;not null" json:"email"`
	IsEmailVerified bool   `gorm:"default:false" json:"is_email_verified"`
	PhoneNumber     string `gorm:"uniqueIndex;default:null" json:"phone_number"`
	IsPhoneVerified bool   `gorm:"default:false" json:"is_phone_verified"`

	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	DisplayName        string     `json:"display_name"`
	ProfilePictureURL  string     `gorm:"default:null" json:"profile_picture_url"`
	Bio                string     `gorm:"type:text;default:null" json:"bio"`
	DateOfBirth        *time.Time `gorm:"default:null" json:"date_of_birth"`
	Gender             string     `gorm:"default:null" json:"gender"`
	Avatar             string     `gorm:"default:null" json:"avatar"`
	LanguagePreference string     `gorm:"default:'en'" json:"language_preference"`
	TimeZone           string     `gorm:"default:'UTC'" json:"time_zone"`
	LastActivityAt     *time.Time `gorm:"default:null" json:"last_activity_at"`

	CurrentStatusID *uint       `gorm:"default:1" json:"current_status_id"`                     // Foreign Key to UserStatuses, nullable if no status set
	CurrentStatus   *UserStatus `gorm:"foreignKey:CurrentStatusID" json:"UserStatus,omitempty"` // Belongs To UserStatus

	Credentials []UserCredentials `gorm:"foreignKey:UserID"` // A user can have many credential types
	Locations   []UserLocation    `gorm:"foreignKey:UserID"` // A user can have many locations

}

type UserCredentials struct {
	gorm.Model

	UserID uint `gorm:"not null;uniqueIndex:idx_user_credential_type_unique" json:"user_id"` // Foreign Key to Users table

	CredentialType string `gorm:"not null;uniqueIndex:idx_user_credential_type_unique" json:"credential_type"`

	Username     string `gorm:"uniqueIndex;default:null" json:"username"` // Optional: If separate username is used for login
	PasswordHash string `gorm:"default:null" json:"password_hash"`        // Hashed password
	PasswordSalt string `gorm:"default:null" json:"password_salt"`        // Salt used for hashing the password

	OAuthProvider   string `gorm:"default:null" json:"oauth_provider"`    // e.g., "google", "facebook"
	OAuthProviderID string `gorm:"default:null" json:"oauth_provider_id"` // ID from the OAuth provider

	User *Users `gorm:"foreignKey:UserID" json:"User,omitempty"`
}

type UserLocation struct {
	gorm.Model

	UserID       uint   `gorm:"not null" json:"user_id"`       // Foreign Key to Users table
	LocationType string `gorm:"not null" json:"location_type"` // e.g., "primary", "shipping", "billing", "work"

	AddressLine1  string   `gorm:"not null" json:"address_line1"`
	AddressLine2  string   `gorm:"default:null" json:"address_line2"`
	City          string   `gorm:"not null" json:"city"`
	StateProvince string   `gorm:"default:null" json:"state_province"`
	PostalCode    string   `gorm:"not null" json:"postal_code"`
	Country       string   `gorm:"not null" json:"country"`
	Latitude      *float64 `gorm:"default:null" json:"latitude"`    // Optional: For geographic coordinates
	Longitude     *float64 `gorm:"default:null" json:"longitude"`   // Optional: For geographic coordinates
	IsDefault     bool     `gorm:"default:false" json:"is_default"` // Flag to mark a default location

	User *Users `gorm:"foreignKey:UserID" json:"User,omitempty"`
}

type UserStatus struct {
	ID          uint   `gorm:"primaryKey" json:"id"`                      // Custom ID for specific status values (e.g., 1 for "Active", 2 for "Inactive")
	StatusName  string `gorm:"unique;not null" json:"status_name"`        // e.g., "Active", "Inactive", "Suspended", "Pending Verification"
	Description string `gorm:"type:text;default:null" json:"description"` // Optional description for the status

}
