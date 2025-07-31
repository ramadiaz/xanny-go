package cache

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IntegrationExample shows how to integrate cache with existing project
type IntegrationExample struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewIntegrationExample creates a new integration example
func NewIntegrationExample(db *gorm.DB, cacheService *CacheService) *IntegrationExample {
	return &IntegrationExample{
		db:           db,
		cacheService: cacheService,
	}
}

// SetupCacheIntegration demonstrates how to integrate cache with existing services
func SetupCacheIntegration(db *gorm.DB, router *gin.Engine) error {
	// 1. Initialize cache from environment
	cache, err := NewCacheFromEnv()
	if err != nil {
		log.Printf("Warning: Failed to initialize cache: %v", err)
		log.Println("Continuing without cache...")
		return nil
	}
	defer cache.Close()

	// 2. Create cache service
	cacheService := NewCacheService(cache)

	// 3. Create integration example
	integration := NewIntegrationExample(db, cacheService)

	// 4. Setup cache routes
	api := router.Group("/api")
	SetupCacheRoutes(api, cacheService)

	// 5. Setup cache middleware for public routes
	SetupCacheMiddleware(router, cacheService)

	// 6. Add cache health check endpoint
	router.GET("/health/cache", integration.CacheHealthHandler)

	log.Println("Cache integration setup completed")
	return nil
}

// CacheHealthHandler handles cache health checks
func (ie *IntegrationExample) CacheHealthHandler(c *gin.Context) {
	healthChecker := NewCacheHealthChecker(ie.cacheService.cache)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := healthChecker.CheckHealth(ctx)
	if err != nil {
		c.JSON(503, gin.H{
			"status":  "unhealthy",
			"cache":   "down",
			"error":   err.Error(),
			"message": "Cache health check failed",
		})
		return
	}

	stats := healthChecker.GetCacheStats(ctx)
	c.JSON(200, gin.H{
		"status": "healthy",
		"cache":  "up",
		"stats":  stats,
	})
}

// ExampleUserServiceWithCache demonstrates how to add cache to existing user service
type ExampleUserServiceWithCache struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewExampleUserServiceWithCache creates a new user service with cache
func NewExampleUserServiceWithCache(db *gorm.DB, cacheService *CacheService) *ExampleUserServiceWithCache {
	return &ExampleUserServiceWithCache{
		db:           db,
		cacheService: cacheService,
	}
}

