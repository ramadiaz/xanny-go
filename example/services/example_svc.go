package services

import (
	"layered-template/exceptions"
	"layered-template/models/dto"

	"github.com/gin-gonic/gin"
)

type CompService interface {
	Create(ctx *gin.Context, data dto.ExampleInput) *exceptions.Exception
}