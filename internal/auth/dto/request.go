package dto

type Login struct {
	Username string `validate:"required,min=4"`
	Password string `validate:"required,min=8"`
}