package cache

import (
	"time"

	"github.com/gin-gonic/gin"
)

// SetupCacheRoutes sets up cache-related routes
func SetupCacheRoutes(router *gin.RouterGroup, cacheService *CacheService) {
	controller := NewCacheController(cacheService)

	// Cache management routes
	cacheGroup := router.Group("/cache")
	{
		cacheGroup.GET("/stats", controller.GetCacheStats)
		cacheGroup.POST("/flush", controller.FlushCache)
		cacheGroup.POST("/invalidate", controller.InvalidatePattern)
		cacheGroup.POST("/set", controller.SetCacheValue)
		cacheGroup.GET("/:key", controller.GetCacheValue)
		cacheGroup.DELETE("/:key", controller.DeleteCacheValue)
	}

	// Example service routes with caching
	usersGroup := router.Group("/users")
	{
		usersGroup.GET("", controller.GetUsers)
		usersGroup.GET("/:id", controller.GetUser)
		usersGroup.PUT("/:id", controller.UpdateUser)
	}

	productsGroup := router.Group("/products")
	{
		productsGroup.GET("/:id", controller.GetProduct)
		productsGroup.GET("/:id/price", controller.GetProductPrice)
	}
}

// SetupCacheMiddleware sets up cache middleware for routes
func SetupCacheMiddleware(router *gin.Engine, cacheService *CacheService) {
	// Add cache middleware to specific routes
	cacheMiddleware := CacheMiddleware(&CacheMiddlewareOptions{
		Cache:      cacheService.cache,
		DefaultTTL: 5 * time.Minute,
		SkipCache: func(c *gin.Context) bool {
			// Skip caching for authenticated requests
			return c.GetHeader("Authorization") != ""
		},
	})

	// Apply cache middleware to public routes
	publicGroup := router.Group("/api/public")
	publicGroup.Use(cacheMiddleware)
	{
		// Add public routes here that should be cached
	}

	// Add cache control headers middleware
	cacheControlMiddleware := CacheControlMiddleware(10*time.Minute, true)
	router.Use(cacheControlMiddleware)
}
