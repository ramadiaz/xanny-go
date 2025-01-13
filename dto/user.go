package dto

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	Password        string `json:"password,omitempty" validate:"required,min=8"`
	PhoneNumber     string `json:"phone_number" validate:"required,e164"`
	Country         string `json:"country" validate:"required"`
	Province        string `json:"province" validate:"required"`
	City            string `json:"city" validate:"required"`
	ZipCode         string `json:"zip_code" validate:"required,numeric,min=4,max=6"`
	CompleteAddress string `json:"complete_address" validate:"required,min=10"`
}