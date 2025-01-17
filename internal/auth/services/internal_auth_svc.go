package services

import (
	"xanny-go-template/internal/auth/dto"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Login(ctx *gin.Context, data dto.Login) (*string, *exceptions.Exception)
}
