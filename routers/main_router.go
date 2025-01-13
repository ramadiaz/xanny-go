package routers

import (
	"net/http"
	"xanny-go-template/pkg/config"
	"xanny-go-template/injectors"
	"xanny-go-template/pkg/middleware"
	"xanny-go-template/dto"

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
			Message:   "pong",
		})
	})

	exampleController := injectors.InitializeExampleController(db, validate)

	ExampleRoutes(api, exampleController)
}
