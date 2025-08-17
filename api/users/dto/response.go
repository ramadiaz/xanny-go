package dto

// Response represents a generic API response
type Response struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"success"`
	Body    interface{} `json:"body,omitempty"`
}

// UserOutput represents user data in responses
type UserOutput struct {
	UUID            string `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email           string `json:"email" example:"user@example.com"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Name            string `json:"name" example:"John Doe"`
}
