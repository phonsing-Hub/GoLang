package schema

import "time"

// CreateTicket represents the schema for creating a new ticket
type CreateTicket struct {
	ProjectID      uint       `json:"project_id" validate:"required"`
	EpicID         *uint      `json:"epic_id"`
	Title          string     `json:"title" validate:"required,min=1,max=500"`
	Description    string     `json:"description"`
	TypeID         uint       `json:"type_id" validate:"omitempty"`
	StatusID       uint       `json:"status_id" validate:"required"`
	PriorityID     uint       `json:"priority_id" validate:"omitempty"`
	AssigneeID     *uint      `json:"assignee_id"`
	ParentID       *uint      `json:"parent_id"`
	EstimatedHours *float64   `json:"estimated_hours" validate:"omitempty,min=0"`
	DueDate        *time.Time `json:"due_date"`
	Labels         []uint     `json:"labels"` // Label IDs
}

// UpdateTicket represents the schema for updating ticket data
type UpdateTicket struct {
	Title          string     `json:"title" validate:"omitempty,min=1,max=500"`
	Description    string     `json:"description"`
	TypeID         uint       `json:"type_id" validate:"omitempty"`
	StatusID       uint       `json:"status_id" validate:"omitempty"`
	PriorityID     uint       `json:"priority_id" validate:"omitempty"`
	AssigneeID     *uint      `json:"assignee_id"`
	EstimatedHours *float64   `json:"estimated_hours" validate:"omitempty,min=0"`
	ActualHours    *float64   `json:"actual_hours" validate:"omitempty,min=0"`
	DueDate        *time.Time `json:"due_date"`
}

// TicketResponse represents ticket data for responses
type TicketResponse struct {
	ID             uint                 `json:"id"`
	ProjectID      uint                 `json:"project_id"`
	EpicID         *uint                `json:"epic_id"`
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	TicketKey      string               `json:"ticket_key"`
	EstimatedHours *float64             `json:"estimated_hours"`
	ActualHours    *float64             `json:"actual_hours"`
	DueDate        *time.Time           `json:"due_date"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Type           TicketTypeResponse   `json:"type"`
	Status         TicketStatusResponse `json:"status"`
	Priority       PriorityResponse     `json:"priority"`
	Assignee       *UserInfo            `json:"assignee"`
	Reporter       UserInfo             `json:"reporter"`
	Project        ProjectResponse      `json:"project"`
}

// CreateTicketComment represents the schema for creating a ticket comment
type CreateTicketComment struct {
	TicketID uint   `json:"ticket_id" validate:"required"`
	Content  string `json:"content" validate:"required,min=1"`
}

// UpdateTicketComment represents the schema for updating a ticket comment
type UpdateTicketComment struct {
	Content string `json:"content" validate:"required,min=1"`
}

// TicketCommentResponse represents ticket comment data for responses
type TicketCommentResponse struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserInfo  `json:"user"`
}

// CreateTimeLog represents the schema for creating a time log
type CreateTimeLog struct {
	TicketID    uint       `json:"ticket_id" validate:"required"`
	Description string     `json:"description"`
	Hours       float64    `json:"hours" validate:"required,min=0.25,max=24"`
	LogDate     *time.Time `json:"log_date"`
}

// UpdateTimeLog represents the schema for updating a time log
type UpdateTimeLog struct {
	Description string     `json:"description"`
	Hours       float64    `json:"hours" validate:"omitempty,min=0.25,max=24"`
	LogDate     *time.Time `json:"log_date"`
}

// TimeLogResponse represents time log data for responses
type TimeLogResponse struct {
	ID          uint      `json:"id"`
	TicketID    uint      `json:"ticket_id"`
	Description string    `json:"description"`
	Hours       float64   `json:"hours"`
	LogDate     time.Time `json:"log_date"`
	CreatedAt   time.Time `json:"created_at"`
	User        UserInfo  `json:"user"`
}

// TicketTypeResponse represents ticket type data for responses
type TicketTypeResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// PriorityResponse represents priority data for responses
type PriorityResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// CreateLabel represents the schema for creating a label
type CreateLabel struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"omitempty,max=200"`
	Color       string `json:"color" validate:"required,hexcolor"`
}

// UpdateLabel represents the schema for updating a label
type UpdateLabel struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=50"`
	Description string `json:"description" validate:"omitempty,max=200"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
}

// LabelResponse represents label data for responses
type LabelResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
