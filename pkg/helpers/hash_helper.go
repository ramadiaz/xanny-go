package helpers

import (
	"net/http"
	"xanny-go-template/pkg/exceptions"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, *exceptions.Exception) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", exceptions.NewException(http.StatusInternalServerError, exceptions.ErrCredentialsHash)
    }

    return string(bytes), nil
}

func CheckPasswordHash(password, hash string) *exceptions.Exception {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        return exceptions.NewException(http.StatusUnauthorized, exceptions.ErrInvalidCredentials)
    }

    return nil
}