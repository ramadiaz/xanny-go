package dto

// Users represents user registration request
type Users struct {
	Email     string `json:"email" example:"user@example.com" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Passoword string `json:"password" example:"password123" validate:"required,min=6"`
}
// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"password123" binding:"required"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
}
