package schema

import "time"

// CreateProject represents the schema for creating a new project
type CreateProject struct {
	Name           string `json:"name" validate:"required,min=1,max=255"`
	Key            string `json:"key" validate:"required,min=2,max=10,uppercase"`
	Description    string `json:"description" validate:"omitempty,max=1000"`
	ProjectType    string `json:"project_type" validate:"omitempty,oneof=kanban scrum"`
	IsPrivate      bool   `json:"is_private"`
	OrganizationID uint   `json:"organization_id" validate:"required"`
}

// UpdateProject represents the schema for updating project data
type UpdateProject struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	ProjectType string `json:"project_type" validate:"omitempty,oneof=kanban scrum"`
	IsPrivate   bool   `json:"is_private"`
}

// ProjectResponse represents project data for responses
type ProjectResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Key            string    `json:"key"`
	Description    string    `json:"description"`
	ProjectType    string    `json:"project_type"`
	IsPrivate      bool      `json:"is_private"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	OrganizationID uint      `json:"organization_id"`
	OwnerID        uint      `json:"owner_id"`
}

// AddProjectMember represents the schema for adding a member to project
type AddProjectMember struct {
	UserID uint   `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=owner admin member viewer"`
}

// UpdateProjectMemberRole represents the schema for updating project member role
type UpdateProjectMemberRole struct {
	Role string `json:"role" validate:"required,oneof=owner admin member viewer"`
}

// ProjectMemberResponse represents project member data for responses
type ProjectMemberResponse struct {
	ID       uint            `json:"id"`
	Role     string          `json:"role"`
	JoinedAt time.Time       `json:"joined_at"`
	User     UserInfo        `json:"user"`
	Project  ProjectResponse `json:"project"`
}

// CreateTicketStatus represents the schema for creating ticket status
type CreateTicketStatus struct {
	ProjectID uint   `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Position  int    `json:"position" validate:"omitempty,min=0"`
	IsDefault bool   `json:"is_default"`
	Color     string `json:"color" validate:"omitempty,hexcolor"`
}

// UpdateTicketStatus represents the schema for updating ticket status
type UpdateTicketStatus struct {
	Name      string `json:"name" validate:"omitempty,min=1,max=100"`
	Position  int    `json:"position" validate:"omitempty,min=0"`
	IsDefault bool   `json:"is_default"`
	Color     string `json:"color" validate:"omitempty,hexcolor"`
}

// TicketStatusResponse represents ticket status data for responses
type TicketStatusResponse struct {
	ID        uint      `json:"id"`
	ProjectID uint      `json:"project_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	IsDefault bool      `json:"is_default"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
