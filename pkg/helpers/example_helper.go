package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func SayHi() string {
	return "Hi!"
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
