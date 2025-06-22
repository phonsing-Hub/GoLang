package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt

	// Essential Identification
	Email             string    `gorm:"uniqueIndex;not null" json:"email"`
	IsEmailVerified   bool      `gorm:"default:false" json:"is_email_verified"`
	PhoneNumber       string    `gorm:"uniqueIndex;default:null" json:"phone_number"`
	IsPhoneVerified   bool      `gorm:"default:false" json:"is_phone_verified"`

	// Personal Information
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	DisplayName       string    `json:"display_name"`
	ProfilePictureURL string    `gorm:"default:null" json:"profile_picture_url"`
	Bio               string    `gorm:"type:text;default:null" json:"bio"`
	DateOfBirth       *time.Time `gorm:"default:null" json:"date_of_birth"`
	Gender            string    `gorm:"default:null" json:"gender"`

	// Account Activity & Preferences (still in Users for simplicity)
	LanguagePreference string   `gorm:"default:'en'" json:"language_preference"`
	TimeZone          string    `gorm:"default:'UTC'" json:"time_zone"`
	LastActivityAt    *time.Time `gorm:"default:null" json:"last_activity_at"`

	// Relationships
	// A user has one current status (you might also track history with another table)
	CurrentStatusID *uint           `gorm:"default:1" json:"current_status_id"` // Foreign Key to UserStatuses, nullable if no status set
	CurrentStatus   *UserStatus      `gorm:"foreignKey:CurrentStatusID" json:"UserStatus,omitempty"`          // Belongs To UserStatus

	Credentials       []UserCredentials `gorm:"foreignKey:UserID"`       // A user can have many credential types
	Locations         []UserLocation    `gorm:"foreignKey:UserID"`       // A user can have many locations

	// You might also have relationships to other tables like:
	// Orders            []Order           `gorm:"foreignKey:UserID"`
	// Posts             []Post            `gorm:"foreignKey:UserID"`
}
