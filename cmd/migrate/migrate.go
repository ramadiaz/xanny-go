package main

import (
	"xanny-go/models"
	"xanny-go/pkg/config"
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(&models.Users{}, &models.Clients{}, &models.RefreshToken{}, &models.BlacklistedToken{})
	if err != nil {
		panic("failed to migrate models: " + err.Error())
	}
}
