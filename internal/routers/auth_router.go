package routers

import (
	"xanny-go/internal/auth/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, internalAuthController controllers.CompControllers) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", internalAuthController.Login)
	}
}
