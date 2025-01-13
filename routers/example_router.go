package routers

import (
	"layered-template/example/controllers"

	"github.com/gin-gonic/gin"
)

func ExampleRoutes(r *gin.RouterGroup, exampleController controllers.CompControllers) {
	hotelGroup := r.Group("/example")
	{	
		hotelGroup.GET("/create", exampleController.Create)
	}
}