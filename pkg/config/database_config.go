package config

import (
	"fmt"
	"log"
	"os"
	"time"
	"xanny-go-template/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	GORMLogger "gorm.io/gorm/logger"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func InitDB() *gorm.DB {
	godotenv.Load()
	logger.Info("Initializing database connection...")
	start := time.Now()

	dbConfig := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}

	dbURI := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	connection, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: GORMLogger.Default.LogMode(GORMLogger.Error),
	})
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	logger.Info("Connected to database in %s", elapsed)

	return connection
}
