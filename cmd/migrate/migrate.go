package main

import (
	"xanny-go-template/pkg/config"
	"xanny-go-template/models"
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(
		&models.Clients{},
		&models.Users{},
	)
	if err != nil {
		panic("failed to migrate models: " + err.Error())
	}
}
