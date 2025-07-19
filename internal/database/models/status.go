package models

import (
	"time"
)

// UserStatus represents user status lookup table
type UserStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // active, inactive, suspended, pending_verification
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Active, Inactive, Suspended, Pending Verification
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Users []User `json:"users,omitempty" gorm:"foreignKey:StatusID"`
}

// ProjectStatus represents project status lookup table
type ProjectStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // active, archived, on_hold, cancelled
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Active, Archived, On Hold, Cancelled
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Projects []Project `json:"projects,omitempty" gorm:"foreignKey:StatusID"`
}

// EpicStatus represents epic status lookup table
type EpicStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // planning, in_progress, completed, cancelled
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Planning, In Progress, Completed, Cancelled
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Epics []Epic `json:"epics,omitempty" gorm:"foreignKey:StatusID"`
}

// SprintStatus represents sprint status lookup table
type SprintStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // planning, active, completed, cancelled
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Planning, Active, Completed, Cancelled
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Sprints []Sprint `json:"sprints,omitempty" gorm:"foreignKey:StatusID"`
}

// OrganizationStatus represents organization status lookup table
type OrganizationStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // active, suspended, trial, expired
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Active, Suspended, Trial, Expired
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Organizations []Organization `json:"organizations,omitempty" gorm:"foreignKey:StatusID"`
}

// MemberStatus represents membership status lookup table (for organization members)
type MemberStatus struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // active, invited, suspended, inactive
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Active, Invited, Suspended, Inactive
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty" gorm:"foreignKey:StatusID"`
}

// Priority represents priority lookup table (for tickets, epics)
type Priority struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // low, medium, high, critical, blocker
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Low, Medium, High, Critical, Blocker
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"`          // hex color for UI
	Level       int       `json:"level" gorm:"not null;unique"` // 1=Low, 2=Medium, 3=High, 4=Critical, 5=Blocker
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:PriorityID"`
	Epics   []Epic   `json:"epics,omitempty" gorm:"foreignKey:PriorityID"`
}

// TicketType represents ticket type lookup table
type TicketType struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`   // task, bug, feature, improvement, epic
	DisplayName string    `json:"display_name" gorm:"not null;size:100"` // Task, Bug, Feature, Improvement, Epic
	Description string    `json:"description" gorm:"type:text"`
	Icon        string    `json:"icon" gorm:"size:50"` // icon name for UI
	Color       string    `json:"color" gorm:"size:7"` // hex color for UI
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Position    int       `json:"position" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:TypeID"`
}

// Table names
func (UserStatus) TableName() string {
	return "user_statuses"
}

func (ProjectStatus) TableName() string {
	return "project_statuses"
}

func (EpicStatus) TableName() string {
	return "epic_statuses"
}

func (SprintStatus) TableName() string {
	return "sprint_statuses"
}

func (OrganizationStatus) TableName() string {
	return "organization_statuses"
}

func (MemberStatus) TableName() string {
	return "member_statuses"
}

func (Priority) TableName() string {
	return "priorities"
}

func (TicketType) TableName() string {
	return "ticket_types"
}
