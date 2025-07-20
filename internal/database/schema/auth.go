package schema

import "time"

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // can be email or username
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"` // Bearer
	ExpiresIn    int      `json:"expires_in"` // seconds
	User         UserInfo `json:"user"`
}

// UserInfo represents basic user information for auth responses
type UserInfo struct {
	ID              uint       `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	DisplayName     string     `json:"display_name"`
	Avatar          string     `json:"avatar"`
	IsEmailVerified bool       `json:"is_email_verified"`
	LastLoginAt     *time.Time `json:"last_login_at"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	FirstName       string `json:"first_name" validate:"required,min=1,max=100"`
	LastName        string `json:"last_name" validate:"omitempty,max=100"`
	DisplayName     string `json:"display_name" validate:"omitempty,max=100"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents password reset request
type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// OAuth login request
type OAuthLoginRequest struct {
	Provider     string `json:"provider" validate:"required,oneof=google github microsoft"`
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleCallbackRequest struct {
    Code  string `json:"code" validate:"required"`
    State string `json:"state,omitempty"`
    Error string `json:"error,omitempty"`
    ErrorDescription string `json:"error_description,omitempty"`
}

type GoogleTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    Scope        string `json:"scope"`
    IDToken      string `json:"id_token,omitempty"`
}

