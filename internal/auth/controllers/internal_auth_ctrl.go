package controllers

import "github.com/gin-gonic/gin"

type CompControllers interface {
	Login(ctx *gin.Context)
}