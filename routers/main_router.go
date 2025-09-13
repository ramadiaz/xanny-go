package routers

import (
	"net/http"
	"xanny-go/injectors"
	"xanny-go/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func CompRouters(r *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) {
	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	// @Summary Health check
	// @Description Get the health status of all services
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} helpers.HealthCheck
	// @Failure 503 {object} helpers.HealthCheck
	// @Router /health [get]
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
