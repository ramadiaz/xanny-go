package routers

import (
	"net/http"
	"xanny-go-template/dto"
	"xanny-go-template/injectors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CompRouters(r *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.Response{
			Status: http.StatusOK,
			Message:   "pong",
		})
	})

	exampleController := injectors.InitializeExampleController(db, validate)

	ExampleRoutes(r, exampleController)
}
