package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents core user entity
type User struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	Email              string         `json:"email" gorm:"uniqueIndex;not null"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
	DisplayName        string         `json:"display_name"`
	Bio                string         `json:"bio" gorm:"type:text;default:null"`
	Avatar             string         `json:"avatar" gorm:"default:null"`
	DateOfBirth        *time.Time     `json:"date_of_birth"`
	Gender             string         `json:"gender" gorm:"default:null"`
	PhoneNumber        string         `json:"phone_number" gorm:"uniqueIndex;default:null"`
	LanguagePreference string         `json:"language_preference" gorm:"default:'en'"`
	TimeZone           string         `json:"time_zone" gorm:"default:'UTC'"`
	IsEmailVerified    bool           `json:"is_email_verified" gorm:"default:false"`
	IsPhoneVerified    bool           `json:"is_phone_verified" gorm:"default:false"`
	StatusID           uint           `json:"status_id" gorm:"not null;index;default:1"` // FK to user_statuses (1=active)
	LastLoginAt        *time.Time     `json:"last_login_at" gorm:"index"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Status              UserStatus           `json:"status" gorm:"foreignKey:StatusID"`
	AuthMethods         []UserAuthMethod     `json:"auth_methods,omitempty" gorm:"foreignKey:UserID"`
	OrganizationMembers []OrganizationMember `json:"organization_members,omitempty" gorm:"foreignKey:UserID"`
	Preferences         []UserPreference     `json:"preferences,omitempty" gorm:"foreignKey:UserID"`
	Locations           []UserLocation       `json:"locations,omitempty" gorm:"foreignKey:UserID"`

	// Ticket System Relationships
	OwnedProjects      []Project          `json:"owned_projects,omitempty" gorm:"foreignKey:OwnerID"`
	ProjectMemberships []ProjectMember    `json:"project_memberships,omitempty" gorm:"foreignKey:UserID"`
	AssignedTickets    []Ticket           `json:"assigned_tickets,omitempty" gorm:"foreignKey:AssigneeID"`
	ReportedTickets    []Ticket           `json:"reported_tickets,omitempty" gorm:"foreignKey:ReporterID"`
	TicketComments     []TicketComment    `json:"ticket_comments,omitempty" gorm:"foreignKey:UserID"`
	WatchedTickets     []Ticket           `json:"watched_tickets,omitempty" gorm:"many2many:ticket_watchers"`
	TimeLogs           []TimeLog          `json:"time_logs,omitempty" gorm:"foreignKey:UserID"`
	UploadedFiles      []TicketAttachment `json:"uploaded_files,omitempty" gorm:"foreignKey:UploadedBy"`
}

// Keep backward compatibility

// UserAuthMethod
type UserAuthMethod struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	AuthType     string         `json:"auth_type" gorm:"not null;index"` // password, oauth, sso
	AuthProvider *string        `json:"auth_provider" gorm:"index;default:null"`  // google, github, etc.
	ProviderID   *string        `json:"provider_id" gorm:"default:null"`          // OAuth ID
	IsPrimary    bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// For password auth
	PasswordHash string `json:"-" gorm:"default:null"`
	PasswordSalt string `json:"-" gorm:"default:null"`

	// For OAuth
	AccessToken  *string    `json:"-" gorm:"default:null"`
	RefreshToken *string    `json:"-" gorm:"default:null"`
	TokenExpiry  *time.Time `json:"-" gorm:"default:null"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Organization for multi-tenancy
type Organization struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	Slug        string         `json:"slug" gorm:"not null;uniqueIndex;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	LogoURL     string         `json:"logo_url"`
	PlanType    string         `json:"plan_type" gorm:"not null;default:'free';index"` // free, pro, enterprise
	StatusID    uint           `json:"status_id" gorm:"not null;index;default:1"`      // FK to organization_statuses (1=active)
	Settings    string         `json:"settings" gorm:"type:json"`                      // JSON settings
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Status   OrganizationStatus   `json:"status" gorm:"foreignKey:StatusID"`
	Members  []OrganizationMember `json:"members,omitempty" gorm:"foreignKey:OrganizationID"`
	Projects []Project            `json:"projects,omitempty" gorm:"foreignKey:OrganizationID"`
}

// OrganizationMember for role-based access
type OrganizationMember struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	OrganizationID uint       `json:"organization_id" gorm:"not null;index"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	Role           string     `json:"role" gorm:"not null;default:'member';index"` // owner, admin, member, guest
	StatusID       uint       `json:"status_id" gorm:"not null;index;default:1"`   // FK to member_statuses (1=active)
	InvitedAt      *time.Time `json:"invited_at"`
	JoinedAt       *time.Time `json:"joined_at"`
	InvitedBy      *uint      `json:"invited_by" gorm:"index"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	Status       MemberStatus `json:"status" gorm:"foreignKey:StatusID"`
	Organization Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	User         User         `json:"user" gorm:"foreignKey:UserID"`
	Inviter      *User        `json:"inviter,omitempty" gorm:"foreignKey:InvitedBy"`
}

// UserPreference for user settings
type UserPreference struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Key       string    `json:"key" gorm:"not null;index"` // notifications.email, theme.mode
	Value     string    `json:"value" gorm:"type:text"`
	Context   string    `json:"context" gorm:"index"` // global, org_123, project_456
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}


type UserLocation struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	LocationType  string         `json:"location_type" gorm:"not null;index"` // primary, shipping, billing, work
	AddressLine1  string         `json:"address_line1" gorm:"not null"`
	AddressLine2  string         `json:"address_line2"`
	City          string         `json:"city" gorm:"not null"`
	StateProvince string         `json:"state_province"`
	PostalCode    string         `json:"postal_code" gorm:"not null"`
	Country       string         `json:"country" gorm:"not null;index"`
	Latitude      *float64       `json:"latitude"`
	Longitude     *float64       `json:"longitude"`
	IsDefault     bool           `json:"is_default" gorm:"default:false;index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
