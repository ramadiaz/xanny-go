package services

import (
	"xanny-go-template/api/users/dto"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Create(ctx *gin.Context, data dto.Users) *exceptions.Exception
}