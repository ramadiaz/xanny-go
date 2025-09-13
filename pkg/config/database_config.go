package config

import (
	"fmt"
	"log"
	"time"
	"xanny-go/pkg/logger"

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
	logger.Info("Initializing database connection...")
	start := time.Now()

	dbConfig := DBConfig{
		User:     GetDBUser(),
		Password: GetDBPassword(),
		Host:     GetDBHost(),
		Port:     GetDBPort(),
		Name:     GetDBName(),
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
