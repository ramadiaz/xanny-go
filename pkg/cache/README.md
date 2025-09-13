# Cache Layer for Xanny Go

A comprehensive, reusable cache layer for Go applications with support for multiple cache backends, HTTP response caching, and advanced cache management features.

## Features

- **Multiple Cache Backends**: Redis and in-memory cache support
- **Unified Interface**: Common interface for all cache implementations
- **HTTP Response Caching**: Middleware for caching HTTP responses
- **Cache Management**: Cache manager for multiple cache instances
- **Health Checking**: Built-in cache health monitoring
- **Pattern Invalidation**: Support for pattern-based cache invalidation
- **JSON Support**: Built-in JSON serialization/deserialization
- **TTL Management**: Configurable time-to-live for cache entries
- **Fallback Strategy**: Support for primary/fallback cache strategies

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "xanny-go/pkg/cache"
)

func main() {
    // Create cache from environment variables
    cache, err := cache.NewCacheFromEnv()
    if err != nil {
        log.Fatalf("Failed to create cache: %v", err)
    }
    defer cache.Close()
    
    // Create cache service
    cacheService := cache.NewCacheService(cache)
    
    ctx := context.Background()
    
    // Set a value
    err = cacheService.helper.SetJSON(ctx, "user:123", map[string]interface{}{
        "id":   "123",
        "name": "John Doe",
    }, 10*time.Minute)
    
    // Get a value
    var user map[string]interface{}
    err = cacheService.helper.GetJSON(ctx, "user:123", &user)
}
```

### Environment Configuration

Set these environment variables to configure the cache:

```bash
# Cache type (redis or memory)
CACHE_TYPE=redis

# Redis configuration
REDIS_ADDR=localhost:6379
REDIS_PASS=your_password
REDIS_DB=0

# Cache options
CACHE_DEFAULT_TTL=5m
CACHE_MAX_SIZE=1000
CACHE_PREFIX=cache:
```

### HTTP Response Caching

```go
package main

import (
    "github.com/gin-gonic/gin"
    "xanny-go/pkg/cache"
)

func main() {
    r := gin.New()
    
    // Create cache
    cache, _ := cache.NewCacheFromEnv()
    cacheService := cache.NewCacheService(cache)
    
    // Setup cache middleware
    cache.SetupCacheMiddleware(r, cacheService)
    
    // Setup cache routes
    api := r.Group("/api")
    cache.SetupCacheRoutes(api, cacheService)
    
    r.Run(":8080")
}
```

## Architecture

### Core Components

1. **Cache Interface**: Defines the contract for all cache implementations
2. **Redis Cache**: Redis-backed cache implementation
3. **Memory Cache**: In-memory cache using sync.Map
4. **Cache Manager**: Manages multiple cache instances
5. **Cache Helper**: Provides convenient JSON operations
6. **Cache Middleware**: HTTP response caching for Gin
7. **Cache Service**: High-level cache operations for business logic

### Cache Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Flush(ctx context.Context) error
    InvalidatePattern(ctx context.Context, pattern string) error
    Close() error
}
```

## Usage Examples

### Service Layer Integration

```go
type UserService struct {
    cacheService *cache.CacheService
    userRepo     UserRepository
}

func (us *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    var user User
    
    cacheKey := fmt.Sprintf("user:%s", userID)
    err := us.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &user, func() (interface{}, time.Duration, error) {
        // This function is called when cache miss occurs
        user, err := us.userRepo.FindByID(ctx, userID)
        if err != nil {
            return nil, 0, err
        }
        
        return user, 10*time.Minute, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
    // Update database
    err := us.userRepo.Update(ctx, userID, updates)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%s", userID)
    us.cacheService.cache.Delete(ctx, cacheKey)
    
    // Invalidate related caches
    us.cacheService.cache.InvalidatePattern(ctx, "users:*")
    
    return nil
}
```

### Cache Manager with Fallback

```go
// Create cache manager with primary and fallback caches
manager, err := cache.NewCacheManagerFromEnv()
if err != nil {
    log.Fatal(err)
}
defer manager.CloseAll()

// Get caches
primaryCache, _ := manager.GetCache("primary")
fallbackCache, _ := manager.GetCache("fallback")

// Try primary cache first, then fallback
data, err := primaryCache.Get(ctx, key)
if err != nil {
    data, err = fallbackCache.Get(ctx, key)
    if err != nil {
        // Fetch from source
        data = fetchFromSource()
        
        // Store in both caches
        primaryCache.Set(ctx, key, data, 10*time.Minute)
        fallbackCache.Set(ctx, key, data, 5*time.Minute)
    }
}
```

### HTTP Response Caching

