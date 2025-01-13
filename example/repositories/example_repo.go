package repositories

import (
	"layered-template/exceptions"
	"layered-template/models/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data database.Example) *exceptions.Exception
}