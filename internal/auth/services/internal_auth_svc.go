package services

import (
	"xanny-go/internal/auth/dto"
	"xanny-go/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Login(ctx *gin.Context, data dto.Login) (*string, *exceptions.Exception)
}
