package schema

type LocationRequest struct {
	LocationType string `json:"location_type"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`

	Location *LocationRequest `json:"location"`
}
