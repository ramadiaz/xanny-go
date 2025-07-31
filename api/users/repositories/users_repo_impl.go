package repositories

import (
	"xanny-go-template/models"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositoriesImpl struct {
}

func NewComponentRepository() CompRepositories {
	return &CompRepositoriesImpl{}
}

func (r *CompRepositoriesImpl) Create(ctx *gin.Context, tx *gorm.DB, data models.Users) *exceptions.Exception {
	result := tx.Create(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(tx, result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) CreateRefreshToken(ctx *gin.Context, tx *gorm.DB, token models.RefreshToken) *exceptions.Exception {
	if err := tx.Create(&token).Error; err != nil {
		return exceptions.ParseGormError(tx, err)
	}
	return nil
}

func (r *CompRepositoriesImpl) FindRefreshToken(ctx *gin.Context, tx *gorm.DB, token string) (*models.RefreshToken, *exceptions.Exception) {
	var refreshToken models.RefreshToken
	err := tx.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exceptions.ParseGormError(nil, err)
	}
	return &refreshToken, nil
}

func (r *CompRepositoriesImpl) DeleteRefreshToken(ctx *gin.Context, tx *gorm.DB, token string) *exceptions.Exception {
	if err := tx.Where("token = ?", token).Delete(&models.RefreshToken{}).Error; err != nil {
		return exceptions.ParseGormError(tx, err)
	}
	return nil
}

func (r *CompRepositoriesImpl) CreateBlacklistedToken(ctx *gin.Context, tx *gorm.DB, token models.BlacklistedToken) *exceptions.Exception {
	if err := tx.Create(&token).Error; err != nil {
		return exceptions.ParseGormError(tx, err)
	}
	return nil
}

func (r *CompRepositoriesImpl) FindBlacklistedToken(ctx *gin.Context, tx *gorm.DB, token string) (bool, *exceptions.Exception) {
	var blacklisted models.BlacklistedToken
	err := tx.Where("token = ?", token).First(&blacklisted).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, exceptions.ParseGormError(nil, err)
	}
	return true, nil
}