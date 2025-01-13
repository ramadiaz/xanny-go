package routers

import (
	"net/http"
	"xanny-go-template/config"
	"xanny-go-template/injectors"
	"xanny-go-template/middleware"
	"xanny-go-template/models/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CompRouters(api *gin.RouterGroup) {
	db := config.InitDB()
	validate := validator.New(validator.WithRequiredStructEnabled())

	api.Use(middleware.ClientTracker(db))
	api.Use(middleware.GzipResponseMiddleware())

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.Response{
			Status: http.StatusOK,
			Data:   "pong",
		})
	})

	exampleController := injectors.InitializeExampleController(db, validate)

	ExampleRoutes(api, exampleController)
}
