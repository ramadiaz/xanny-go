package controllers

import "github.com/gin-gonic/gin"

type CompControllers interface {
	Create(ctx *gin.Context)
	ResendVerificationEmail(ctx *gin.Context)
	VerificationEmail(ctx *gin.Context)
	Login(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	Logout(ctx *gin.Context)
}
