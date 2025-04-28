package repositories

import (
	"xanny-go-template/models"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data models.Users) *exceptions.Exception 
}