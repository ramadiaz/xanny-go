package cache

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates how to use the cache layer
func Example() {
	// 1. Create cache from environment variables
	cache, err := NewCacheFromEnv()
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	// 2. Create cache service
	cacheService := NewCacheService(cache)

	// 3. Create user service with cache
	userService := NewUserService(cacheService)

	// 4. Example usage
	ctx := context.Background()

	// Get user (will be cached)
	user, err := userService.GetUser(ctx, "123")
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	// Get user again (should be from cache)
	user, err = userService.GetUser(ctx, "123")
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("User (cached): %+v\n", user)
	}

	// Update user (will invalidate cache)
	err = userService.UpdateUser(ctx, "123", map[string]interface{}{
		"name": "Jane Doe",
	})
	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	// Get user again (will fetch from database and cache)
	user, err = userService.GetUser(ctx, "123")
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("User (after update): %+v\n", user)
	}
}

// ExampleWithRedis demonstrates Redis cache usage
func ExampleWithRedis() {
	// Create Redis cache configuration
	config := &CacheConfig{
		Type:       CacheTypeRedis,
		RedisAddr:  "localhost:6379",
		RedisDB:    0,
		DefaultTTL: 10 * time.Minute,
		Prefix:     "app:",
	}

	// Create cache
	cache, err := NewCache(config)
	if err != nil {
		log.Fatalf("Failed to create Redis cache: %v", err)
	}
	defer cache.Close()

	// Create cache service
	cacheService := NewCacheService(cache)

	// Example operations
	ctx := context.Background()

	// Set a value
	err = cacheService.helper.SetJSON(ctx, "test_key", map[string]interface{}{
		"message": "Hello from Redis!",
		"time":    time.Now(),
	}, 5*time.Minute)
	if err != nil {
		log.Printf("Error setting value: %v", err)
	}

	// Get the value
	var result map[string]interface{}
	err = cacheService.helper.GetJSON(ctx, "test_key", &result)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		fmt.Printf("Retrieved: %+v\n", result)
	}

	// Check if key exists
	exists, err := cache.Exists(ctx, "test_key")
	if err != nil {
		log.Printf("Error checking existence: %v", err)
	} else {
		fmt.Printf("Key exists: %v\n", exists)
	}

	// Delete the key
	err = cache.Delete(ctx, "test_key")
	if err != nil {
		log.Printf("Error deleting key: %v", err)
	}
}

// ExampleWithMemory demonstrates in-memory cache usage
func ExampleWithMemory() {
	// Create memory cache configuration
	config := &CacheConfig{
		Type:       CacheTypeMemory,
		DefaultTTL: 5 * time.Minute,
		MaxSize:    1000,
		Prefix:     "mem:",
	}

	// Create cache
	cache, err := NewCache(config)
	if err != nil {
		log.Fatalf("Failed to create memory cache: %v", err)
	}
	defer cache.Close()

	// Example operations
	ctx := context.Background()

	// Set multiple values
	data := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	for key, value := range data {
		err := cache.Set(ctx, key, value, 1*time.Minute)
		if err != nil {
			log.Printf("Error setting %s: %v", key, err)
		}
	}

	// Get multiple values
	keys := []string{"key1", "key2", "key3"}
	if redisCache, ok := cache.(*RedisCache); ok {
		results, err := redisCache.GetMultiple(ctx, keys)
		if err != nil {
			log.Printf("Error getting multiple values: %v", err)
		} else {
			fmt.Printf("Multiple results: %+v\n", results)
		}
	}

	// Get cache statistics
	if memoryCache, ok := cache.(*MemoryCache); ok {
		stats := memoryCache.GetStats()
		fmt.Printf("Memory cache stats: %+v\n", stats)
	}
}

// ExampleWithCacheManager demonstrates cache manager usage
func ExampleWithCacheManager() {
	// Create cache manager
	manager, err := NewCacheManagerFromEnv()
	if err != nil {
		log.Fatalf("Failed to create cache manager: %v", err)
	}
	defer manager.CloseAll()

	// Get primary cache
	primaryCache, err := manager.GetCache("primary")
	if err != nil {
		log.Printf("Error getting primary cache: %v", err)
		return
	}

	// Get fallback cache
	fallbackCache, err := manager.GetCache("fallback")
	if err != nil {
		log.Printf("Error getting fallback cache: %v", err)
		return
	}

	// Example with fallback strategy
	ctx := context.Background()
	key := "important_data"

	// Try primary cache first
	data, err := primaryCache.Get(ctx, key)
	if err != nil {
		fmt.Println("Primary cache miss, trying fallback...")

		// Try fallback cache
		data, err = fallbackCache.Get(ctx, key)
		if err != nil {
			fmt.Println("Fallback cache miss, fetching from source...")

			// Fetch from source and store in both caches
			data = []byte("fetched from database")

			// Store in primary cache
			primaryCache.Set(ctx, key, data, 10*time.Minute)

			// Store in fallback cache
			fallbackCache.Set(ctx, key, data, 5*time.Minute)
		}
	}

	fmt.Printf("Retrieved data: %s\n", string(data))
}

// ExampleWithHealthCheck demonstrates cache health checking
func ExampleWithHealthCheck() {
	// Create cache
	cache, err := NewCacheFromEnv()
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	// Create health checker
	healthChecker := NewCacheHealthChecker(cache)

	// Check health
	ctx := context.Background()
	err = healthChecker.CheckHealth(ctx)
	if err != nil {
		log.Printf("Cache health check failed: %v", err)
	} else {
		fmt.Println("Cache is healthy")
	}

	// Get cache statistics
	stats := healthChecker.GetCacheStats(ctx)
	fmt.Printf("Cache stats: %+v\n", stats)
}
