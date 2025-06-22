package models

import (
	"gorm.io/gorm"
)


type UserLocation struct {
	gorm.Model // Provides ID, CreatedAt, UpdatedAt, DeletedAt

	UserID        uint   `gorm:"not null" json:"user_id"` // Foreign Key to Users table
	LocationType  string `gorm:"not null" json:"location_type"` // e.g., "primary", "shipping", "billing", "work"

	AddressLine1  string `gorm:"not null" json:"address_line1"`
	AddressLine2  string `gorm:"default:null" json:"address_line2"`
	City          string `gorm:"not null" json:"city"`
	StateProvince string `gorm:"default:null" json:"state_province"`
	PostalCode    string `gorm:"not null" json:"postal_code"`
	Country       string `gorm:"not null" json:"country"`
	Latitude      *float64 `gorm:"default:null" json:"latitude"`  // Optional: For geographic coordinates
	Longitude     *float64 `gorm:"default:null" json:"longitude"` // Optional: For geographic coordinates
	IsDefault     bool   `gorm:"default:false" json:"is_default"` // Flag to mark a default location

	 User *Users `gorm:"foreignKey:UserID" json:"User,omitempty"`
}