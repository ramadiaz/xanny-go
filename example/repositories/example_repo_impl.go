package repositories

import (
	"xanny-go-template/exceptions"
	"xanny-go-template/models/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositoriesImpl struct {
}

func NewComponentRepository() CompRepositories {
	return &CompRepositoriesImpl{}
}

func (r *CompRepositoriesImpl) Create(ctx *gin.Context, tx *gorm.DB, data database.Example) *exceptions.Exception {
	return nil
}
