package services

import (
	"xanny-go-template/dto"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompService interface {
	Create(ctx *gin.Context, data dto.ExampleInput) *exceptions.Exception
}
