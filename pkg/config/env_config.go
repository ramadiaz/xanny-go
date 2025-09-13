package config

import (
	"os"

	"xanny-go/pkg/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_USER         string
	DB_PASSWORD     string
	DB_HOST         string
	DB_PORT         string
	DB_NAME         string
	PORT            string
	JWT_SECRET      string
	INTERNAL_SECRET string
	ENVIRONMENT     string
	ADMIN_USERNAME  string
	ADMIN_PASSWORD  string
	REDIS_ADDR      string
	REDIS_PASS      string
	FRONTEND_URL    string
	SMTP_EMAIL      string
	SMTP_PASSWORD   string
	SMTP_SERVER     string
	SMTP_PORT       string
}

var globalConfig *Config

func InitConfig() {
	godotenv.Load()
	config := &Config{
		DB_USER:         getEnv("DB_USER"),
		DB_PASSWORD:     getEnv("DB_PASSWORD"),
		DB_HOST:         getEnv("DB_HOST"),
		DB_PORT:         getEnv("DB_PORT"),
		DB_NAME:         getEnv("DB_NAME"),
		PORT:            getEnv("PORT"),
		JWT_SECRET:      getEnv("JWT_SECRET"),
		INTERNAL_SECRET: getEnv("INTERNAL_SECRET"),
		ENVIRONMENT:     getEnv("ENVIRONMENT"),
		ADMIN_USERNAME:  getEnv("ADMIN_USERNAME"),
		ADMIN_PASSWORD:  getEnv("ADMIN_PASSWORD"),
		REDIS_ADDR:      getEnv("REDIS_ADDR"),
		REDIS_PASS:      getEnv("REDIS_PASS"),
		FRONTEND_URL:    getEnv("FRONTEND_URL"),
		SMTP_EMAIL:      getEnv("SMTP_EMAIL"),
		SMTP_PASSWORD:   getEnv("SMTP_PASSWORD"),
		SMTP_SERVER:     getEnv("SMTP_SERVER"),
		SMTP_PORT:       getEnv("SMTP_PORT"),
	}

	globalConfig = config
	logger.Info("Configuration initialized successfully")
}

func GetConfig() *Config {
	if globalConfig == nil {
		logger.PanicError("Configuration not initialized. Call InitConfig() first.")
	}
	return globalConfig
}

func GetDBUser() string         { return GetConfig().DB_USER }
func GetDBPassword() string     { return GetConfig().DB_PASSWORD }
func GetDBHost() string         { return GetConfig().DB_HOST }
func GetDBPort() string         { return GetConfig().DB_PORT }
func GetDBName() string         { return GetConfig().DB_NAME }
func GetPort() string           { return GetConfig().PORT }
func GetJWTSecret() string      { return GetConfig().JWT_SECRET }
func GetInternalSecret() string { return GetConfig().INTERNAL_SECRET }
func GetEnvironment() string    { return GetConfig().ENVIRONMENT }
func GetAdminUsername() string  { return GetConfig().ADMIN_USERNAME }
func GetAdminPassword() string  { return GetConfig().ADMIN_PASSWORD }
func GetRedisAddr() string      { return GetConfig().REDIS_ADDR }
func GetRedisPass() string      { return GetConfig().REDIS_PASS }
func GetFrontendURL() string    { return GetConfig().FRONTEND_URL }
func GetSMTPEmail() string      { return GetConfig().SMTP_EMAIL }
func GetSMTPPassword() string   { return GetConfig().SMTP_PASSWORD }
func GetSMTPServer() string     { return GetConfig().SMTP_SERVER }
func GetSMTPPort() string       { return GetConfig().SMTP_PORT }

func IsProduction() bool  { return GetEnvironment() == "production" }
func IsDevelopment() bool { return GetEnvironment() == "development" }

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.PanicError("Environment variable %s is required but not set", key)
	}
	return value
}
