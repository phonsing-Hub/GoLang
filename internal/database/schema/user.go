package schema

import "time"

// UserLocation represents location data for user creation
type UserLocation struct {
	LocationType  string   `json:"location_type" validate:"required,oneof=primary shipping billing work"`
	AddressLine1  string   `json:"address_line1" validate:"required"`
	AddressLine2  string   `json:"address_line2"`
	City          string   `json:"city" validate:"required"`
	StateProvince string   `json:"state_province"`
	PostalCode    string   `json:"postal_code" validate:"required"`
	Country       string   `json:"country" validate:"required"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	IsDefault     bool     `json:"is_default"`
}

// UserPreference represents user preference data
type UserPreference struct {
	Key     string `json:"key" validate:"required"`
	Value   string `json:"value" validate:"required"`
	Context string `json:"context"` // global, org_123, project_456
}

// CreateUser represents the schema for creating a new user
type CreateUser struct {
	Email              string     `json:"email" validate:"required,email"`
	Username           string     `json:"username" validate:"omitempty,min=3,max=50"`
	Password           string     `json:"password" validate:"required,min=8"`
	FirstName          string     `json:"first_name" validate:"required,min=1,max=100"`
	LastName           string     `json:"last_name" validate:"omitempty,max=100"`
	DisplayName        string     `json:"display_name" validate:"omitempty,max=100"`
	Bio                string     `json:"bio" validate:"omitempty,max=500"`
	PhoneNumber        string     `json:"phone_number" validate:"omitempty,min=10,max=20"`
	Gender             string     `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth        *time.Time `json:"date_of_birth"`
	LanguagePreference string     `json:"language_preference" validate:"omitempty,min=2,max=5"`
	TimeZone           string     `json:"time_zone" validate:"omitempty"`

	// Optional data
	Locations   []UserLocation   `json:"locations" validate:"omitempty,dive"`
	Preferences []UserPreference `json:"preferences" validate:"omitempty,dive"`
}

// UpdateUser represents the schema for updating user data
type UpdateUser struct {
	Email              string     `json:"email" validate:"omitempty,email"`
	FirstName          string     `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName           string     `json:"last_name" validate:"omitempty,max=100"`
	DisplayName        string     `json:"display_name" validate:"omitempty,max=100"`
	Bio                string     `json:"bio" validate:"omitempty,max=500"`
	PhoneNumber        string     `json:"phone_number" validate:"omitempty,min=10,max=20"`
	Gender             string     `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth        *time.Time `json:"date_of_birth"`
	LanguagePreference string     `json:"language_preference" validate:"omitempty,min=2,max=5"`
	TimeZone           string     `json:"time_zone" validate:"omitempty"`
}

// CreateUserResponse represents the response after creating a user
type CreateUserResponse struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DisplayName     string    `json:"display_name"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
}

// UserProfile represents user profile data for responses
type UserProfileResponse struct {
	ID                 uint       `json:"id"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	DisplayName        string     `json:"display_name"`
	Bio                string     `json:"bio"`
	Avatar             string     `json:"avatar"`
	DateOfBirth        *time.Time `json:"date_of_birth"`
	Gender             string     `json:"gender"`
	PhoneNumber        string     `json:"phone_number"`
	IsPhoneVerified    bool       `json:"is_phone_verified"`
	LanguagePreference string     `json:"language_preference"`
	TimeZone           string     `json:"time_zone"`
	LastActivityAt     *time.Time `json:"last_activity_at"`
}

// UserResponse represents complete user data for responses
type UserResponse struct {
	ID              uint                 `json:"id"`
	Email           string               `json:"email"`
	Username        *string              `json:"username"`
	IsEmailVerified bool                 `json:"is_email_verified"`
	LastLoginAt     *time.Time           `json:"last_login_at"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Profile         *UserProfileResponse `json:"profile,omitempty"`
	Locations       []UserLocation       `json:"locations,omitempty"`
	Preferences     []UserPreference     `json:"preferences,omitempty"`
}
