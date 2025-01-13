package main

import (
	"xanny-go-template/pkg/config"
	"xanny-go-template/models"
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(
		&models.Client{},
		&models.Example{},
	)
	if err != nil {
		panic("failed to migrate models: " + err.Error())
	}
}
