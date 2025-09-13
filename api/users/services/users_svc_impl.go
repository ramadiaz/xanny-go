package services

import (
	"xanny-go/api/users/dto"
	"xanny-go/api/users/repositories"
	"xanny-go/models"
	"xanny-go/pkg/config"
	"xanny-go/pkg/exceptions"
	"xanny-go/pkg/helpers"
	"xanny-go/pkg/logger"

	emailDTO "xanny-go/emails/dto"
	emails "xanny-go/emails/services"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	userUUID := uuid.NewString()

	err = s.repo.Create(ctx, tx, models.Users{
		UUID:           userUUID,
		Name:           data.Name,
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return err
	}

	go func() {
		token, err := s.CreateVerificationToken(ctx, userUUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		err = emails.VerificationEmail(emailDTO.EmailVerification{
			Email:           data.Email,
			Name:            data.Name,
			VerificationURL: config.GetFrontendURL() + "/auth/verify?token=" + *token,
			SupportEmail:    "support@xanware.id",
		})
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}()

	return nil
}

func (s *CompServicesImpl) CreateVerificationToken(ctx *gin.Context, userUUID string) (*string, *exceptions.Exception) {
	token := helpers.GenerateRandomString(32)

	err := s.repo.CreateVerificationToken(ctx, s.DB, models.VerificationToken{
		Token:     token,
		UserUUID:  userUUID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *CompServicesImpl) ResendVerificationEmail(ctx *gin.Context, email string) *exceptions.Exception {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	user, err := s.repo.FindByEmail(ctx, tx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return exceptions.NewException(404, "User not found")
	}

	if user.IsEmailVerified {
		return exceptions.NewException(400, "Email is already verified")
	}

	existToken, err := s.repo.FindVerificationTokenByUserUUID(ctx, tx, user.UUID)
	if err == nil {
		if existToken.ExpiresAt.After(time.Now()) {
			return exceptions.NewException(400, "Verification token is still valid")
		}
	}

	token, err := s.CreateVerificationToken(ctx, user.UUID)
	if err != nil {
		return err
	}

	err = emails.VerificationEmail(emailDTO.EmailVerification{
		Email:           user.Email,
		Name:            user.Name,
		VerificationURL: config.GetFrontendURL() + "/auth/verify?token=" + *token,
		SupportEmail:    "support@xanware.id",
	})
	if err != nil {
		return exceptions.NewException(500, "Failed to send verification email")
	}

	return nil
}

func (s *CompServicesImpl) VerificationEmail(ctx *gin.Context, token string) *exceptions.Exception {
	tx := s.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	tokenData, err := s.repo.FindVerificationToken(ctx, tx, token)
	if err != nil {
		return err
	}

	if tokenData == nil {
		return exceptions.NewException(400, "Invalid or expired verification token")
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return exceptions.NewException(400, "Verification token has expired")
	}

	err = s.repo.Update(ctx, tx, models.Users{
		UUID:            tokenData.UserUUID,
		IsEmailVerified: true,
	})
	if err != nil {
		return err
	}

	err = s.repo.DeleteVerificationToken(ctx, tx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *CompServicesImpl) Login(ctx *gin.Context, email, password string) (accessToken, refreshToken string, err *exceptions.Exception) {
	user, err := s.repo.FindByEmail(ctx, s.DB, email)
	if err != nil {
		return "", "", err
	}

	if hashErr := helpers.CheckPasswordHash(password, user.HashedPassword); hashErr != nil {
		return "", "", exceptions.NewException(401, "Invalid email or password")
	}

	if !user.IsEmailVerified {
		return "", "", exceptions.NewException(401, "Email is not verified")
	}

	jwtSecret := config.GetJWTSecret()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["email"] = user.Email
	claims["is_email_verified"] = user.IsEmailVerified
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessTokenStr, signErr := token.SignedString([]byte(jwtSecret))
	if signErr != nil {
		return "", "", exceptions.NewException(500, "Failed to generate access token")
	}

	refreshTokenRaw := helpers.GenerateRandomString(64)
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7)
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

	user, err := s.repo.FindByUUID(ctx, s.DB, tokenModel.UserUUID)
	if err != nil {
		return "", err
	}

	jwtSecret := config.GetJWTSecret()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["email"] = user.Email
	claims["is_email_verified"] = user.IsEmailVerified
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessTokenStr, signErr := token.SignedString([]byte(jwtSecret))
	if signErr != nil {
		return "", exceptions.NewException(500, "Failed to generate access token")
	}

	return accessTokenStr, nil
}

func (s *CompServicesImpl) Logout(ctx *gin.Context, accessToken, refreshToken string) *exceptions.Exception {
	tx := s.DB.Begin()

	claims := jwt.MapClaims{}
	jwtSecret := config.GetJWTSecret()
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err == nil {
		exp, _ := claims["exp"].(float64)
		helpers.SetBlacklistedToken(accessToken, time.Unix(int64(exp), 0))
	}

	s.repo.DeleteRefreshToken(ctx, tx, refreshToken)
	tx.Commit()
	return nil
}
