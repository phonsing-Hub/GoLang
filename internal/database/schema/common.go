package schema

import "time"

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents error details in API responses
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page    int `form:"page" validate:"omitempty,min=1"`
	PerPage int `form:"per_page" validate:"omitempty,min=1,max=100"`
}

// SortRequest represents sorting parameters
type SortRequest struct {
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// FilterRequest represents common filtering parameters
type FilterRequest struct {
	Search    string     `form:"search"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	Status    string     `form:"status"`
}

// IDResponse represents a simple ID response
type IDResponse struct {
	ID uint `json:"id"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// BoolResponse represents a simple boolean response
type BoolResponse struct {
	Success bool `json:"success"`
}

// CountResponse represents a count response
type CountResponse struct {
	Count int `json:"count"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}
