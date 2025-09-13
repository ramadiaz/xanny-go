package repositories

import (
	"xanny-go/models"
	"xanny-go/pkg/exceptions"

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

func (r *CompRepositoriesImpl) FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) (*models.Users, *exceptions.Exception) {
	var user models.Users
	err := tx.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, exceptions.ParseGormError(tx, err)
	}
	return &user, nil
}

func (r *CompRepositoriesImpl) FindByEmail(ctx *gin.Context, tx *gorm.DB, email string) (*models.Users, *exceptions.Exception) {
	var user models.Users
	err := tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, exceptions.ParseGormError(tx, err)
	}
	return &user, nil
}

func (r *CompRepositoriesImpl) Update(ctx *gin.Context, tx *gorm.DB, data models.Users) *exceptions.Exception {
	result := tx.Where("uuid = ?", data.UUID).Updates(&data)
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

func (r *CompRepositoriesImpl) CreateVerificationToken(ctx *gin.Context, tx *gorm.DB, token models.VerificationToken) *exceptions.Exception {
	if err := tx.Create(&token).Error; err != nil {
		return exceptions.ParseGormError(tx, err)
	}
	return nil
}

func (r *CompRepositoriesImpl) FindVerificationToken(ctx *gin.Context, tx *gorm.DB, token string) (*models.VerificationToken, *exceptions.Exception) {
	var verification models.VerificationToken
	err := tx.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return nil, exceptions.ParseGormError(tx, err)
	}
	return &verification, nil
}

func (r *CompRepositoriesImpl) FindVerificationTokenByUserUUID(ctx *gin.Context, tx *gorm.DB, userUUID string) (*models.VerificationToken, *exceptions.Exception) {
	var verification models.VerificationToken
	err := tx.Where("user_uuid = ?", userUUID).First(&verification).Error
	if err != nil {
		return nil, exceptions.ParseGormError(tx, err)
	}
	return &verification, nil
}

func (r *CompRepositoriesImpl) DeleteVerificationToken(ctx *gin.Context, tx *gorm.DB, token string) *exceptions.Exception {
	if err := tx.Where("token = ?", token).Delete(&models.VerificationToken{}).Error; err != nil {
		return exceptions.ParseGormError(tx, err)
	}
	return nil
}