// GetUserWithCache demonstrates cached user retrieval
func (us *ExampleUserServiceWithCache) GetUserWithCache(ctx context.Context, userID string) (map[string]interface{}, error) {
	var user map[string]interface{}

	cacheKey := "user:" + userID
	err := us.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &user, func() (interface{}, time.Duration, error) {
		// This would be your actual database query
		// For example:
		// var userModel models.User
		// result := us.db.WithContext(ctx).First(&userModel, userID)
		// if result.Error != nil {
		//     return nil, 0, result.Error
		// }

		// Simulate database query
		user = map[string]interface{}{
			"id":    userID,
			"name":  "John Doe",
			"email": "john@example.com",
		}

		return user, 10 * time.Minute, nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserWithCache demonstrates cached user update with invalidation
func (us *ExampleUserServiceWithCache) UpdateUserWithCache(ctx context.Context, userID string, updates map[string]interface{}) error {
	// This would be your actual database update
	// For example:
	// result := us.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	// if result.Error != nil {
	//     return result.Error
	// }

	// Invalidate user cache
	cacheKey := "user:" + userID
	if err := us.cacheService.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("Warning: Failed to invalidate user cache: %v", err)
	}

	// Invalidate user list caches
	if err := us.cacheService.cache.InvalidatePattern(ctx, "users:*"); err != nil {
		log.Printf("Warning: Failed to invalidate user list cache: %v", err)
	}

	return nil
}

// ExampleProductServiceWithCache demonstrates product service with different TTL strategies
type ExampleProductServiceWithCache struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewExampleProductServiceWithCache creates a new product service with cache
func NewExampleProductServiceWithCache(db *gorm.DB, cacheService *CacheService) *ExampleProductServiceWithCache {
	return &ExampleProductServiceWithCache{
		db:           db,
		cacheService: cacheService,
	}
}

// GetProductWithCache demonstrates long-term product caching
func (ps *ExampleProductServiceWithCache) GetProductWithCache(ctx context.Context, productID string) (map[string]interface{}, error) {
	var product map[string]interface{}

	cacheKey := "product:" + productID
	err := ps.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &product, func() (interface{}, time.Duration, error) {
		// Simulate database query for product details
		product = map[string]interface{}{
			"id":          productID,
			"name":        "Sample Product",
			"description": "A sample product description",
			"category":    "Electronics",
		}

		// Products can be cached longer since they don't change frequently
		return product, 1 * time.Hour, nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductPriceWithCache demonstrates short-term price caching
func (ps *ExampleProductServiceWithCache) GetProductPriceWithCache(ctx context.Context, productID string) (float64, error) {
	var price float64

	cacheKey := "product_price:" + productID
	err := ps.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &price, func() (interface{}, time.Duration, error) {
		// Simulate database query for price
		price = 99.99

		// Prices can change frequently, so cache for shorter time
		return price, 5 * time.Minute, nil
	})

	if err != nil {
		return 0, err
	}

	return price, nil
}

// ExampleConfigurationService demonstrates configuration caching
type ExampleConfigurationService struct {
	cacheService *CacheService
}

// NewExampleConfigurationService creates a new configuration service with cache
func NewExampleConfigurationService(cacheService *CacheService) *ExampleConfigurationService {
	return &ExampleConfigurationService{
		cacheService: cacheService,
	}
}

// GetConfigurationWithCache demonstrates configuration caching
func (cs *ExampleConfigurationService) GetConfigurationWithCache(ctx context.Context, configKey string) (map[string]interface{}, error) {
	var config map[string]interface{}

	cacheKey := "config:" + configKey
	err := cs.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &config, func() (interface{}, time.Duration, error) {
		// Simulate configuration loading
		config = map[string]interface{}{
			"key":          configKey,
			"value":        "config_value",
			"last_updated": time.Now(),
		}

		// Configuration can be cached for a very long time
		return config, 24 * time.Hour, nil
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}

// RefreshConfiguration demonstrates configuration refresh
func (cs *ExampleConfigurationService) RefreshConfiguration(ctx context.Context, configKey string) error {
	// Invalidate configuration cache
	cacheKey := "config:" + configKey
	return cs.cacheService.cache.Delete(ctx, cacheKey)
}

// ExampleRateLimitService demonstrates rate limiting with cache
type ExampleRateLimitService struct {
	cacheService *CacheService
}

// NewExampleRateLimitService creates a new rate limit service with cache
func NewExampleRateLimitService(cacheService *CacheService) *ExampleRateLimitService {
	return &ExampleRateLimitService{
		cacheService: cacheService,
	}
}

// CheckRateLimit demonstrates rate limiting with cache
func (rls *ExampleRateLimitService) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Try to increment the counter
	if redisCache, ok := rls.cacheService.cache.(*RedisCache); ok {
		current, err := redisCache.Increment(ctx, "rate_limit:"+key, 1)
		if err != nil {
			return false, err
		}

		// Set expiration if this is the first request
		if current == 1 {
			redisCache.Set(ctx, "rate_limit:"+key, []byte("1"), window)
		}

		return current <= int64(limit), nil
	}

	// Fallback for non-Redis caches
	return true, nil
}

// ExampleSessionService demonstrates session management with cache
type ExampleSessionService struct {
	cacheService *CacheService
}

// NewExampleSessionService creates a new session service with cache
func NewExampleSessionService(cacheService *CacheService) *ExampleSessionService {
	return &ExampleSessionService{
		cacheService: cacheService,
	}
}

// StoreSession demonstrates session storage with cache
func (ss *ExampleSessionService) StoreSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	cacheKey := "session:" + sessionID
	return ss.cacheService.helper.SetJSON(ctx, cacheKey, data, 30*time.Minute)
}

// GetSession demonstrates session retrieval with cache
func (ss *ExampleSessionService) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	var session map[string]interface{}

	cacheKey := "session:" + sessionID
	err := ss.cacheService.helper.GetJSON(ctx, cacheKey, &session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// DeleteSession demonstrates session deletion
func (ss *ExampleSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	cacheKey := "session:" + sessionID
	return ss.cacheService.cache.Delete(ctx, cacheKey)
}
