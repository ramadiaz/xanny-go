package repositories

import (
	"database/sql"
	"log"
	"layered-template/config"
)

type CompRepository interface{}

type compRepository struct {
	DB *sql.DB
}

func NewComponentRepository(DB *sql.DB) *compRepository {
	db := config.InitDB()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS client_track (
			id BIGSERIAL PRIMARY KEY NOT NULL,
			ip VARCHAR(255),
			browser VARCHAR(255),
			version VARCHAR(255),
			os VARCHAR(255),
			device VARCHAR(255),
			origin VARCHAR(255),
			api VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return &compRepository{
		DB: db,
	}
}
