package repositories

import (
	"xanny-go-template/models"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data models.Users) *exceptions.Exception
	FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) (*models.Users, *exceptions.Exception)
	FindByEmail(ctx *gin.Context, tx *gorm.DB, email string) (*models.Users, *exceptions.Exception)
	Update(ctx *gin.Context, tx *gorm.DB, data models.Users) *exceptions.Exception
	CreateRefreshToken(ctx *gin.Context, tx *gorm.DB, token models.RefreshToken) *exceptions.Exception
	FindRefreshToken(ctx *gin.Context, tx *gorm.DB, token string) (*models.RefreshToken, *exceptions.Exception)
	DeleteRefreshToken(ctx *gin.Context, tx *gorm.DB, token string) *exceptions.Exception
	CreateBlacklistedToken(ctx *gin.Context, tx *gorm.DB, token models.BlacklistedToken) *exceptions.Exception
	FindBlacklistedToken(ctx *gin.Context, tx *gorm.DB, token string) (bool, *exceptions.Exception)
	CreateVerificationToken(ctx *gin.Context, tx *gorm.DB, token models.VerificationToken) *exceptions.Exception
	FindVerificationToken(ctx *gin.Context, tx *gorm.DB, token string) (*models.VerificationToken, *exceptions.Exception)
	FindVerificationTokenByUserUUID(ctx *gin.Context, tx *gorm.DB, userUUID string) (*models.VerificationToken, *exceptions.Exception)
	DeleteVerificationToken(ctx *gin.Context, tx *gorm.DB, token string) *exceptions.Exception
}
