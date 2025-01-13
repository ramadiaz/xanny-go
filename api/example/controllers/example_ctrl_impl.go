package controllers

import (
	"net/http"
	"xanny-go-template/api/example/services"
	"xanny-go-template/dto"
	"xanny-go-template/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompControllersImpl struct {
	services services.CompService
}

func NewCompController(compServices services.CompService) CompControllers {
	return &CompControllersImpl{
		services: compServices,
	}
}

func (h *CompControllersImpl) Create(ctx *gin.Context) {
	var exampleInput dto.ExampleInput

	if jsonErr := ctx.ShouldBindJSON(&exampleInput); jsonErr != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, jsonErr.Error()))
		return
	}

	err := h.services.Create(ctx, exampleInput)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: "data created successfully",
	})
}
