package controllers

import (
	"net/http"
	"xanny-go-template/api/users/dto"
	"xanny-go-template/api/users/services"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompControllersImpl struct {
	services services.CompServices
}

func NewCompController(compServices services.CompServices) CompControllers {
	return &CompControllersImpl{
		services: compServices,
	}
}

// Create godoc
// @Summary Create a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.Users true "User registration data"
// @Success 201 {object} dto.Response
// @Failure 400 {object} exceptions.Exception
// @Failure 409 {object} exceptions.Exception
// @Router /user/create [post]
func (h *CompControllersImpl) Create(ctx *gin.Context) {
	var data dto.Users

	jsonErr := ctx.ShouldBindJSON(&data)
	if jsonErr != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, exceptions.ErrBadRequest))
		return
	}

	err := h.services.Create(ctx, data)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: "success",
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} exceptions.Exception
// @Failure 401 {object} exceptions.Exception
// @Router /user/login [post]
func (h *CompControllersImpl) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, "Invalid request body"))
		return
	}
	accessToken, refreshToken, err := h.services.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param refresh body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} exceptions.Exception
// @Failure 401 {object} exceptions.Exception
// @Router /user/refresh [post]
func (h *CompControllersImpl) Refresh(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, "Invalid request body"))
		return
	}
	accessToken, err := h.services.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate tokens
// @Tags users
// @Accept json
// @Produce json
// @Param logout body dto.LogoutRequest true "Logout request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} exceptions.Exception
// @Router /user/logout [post]
func (h *CompControllersImpl) Logout(ctx *gin.Context) {
	var req dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, "Invalid request body"))
		return
	}
	err := h.services.Logout(ctx, req.AccessToken, req.RefreshToken)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "logout success"})
}
