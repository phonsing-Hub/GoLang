package models


type UserStatus struct {
	ID          uint   `gorm:"primaryKey" json:"id"` // Custom ID for specific status values (e.g., 1 for "Active", 2 for "Inactive")
	StatusName  string `gorm:"unique;not null" json:"status_name"` // e.g., "Active", "Inactive", "Suspended", "Pending Verification"
	Description string `gorm:"type:text;default:null" json:"description"` // Optional description for the status
	// You might add an IsDefault bool `gorm:"default:false"` to indicate default status
}