package services

import (
	"xanny-go-template/exceptions"
	"xanny-go-template/models/dto"

	"github.com/gin-gonic/gin"
)

type CompService interface {
	Create(ctx *gin.Context, data dto.ExampleInput) *exceptions.Exception
}
