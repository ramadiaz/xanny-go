package exceptions

import (
	"fmt"
	"net/http"
)

type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (err *Exception) Error() string {
	return fmt.Sprintf("Error %d: %s", err.Status, err.Message)
}

func NewException(status int, message string) *Exception {
	return &Exception{
		Status:  status,
		Message: message,
	}
}

func NewValidationException(err error) *Exception {
	return &Exception{
		Status:  http.StatusBadRequest,
		Message: err.Error(),
	}
}