package response

import (
	"encoding/json"

	"github.com/phonsing-Hub/GoLang/internal/database/models"
)

type SWListDataDetail struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Data  json.RawMessage `json:"data" swaggertype:"array,object"`
}

type SWListSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    *SWListDataDetail `json:"data"`
	Error   json.RawMessage   `json:"error" swaggertype:"object"`
}

type SWSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    json.RawMessage `json:"data" swaggertype:"object"`
	Error   json.RawMessage `json:"error" swaggertype:"object"`
}

type SWErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SWErrorResponse struct {
	Success bool           `json:"success" example:"false"`
	Data    any            `json:"data"`
	Error   *SWErrorDetail `json:"error"`
}

// LoginResponse represents the response structure for login endpoint
type LoginResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		User  models.User `json:"user"`
	} `json:"data"`
	Error interface{} `json:"error" example:"null"`
}

// UserResponse represents the response structure for user endpoints
type UserResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    models.User `json:"data"`
	Error   interface{}  `json:"error" example:"null"`
}

// UsersResponse represents the response structure for users list endpoint
type UsersResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []models.User `json:"data"`
	Error   interface{}    `json:"error" example:"null"`
}

// CreateUserResponse represents the response structure for user creation
type CreateUserResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		UserID uint   `json:"user_id" example:"1"`
		Email  string `json:"email" example:"user@example.com"`
	} `json:"data"`
	Error interface{} `json:"error" example:"null"`
}

// AvatarResponse represents the response structure for avatar operations
type AvatarResponse struct {
	Success bool `json:"success" example:"true"`
	Data    struct {
		Message   string `json:"message" example:"Avatar uploaded successfully"`
		AvatarURL string `json:"avatarUrl,omitempty" example:"/static/avatars/filename.png"`
	} `json:"data"`
	Error interface{} `json:"error" example:"null"`
}
