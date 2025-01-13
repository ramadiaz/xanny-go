package services

import (
	"layered-template/example/repositories"
	"layered-template/exceptions"
	"layered-template/mapper"
	"layered-template/models/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CompServicesImpl struct {
	repo     repositories.CompRepositories
	DB       *gorm.DB
	validate *validator.Validate
}

func NewComponentServices(compRepositories repositories.CompRepositories, db *gorm.DB, validate *validator.Validate) CompService {
	return &CompServicesImpl{
		repo:     compRepositories,
		DB:       db,
		validate: validate,
	}
}

func(s *CompServicesImpl) Create(ctx *gin.Context, data dto.ExampleInput) *exceptions.Exception {
	input := mapper.MapExampleInputToModel(data)

	return s.repo.Create(ctx, s.DB, input)
}