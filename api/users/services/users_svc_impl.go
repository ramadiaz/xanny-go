package services

import (
	"xanny-go-template/api/users/dto"
	"xanny-go-template/api/users/repositories"
	"xanny-go-template/models"
	"xanny-go-template/pkg/exceptions"
	"xanny-go-template/pkg/helpers"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CompServicesImpl struct {
	repo     repositories.CompRepositories
	DB       *gorm.DB
	validate *validator.Validate
}

func NewComponentServices(compRepositories repositories.CompRepositories, db *gorm.DB, validate *validator.Validate) CompServices {
	return &CompServicesImpl{
		repo:     compRepositories,
		DB:       db,
		validate: validate,
	}
}

func (s *CompServicesImpl) Create(ctx *gin.Context, data dto.Users) *exceptions.Exception {
	validateErr := s.validate.Struct(data)
	if validateErr != nil {
		return exceptions.NewValidationException(validateErr)
	}

	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	hashedPassword, err := helpers.HashPassword(data.Passoword)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, tx, models.Users{
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *CompServicesImpl) Login(ctx *gin.Context, email, password string) (accessToken, refreshToken string, err *exceptions.Exception) {
	var user models.Users
	errFind := s.DB.Where("email = ?", email).First(&user).Error
	if errFind != nil {
		return "", "", exceptions.NewException(401, "Invalid email or password")
	}

	if hashErr := helpers.CheckPasswordHash(password, user.HashedPassword); hashErr != nil {
		return "", "", exceptions.NewException(401, "Invalid email or password")
	}

	// Generate Access Token
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessTokenStr, signErr := token.SignedString([]byte(jwtSecret))
	if signErr != nil {
		return "", "", exceptions.NewException(500, "Failed to generate access token")
	}

	// Generate Refresh Token
	refreshTokenRaw := helpers.GenerateRandomString(64)
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7) // 7 days
	refreshTokenModel := models.RefreshToken{
		UserUUID:  user.UUID,
		Token:     refreshTokenRaw,
		ExpiresAt: refreshTokenExp,
		CreatedAt: time.Now(),
	}
	tx := s.DB.Begin()
	if err := s.repo.CreateRefreshToken(ctx, tx, refreshTokenModel); err != nil {
		tx.Rollback()
		return "", "", exceptions.NewException(500, "Failed to save refresh token")
	}
	tx.Commit()

	return accessTokenStr, refreshTokenRaw, nil
}

func (s *CompServicesImpl) RefreshToken(ctx *gin.Context, refreshToken string) (accessToken string, err *exceptions.Exception) {
	tokenModel, errFind := s.repo.FindRefreshToken(ctx, s.DB, refreshToken)
	if errFind != nil {
		return "", errFind
	}
	if tokenModel == nil || tokenModel.ExpiresAt.Before(time.Now()) {
		return "", exceptions.NewException(401, "Refresh token expired or not found")
	}

	// Find user
	var user models.Users
	errUser := s.DB.Where("uuid = ?", tokenModel.UserUUID).First(&user).Error
	if errUser != nil {
		return "", exceptions.NewException(401, "User not found")
	}

	// Generate new Access Token
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessTokenStr, signErr := token.SignedString([]byte(jwtSecret))
	if signErr != nil {
		return "", exceptions.NewException(500, "Failed to generate access token")
	}

	return accessTokenStr, nil
}

func (s *CompServicesImpl) Logout(ctx *gin.Context, accessToken, refreshToken string) *exceptions.Exception {
	tx := s.DB.Begin()
	// Blacklist access token
	claims := jwt.MapClaims{}
	jwtSecret := os.Getenv("JWT_SECRET")
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err == nil {
		exp, _ := claims["exp"].(float64)
		blacklisted := models.BlacklistedToken{
			Token:     accessToken,
			ExpiresAt: time.Unix(int64(exp), 0),
			CreatedAt: time.Now(),
		}
		s.repo.CreateBlacklistedToken(ctx, tx, blacklisted)
	}
	// Delete refresh token
	s.repo.DeleteRefreshToken(ctx, tx, refreshToken)
	tx.Commit()
	return nil
}
