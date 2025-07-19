package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the system
type Project struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrganizationID *uint          `json:"organization_id" gorm:"index"` // nullable สำหรับ personal projects
	Name           string         `json:"name" gorm:"not null;size:255"`
	Description    string         `json:"description" gorm:"type:text"`
	Key            string         `json:"key" gorm:"not null;unique;size:10"` // เช่น "PROJ"
	OwnerID        uint           `json:"owner_id" gorm:"not null;index"`
	StatusID       uint           `json:"status_id" gorm:"not null;index;default:1"` // FK to project_statuses (1=active)
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Status       ProjectStatus   `json:"status" gorm:"foreignKey:StatusID"`
	Organization *Organization   `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	Owner        User            `json:"owner" gorm:"foreignKey:OwnerID"`
	Members      []ProjectMember `json:"members,omitempty"`
	Tickets      []Ticket        `json:"tickets,omitempty"`
	Statuses     []TicketStatus  `json:"statuses,omitempty"`
	Labels       []Label         `json:"labels,omitempty"`
	Sprints      []Sprint        `json:"sprints,omitempty"`
}

// ProjectMember represents project membership
type ProjectMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Role      string    `json:"role" gorm:"not null;default:'developer';index"` // admin, developer, viewer
	JoinedAt  time.Time `json:"joined_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships
	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}

// Indexes
func (ProjectMember) TableName() string {
	return "project_members"
}

// Add unique constraint for project_id + user_id
type ProjectMemberIndex struct {
	ProjectID uint `gorm:"uniqueIndex:idx_project_user"`
	UserID    uint `gorm:"uniqueIndex:idx_project_user"`
}
