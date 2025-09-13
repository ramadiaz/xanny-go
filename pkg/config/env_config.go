package config

import (
	"os"
	"reflect"
	"xanny-go/pkg/logger"
)

type Env struct {
	DB_USER        string
	DB_PASSWORD    string
	DB_HOST        string
	DB_PORT        string
	DB_NAME        string
	PORT           string
	JWT_SECRET     string
	ENVIRONMENT    string
	ADMIN_USERNAME string
	ADMIN_PASSWORD string
	// FONNTE_API_KEY string
}

func InitEnvCheck() {
	logger.Info("Checking environment variables...")
	environment := Env{
		DB_USER:        os.Getenv("DB_USER"),
		DB_PASSWORD:    os.Getenv("DB_PASSWORD"),
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_NAME:        os.Getenv("DB_NAME"),
		PORT:           os.Getenv("PORT"),
		JWT_SECRET:     os.Getenv("JWT_SECRET"),
		ENVIRONMENT:    os.Getenv("ENVIRONMENT"),
		ADMIN_USERNAME: os.Getenv("ADMIN_USERNAME"),
		ADMIN_PASSWORD: os.Getenv("ADMIN_PASSWORD"),
		// FONNTE_API_KEY: os.Getenv("FONNTE_API_KEY"),
	}

	isEmpty, emptyFields := checkEmptyFields(environment)
	if isEmpty {
		logger.PanicError("Missing environment variables: %v", emptyFields)
	} else {
		logger.Info("Environment variables are set!")
	}
}

func checkEmptyFields(env Env) (bool, []string) {
	v := reflect.ValueOf(env)
	typeOfEnv := v.Type()
	var emptyFields []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			emptyFields = append(emptyFields, typeOfEnv.Field(i).Name)
		}
	}

	return len(emptyFields) > 0, emptyFields
}
