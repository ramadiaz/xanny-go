package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func InitDB() *sql.DB {
	godotenv.Load()

	dbConfig := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}

	dbURI := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
