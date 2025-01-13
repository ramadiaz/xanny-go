package routers

import (
	"xanny-go-template/example/controllers"

	"github.com/gin-gonic/gin"
)

func ExampleRoutes(r *gin.RouterGroup, exampleController controllers.CompControllers) {
	hotelGroup := r.Group("/example")
	{
		hotelGroup.POST("/create", exampleController.Create)
	}
}
