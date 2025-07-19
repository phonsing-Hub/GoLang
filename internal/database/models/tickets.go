package models

import (
	"time"

	"gorm.io/gorm"
)

// Epic represents large features or user stories
type Epic struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ProjectID   uint           `json:"project_id" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	EpicKey     string         `json:"epic_key" gorm:"not null;unique;size:20"`     // เช่น "PROJ-E1"
	StatusID    uint           `json:"status_id" gorm:"not null;index;default:1"`   // FK to epic_statuses (1=planning)
	PriorityID  uint           `json:"priority_id" gorm:"not null;index;default:2"` // FK to priorities (2=medium)
	OwnerID     uint           `json:"owner_id" gorm:"not null;index"`
	StartDate   *time.Time     `json:"start_date"`
	TargetDate  *time.Time     `json:"target_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Status   EpicStatus `json:"status" gorm:"foreignKey:StatusID"`
	Priority Priority   `json:"priority" gorm:"foreignKey:PriorityID"`
	Project  Project    `json:"project" gorm:"foreignKey:ProjectID"`
	Owner    User       `json:"owner" gorm:"foreignKey:OwnerID"`
	Tickets  []Ticket   `json:"tickets,omitempty" gorm:"foreignKey:EpicID"`
	Labels   []Label    `json:"labels,omitempty" gorm:"many2many:epic_labels"`
}

// Ticket represents a ticket/issue in the system
type Ticket struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ProjectID      uint           `json:"project_id" gorm:"not null;index"`
	EpicID         *uint          `json:"epic_id" gorm:"index"`
	Title          string         `json:"title" gorm:"not null;size:500"`
	Description    string         `json:"description" gorm:"type:text"`
	TicketKey      string         `json:"ticket_key" gorm:"not null;unique;size:20"`
	TypeID         uint           `json:"type_id" gorm:"not null;index;default:1"` // FK to ticket_types (1=task)
	StatusID       uint           `json:"status_id" gorm:"not null;index"`
	PriorityID     uint           `json:"priority_id" gorm:"not null;index;default:2"` // FK to priorities (2=medium)
	AssigneeID     *uint          `json:"assignee_id" gorm:"index"`
	ReporterID     uint           `json:"reporter_id" gorm:"not null;index"`
	ParentID       *uint          `json:"parent_id" gorm:"index"`
	EstimatedHours *float64       `json:"estimated_hours"`
	ActualHours    *float64       `json:"actual_hours"`
	DueDate        *time.Time     `json:"due_date"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Project     Project            `json:"project" gorm:"foreignKey:ProjectID"`
	Epic        *Epic              `json:"epic,omitempty" gorm:"foreignKey:EpicID"`
	Type        TicketType         `json:"type" gorm:"foreignKey:TypeID"`
	Status      TicketStatus       `json:"status" gorm:"foreignKey:StatusID"`
	Priority    Priority           `json:"priority" gorm:"foreignKey:PriorityID"`
	Assignee    *User              `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	Reporter    User               `json:"reporter" gorm:"foreignKey:ReporterID"`
	Parent      *Ticket            `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []Ticket           `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Comments    []TicketComment    `json:"comments,omitempty"`
	Attachments []TicketAttachment `json:"attachments,omitempty"`
	Labels      []Label            `json:"labels,omitempty" gorm:"many2many:ticket_labels"`
	Watchers    []User             `json:"watchers,omitempty" gorm:"many2many:ticket_watchers"`
	TimeLogs    []TimeLog          `json:"time_logs,omitempty"`
}

// TicketStatus represents status of tickets
type TicketStatus struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Position  int       `json:"position" gorm:"not null;default:0;index"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
	Color     string    `json:"color" gorm:"size:7"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:StatusID"`
}

// TicketComment represents comments on tickets
type TicketComment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TicketID  uint           `json:"ticket_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Ticket Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

// TicketAttachment represents file attachments on tickets
type TicketAttachment struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TicketID   uint      `json:"ticket_id" gorm:"not null;index"`
	Filename   string    `json:"filename" gorm:"not null;size:255"`
	FilePath   string    `json:"file_path" gorm:"not null;size:500"`
	FileSize   int64     `json:"file_size" gorm:"not null"`
	MimeType   string    `json:"mime_type" gorm:"size:100"`
	UploadedBy uint      `json:"uploaded_by" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Ticket   Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	Uploader User   `json:"uploader" gorm:"foreignKey:UploadedBy"`
}

// Label represents labels/tags for tickets
type Label struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ProjectID   uint      `json:"project_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null;size:100;index"`
	Color       string    `json:"color" gorm:"size:7"` // hex color
	Description string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Project Project  `json:"project" gorm:"foreignKey:ProjectID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"many2many:ticket_labels"`
}

// TimeLog represents time tracking for tickets
type TimeLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TicketID    uint      `json:"ticket_id" gorm:"not null;index"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Hours       float64   `json:"hours" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	LoggedDate  time.Time `json:"logged_date" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Ticket Ticket `json:"ticket" gorm:"foreignKey:TicketID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

// Sprint represents agile sprints
type Sprint struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	ProjectID uint       `json:"project_id" gorm:"not null;index"`
	Name      string     `json:"name" gorm:"not null;size:255"`
	StartDate *time.Time `json:"start_date" gorm:"index"`
	EndDate   *time.Time `json:"end_date" gorm:"index"`
	StatusID  uint       `json:"status_id" gorm:"not null;index;default:1"` // FK to sprint_statuses (1=planning)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	Status  SprintStatus `json:"status" gorm:"foreignKey:StatusID"`
	Project Project      `json:"project" gorm:"foreignKey:ProjectID"`
	Tickets []Ticket     `json:"tickets,omitempty" gorm:"many2many:sprint_tickets"`
}

// Table names

func (Epic) TableName() string {
	return "epics"
}

func (TicketStatus) TableName() string {
	return "ticket_statuses"
}

func (TicketComment) TableName() string {
	return "ticket_comments"
}

func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}

func (TimeLog) TableName() string {
	return "time_logs"
}
