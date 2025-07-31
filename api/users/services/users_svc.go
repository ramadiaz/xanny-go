package services

import (
	"xanny-go-template/api/users/dto"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Create(ctx *gin.Context, data dto.Users) *exceptions.Exception
	Login(ctx *gin.Context, email, password string) (accessToken, refreshToken string, err *exceptions.Exception)
	RefreshToken(ctx *gin.Context, refreshToken string) (accessToken string, err *exceptions.Exception)
	Logout(ctx *gin.Context, accessToken, refreshToken string) *exceptions.Exception
}