```go
// Create cache middleware
cacheMiddleware := cache.CacheMiddleware(&cache.CacheMiddlewareOptions{
    Cache:      cacheInstance,
    DefaultTTL: 5 * time.Minute,
    SkipCache: func(c *gin.Context) bool {
        // Skip caching for authenticated requests
        return c.GetHeader("Authorization") != ""
    },
})

// Apply to routes
publicGroup := router.Group("/api/public")
publicGroup.Use(cacheMiddleware)
```

## API Endpoints

The cache layer provides HTTP endpoints for cache management:

- `GET /api/cache/stats` - Get cache statistics
- `POST /api/cache/flush` - Flush all cache
- `POST /api/cache/invalidate` - Invalidate cache by pattern
- `POST /api/cache/set` - Set cache value
- `GET /api/cache/:key` - Get cache value
- `DELETE /api/cache/:key` - Delete cache value

## Configuration Options

### Cache Options

```go
type CacheOptions struct {
    DefaultTTL time.Duration // Default time-to-live for cache entries
    MaxSize    int           // Maximum number of entries (for memory cache)
    Prefix     string        // Key prefix for all cache entries
}
```

### Cache Config

```go
type CacheConfig struct {
    Type           CacheType     // Cache type (redis or memory)
    RedisAddr      string        // Redis server address
    RedisPassword  string        // Redis password
    RedisDB        int           // Redis database number
    DefaultTTL     time.Duration // Default TTL
    MaxSize        int           // Maximum size
    Prefix         string        // Key prefix
}
```

## Best Practices

### 1. Cache Key Naming

Use consistent, hierarchical key naming:

```go
// Good
"user:123"
"user:123:profile"
"users:page:1:limit:10"
"product:456:price"

// Avoid
"user123"
"user_profile_123"
"page1_users"
```

### 2. TTL Strategy

- **Short TTL (1-5 minutes)**: Frequently changing data (prices, stock)
- **Medium TTL (10-30 minutes)**: Moderately changing data (user profiles)
- **Long TTL (1+ hours)**: Rarely changing data (product details, configurations)

### 3. Cache Invalidation

```go
// Invalidate specific keys
cache.Delete(ctx, "user:123")

// Invalidate patterns
cache.InvalidatePattern(ctx, "users:*")
cache.InvalidatePattern(ctx, "product:*:price")

// Invalidate after updates
func (s *Service) UpdateUser(ctx context.Context, userID string, data map[string]interface{}) error {
    // Update database
    err := s.repo.Update(ctx, userID, data)
    if err != nil {
        return err
    }
    
    // Invalidate related caches
    cache.Delete(ctx, fmt.Sprintf("user:%s", userID))
    cache.InvalidatePattern(ctx, "users:*")
    
    return nil
}
```

### 4. Error Handling

```go
// Cache operations should not fail the main operation
func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
    var user User
    
    err := s.cacheService.helper.GetOrSetJSON(ctx, cacheKey, &user, func() (interface{}, time.Duration, error) {
        // Fetch from database
        return s.repo.FindByID(ctx, userID), 10*time.Minute, nil
    })
    
    if err != nil {
        // Log error but continue with database fetch
        log.Printf("Cache error: %v", err)
        
        // Fetch directly from database
        user, err = s.repo.FindByID(ctx, userID)
        if err != nil {
            return nil, err
        }
    }
    
    return &user, nil
}
```

## Monitoring and Health Checks

```go
// Create health checker
healthChecker := cache.NewCacheHealthChecker(cacheInstance)

// Check health
err := healthChecker.CheckHealth(ctx)
if err != nil {
    log.Printf("Cache health check failed: %v", err)
}

// Get statistics
stats := healthChecker.GetCacheStats(ctx)
fmt.Printf("Cache stats: %+v\n", stats)
```

## Performance Considerations

1. **Key Size**: Keep cache keys short and meaningful
2. **Value Size**: Avoid caching large objects unnecessarily
3. **TTL**: Set appropriate TTL based on data volatility
4. **Pattern Matching**: Use specific patterns for invalidation
5. **Connection Pooling**: Redis connection pooling for high throughput

## Troubleshooting

### Common Issues

1. **Cache Miss**: Check if keys are being set correctly
2. **Memory Usage**: Monitor memory cache size and TTL
3. **Redis Connection**: Verify Redis connectivity and configuration
4. **Pattern Invalidation**: Ensure patterns match your key naming convention

### Debugging

```go
// Enable debug logging
log.SetLevel(log.DebugLevel)

// Check cache existence
exists, err := cache.Exists(ctx, "test_key")
fmt.Printf("Key exists: %v\n", exists)

// Get cache statistics
if memoryCache, ok := cache.(*cache.MemoryCache); ok {
    stats := memoryCache.GetStats()
    fmt.Printf("Memory cache stats: %+v\n", stats)
}
```

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Ensure backward compatibility

## License

This cache layer is part of the Xanny Go project and follows the same license terms. 