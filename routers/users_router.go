package routers

import (
	"xanny-go-template/api/users/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, userController controllers.CompControllers) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/create", userController.Create)
		userGroup.POST("/login", userController.Login)
		userGroup.POST("/refresh", userController.Refresh)
		userGroup.POST("/logout", userController.Logout)
	}
}
