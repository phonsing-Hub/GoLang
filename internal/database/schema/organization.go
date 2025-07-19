package schema

import "time"

// CreateOrganization represents the schema for creating a new organization
type CreateOrganization struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Slug        string `json:"slug" validate:"required,min=3,max=100,alphanum"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	PlanType    string `json:"plan_type" validate:"omitempty,oneof=free pro enterprise"`
}

// UpdateOrganization represents the schema for updating organization data
type UpdateOrganization struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	PlanType    string `json:"plan_type" validate:"omitempty,oneof=free pro enterprise"`
}

// OrganizationResponse represents organization data for responses
type OrganizationResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	PlanType    string    `json:"plan_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InviteMember represents the schema for inviting a member to organization
type InviteMember struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=owner admin member guest"`
}

// UpdateMemberRole represents the schema for updating member role
type UpdateMemberRole struct {
	Role string `json:"role" validate:"required,oneof=owner admin member guest"`
}

// OrganizationMemberResponse represents organization member data for responses
type OrganizationMemberResponse struct {
	ID           uint                 `json:"id"`
	Role         string               `json:"role"`
	InvitedAt    *time.Time           `json:"invited_at"`
	JoinedAt     *time.Time           `json:"joined_at"`
	CreatedAt    time.Time            `json:"created_at"`
	User         UserInfo             `json:"user"`
	Organization OrganizationResponse `json:"organization"`
}
