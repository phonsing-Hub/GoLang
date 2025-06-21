package models
import (
  "gorm.io/gorm"
  
)
type Users struct {
	gorm.Model
	
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Role     string `gorm:"not null" json:"role"`
	// Add other fields as necessary
}