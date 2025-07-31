package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheType represents the type of cache to use
type CacheType string

const (
	CacheTypeRedis  CacheType = "redis"
	CacheTypeMemory CacheType = "memory"
)

// CacheConfig holds configuration for cache initialization
type CacheConfig struct {
	Type          CacheType
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DefaultTTL    time.Duration
	MaxSize       int
	Prefix        string
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Type:       CacheTypeMemory,
		DefaultTTL: 5 * time.Minute,
		MaxSize:    1000,
		Prefix:     "cache:",
	}
}

// NewCache creates a new cache instance based on configuration
func NewCache(config *CacheConfig) (Cache, error) {
	if config == nil {
		config = DefaultCacheConfig()
	}

	opts := &CacheOptions{
		DefaultTTL: config.DefaultTTL,
		MaxSize:    config.MaxSize,
		Prefix:     config.Prefix,
	}

	switch config.Type {
	case CacheTypeRedis:
		return NewRedisCacheFromConfig(config, opts)
	case CacheTypeMemory:
		return NewMemoryCache(opts), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", config.Type)
	}
}

// NewRedisCacheFromConfig creates a Redis cache from configuration
func NewRedisCacheFromConfig(config *CacheConfig, opts *CacheOptions) (Cache, error) {
	if config.RedisAddr == "" {
		return nil, fmt.Errorf("Redis address is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return NewRedisCache(client, opts), nil
}

// NewCacheFromEnv creates a cache instance from environment variables
func NewCacheFromEnv() (Cache, error) {
	config := DefaultCacheConfig()

	// Read cache type from environment
	if cacheType := os.Getenv("CACHE_TYPE"); cacheType != "" {
		config.Type = CacheType(cacheType)
	}

	// Read Redis configuration from environment
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		config.RedisAddr = redisAddr
	}

	if redisPass := os.Getenv("REDIS_PASS"); redisPass != "" {
		config.RedisPassword = redisPass
	}

	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		if db, err := strconv.Atoi(redisDB); err == nil {
			config.RedisDB = db
		}
	}

	// Read cache options from environment
	if ttl := os.Getenv("CACHE_DEFAULT_TTL"); ttl != "" {
		if duration, err := time.ParseDuration(ttl); err == nil {
			config.DefaultTTL = duration
		}
	}

	if maxSize := os.Getenv("CACHE_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.Atoi(maxSize); err == nil {
			config.MaxSize = size
		}
	}

	if prefix := os.Getenv("CACHE_PREFIX"); prefix != "" {
		config.Prefix = prefix
	}

	return NewCache(config)
}

// NewCacheManagerFromEnv creates a cache manager with caches from environment
func NewCacheManagerFromEnv() (*CacheManager, error) {
	manager := NewCacheManager(nil)

	// Create primary cache
	primaryCache, err := NewCacheFromEnv()
	if err != nil {
		return nil, err
	}

	manager.RegisterCache("primary", primaryCache)

	// Create memory cache as fallback
	memoryCache := NewMemoryCache(&CacheOptions{
		DefaultTTL: 1 * time.Minute,
		MaxSize:    100,
		Prefix:     "fallback:",
	})

	manager.RegisterCache("fallback", memoryCache)

	return manager, nil
}

// CacheHealthChecker checks the health of cache instances
type CacheHealthChecker struct {
	cache Cache
}

// NewCacheHealthChecker creates a new cache health checker
func NewCacheHealthChecker(cache Cache) *CacheHealthChecker {
	return &CacheHealthChecker{cache: cache}
}

// CheckHealth checks if the cache is healthy
func (chc *CacheHealthChecker) CheckHealth(ctx context.Context) error {
	// Try to set and get a test value
	testKey := "health_check"
	testValue := []byte("ok")

	if err := chc.cache.Set(ctx, testKey, testValue, time.Second); err != nil {
		return fmt.Errorf("failed to set test value: %w", err)
	}

	if _, err := chc.cache.Get(ctx, testKey); err != nil {
		return fmt.Errorf("failed to get test value: %w", err)
	}

	if err := chc.cache.Delete(ctx, testKey); err != nil {
		return fmt.Errorf("failed to delete test value: %w", err)
	}

	return nil
}

// GetCacheStats returns statistics about the cache
func (chc *CacheHealthChecker) GetCacheStats(ctx context.Context) map[string]interface{} {
	stats := map[string]interface{}{
		"type": "unknown",
	}

	// Try to get specific stats based on cache type
	if _, ok := chc.cache.(*RedisCache); ok {
		stats["type"] = "redis"
		// Add Redis-specific stats if needed
	} else if memoryCache, ok := chc.cache.(*MemoryCache); ok {
		stats["type"] = "memory"
		stats["stats"] = memoryCache.GetStats()
	}

	return stats
}
