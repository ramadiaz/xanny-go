package routers

import (
	"net/http"
	"xanny-go-template/injectors"
	"xanny-go-template/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CompRouters(r *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) {
	r.GET("/health", func(ctx *gin.Context) {
		health := helpers.PerformHealthCheck(db)

		statusCode := http.StatusOK
		if health.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		ctx.JSON(statusCode, health)
	})

	userController := injectors.InitializeUserController(db, validate)

	UserRoutes(r, userController)
}
