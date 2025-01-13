package main

import (
	"xanny-go-template/config"
	"xanny-go-template/models/database"
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(
		&database.Client{},
		&database.Example{},
	)
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}
