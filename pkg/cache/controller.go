package cache

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheController handles cache-related HTTP requests
type CacheController struct {
	cacheService   *CacheService
	userService    *UserService
	productService *ProductService
	statsService   *CacheStatsService
}

// NewCacheController creates a new cache controller
func NewCacheController(cacheService *CacheService) *CacheController {
	return &CacheController{
		cacheService:   cacheService,
		userService:    NewUserService(cacheService),
		productService: NewProductService(cacheService),
		statsService:   NewCacheStatsService(cacheService),
	}
}

// GetUser handles GET /users/:id
func (cc *CacheController) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := cc.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   user,
		"cached": true, // This would be determined by cache hit/miss
	})
}

// GetUsers handles GET /users
func (cc *CacheController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := cc.userService.GetUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdateUser handles PUT /users/:id
func (cc *CacheController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cc.userService.UpdateUser(c.Request.Context(), userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "User updated successfully",
		"cache_invalidated": true,
	})
}

// GetProduct handles GET /products/:id
func (cc *CacheController) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	product, err := cc.productService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

// GetProductPrice handles GET /products/:id/price
func (cc *CacheController) GetProductPrice(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	price, err := cc.productService.GetProductPrice(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"price":      price,
	})
}

// GetCacheStats handles GET /cache/stats
func (cc *CacheController) GetCacheStats(c *gin.Context) {
	stats := cc.statsService.GetStats(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"cache_stats": stats,
	})
}

// FlushCache handles POST /cache/flush
func (cc *CacheController) FlushCache(c *gin.Context) {
	err := cc.statsService.FlushCache(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache flushed successfully",
	})
}

// InvalidatePattern handles POST /cache/invalidate
func (cc *CacheController) InvalidatePattern(c *gin.Context) {
	var request struct {
		Pattern string `json:"pattern" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cc.statsService.InvalidatePattern(c.Request.Context(), request.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache pattern invalidated successfully",
		"pattern": request.Pattern,
	})
}

// SetCacheValue handles POST /cache/set
func (cc *CacheController) SetCacheValue(c *gin.Context) {
	var request struct {
		Key        string      `json:"key" binding:"required"`
		Value      interface{} `json:"value" binding:"required"`
		Expiration string      `json:"expiration"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiration time.Duration
	if request.Expiration != "" {
		var err error
		expiration, err = time.ParseDuration(request.Expiration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiration format"})
			return
		}
	}

	err := cc.cacheService.helper.SetJSON(c.Request.Context(), request.Key, request.Value, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Value cached successfully",
		"key":     request.Key,
	})
}

// GetCacheValue handles GET /cache/:key
func (cc *CacheController) GetCacheValue(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cache key is required"})
		return
	}

	var value interface{}
	err := cc.cacheService.helper.GetJSON(c.Request.Context(), key, &value)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

// DeleteCacheValue handles DELETE /cache/:key
func (cc *CacheController) DeleteCacheValue(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cache key is required"})
		return
	}

	err := cc.cacheService.cache.Delete(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache key deleted successfully",
		"key":     key,
	})
}
