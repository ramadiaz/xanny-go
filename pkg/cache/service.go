package cache

import (
	"context"
	"fmt"
	"time"
)

// CacheService provides high-level cache operations for business logic
type CacheService struct {
	cache  Cache
	helper *CacheHelper
}

// NewCacheService creates a new cache service
func NewCacheService(cache Cache) *CacheService {
	return &CacheService{
		cache:  cache,
		helper: NewCacheHelper(cache, nil),
	}
}

// User represents a user entity
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// UserService demonstrates how to use cache in a service layer
type UserService struct {
	cacheService *CacheService
	// In a real application, you would have a repository here
	// userRepo UserRepository
}

// NewUserService creates a new user service with cache
func NewUserService(cacheService *CacheService) *UserService {
	return &UserService{
		cacheService: cacheService,
	}
}

// GetUser retrieves a user with caching
func (us *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User

	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%s", userID)
	err := us.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &user, func() (interface{}, time.Duration, error) {
		// This function is called when cache miss occurs
		// In a real application, you would fetch from database here
		// user, err := us.userRepo.FindByID(ctx, userID)

		// Simulate database fetch
		user = User{
			ID:       userID,
			Name:     "John Doe",
			Email:    "john@example.com",
			Created:  time.Now().Add(-24 * time.Hour),
			Modified: time.Now(),
		}

		// Return the user, cache TTL, and any error
		return user, 10 * time.Minute, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates a user and invalidates cache
func (us *UserService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	// In a real application, you would update the database here
	// err := us.userRepo.Update(ctx, userID, updates)

	// Invalidate user cache
	cacheKey := fmt.Sprintf("user:%s", userID)
	if err := us.cacheService.cache.Delete(ctx, cacheKey); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to invalidate cache for user %s: %v\n", userID, err)
	}

	// Invalidate user list cache
	if err := us.cacheService.cache.InvalidatePattern(ctx, "users:*"); err != nil {
		fmt.Printf("Failed to invalidate user list cache: %v\n", err)
	}

	return nil
}

// GetUsers retrieves a list of users with caching
func (us *UserService) GetUsers(ctx context.Context, page, limit int) ([]User, error) {
	var users []User

	cacheKey := fmt.Sprintf("users:page:%d:limit:%d", page, limit)
	err := us.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &users, func() (interface{}, time.Duration, error) {
		// Simulate database fetch
		users = []User{
			{ID: "1", Name: "John Doe", Email: "john@example.com", Created: time.Now().Add(-24 * time.Hour), Modified: time.Now()},
			{ID: "2", Name: "Jane Smith", Email: "jane@example.com", Created: time.Now().Add(-12 * time.Hour), Modified: time.Now()},
		}

		return users, 5 * time.Minute, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// Product represents a product entity
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

// ProductService demonstrates cache with different TTL strategies
type ProductService struct {
	cacheService *CacheService
}

// NewProductService creates a new product service with cache
func NewProductService(cacheService *CacheService) *ProductService {
	return &ProductService{
		cacheService: cacheService,
	}
}

// GetProduct retrieves a product with long-term caching
func (ps *ProductService) GetProduct(ctx context.Context, productID string) (*Product, error) {
	var product Product

	cacheKey := fmt.Sprintf("product:%s", productID)
	err := ps.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &product, func() (interface{}, time.Duration, error) {
		// Simulate database fetch
		product = Product{
			ID:          productID,
			Name:        "Sample Product",
			Price:       99.99,
			Description: "A sample product description",
			Category:    "Electronics",
		}

		// Products can be cached longer since they don't change frequently
		return product, 1 * time.Hour, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// GetProductPrice retrieves just the price with short-term caching
func (ps *ProductService) GetProductPrice(ctx context.Context, productID string) (float64, error) {
	var price float64

	cacheKey := fmt.Sprintf("product_price:%s", productID)
	err := ps.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &price, func() (interface{}, time.Duration, error) {
		// Simulate database fetch
		price = 99.99

		// Prices can change frequently, so cache for shorter time
		return price, 5 * time.Minute, nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}

	return price, nil
}

// CacheStatsService provides cache statistics and management
type CacheStatsService struct {
	cacheService *CacheService
}

// NewCacheStatsService creates a new cache stats service
func NewCacheStatsService(cacheService *CacheService) *CacheStatsService {
	return &CacheStatsService{
		cacheService: cacheService,
	}
}

// GetStats returns cache statistics
func (css *CacheStatsService) GetStats(ctx context.Context) map[string]interface{} {
	stats := map[string]interface{}{
		"cache_type": "unknown",
	}

	// Try to get specific stats based on cache type
	if _, ok := css.cacheService.cache.(*RedisCache); ok {
		stats["cache_type"] = "redis"
		// Add Redis-specific stats
	} else if memoryCache, ok := css.cacheService.cache.(*MemoryCache); ok {
		stats["cache_type"] = "memory"
		stats["memory_stats"] = memoryCache.GetStats()
	}

	return stats
}

// FlushCache clears all cached data
func (css *CacheStatsService) FlushCache(ctx context.Context) error {
	return css.cacheService.cache.Flush(ctx)
}

// InvalidatePattern removes all keys matching a pattern
func (css *CacheStatsService) InvalidatePattern(ctx context.Context, pattern string) error {
	return css.cacheService.cache.InvalidatePattern(ctx, pattern)
}
