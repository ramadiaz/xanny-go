package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConfig := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}

	if dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.Name == "" {
		log.Fatal("Database configuration is missing")
	}

	dbURI := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",  dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)


	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
