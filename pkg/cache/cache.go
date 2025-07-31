package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Cache interface defines the contract for all cache implementations
type Cache interface {
	// Get retrieves a value from cache by key
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with optional expiration
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Delete removes a value from cache by key
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// Flush removes all keys from cache
	Flush(ctx context.Context) error

	// InvalidatePattern removes all keys matching a pattern
	InvalidatePattern(ctx context.Context, pattern string) error

	// Close closes the cache connection
	Close() error
}

// CacheOptions holds configuration options for cache implementations
type CacheOptions struct {
	DefaultTTL time.Duration
	MaxSize    int
	Prefix     string
}

// DefaultCacheOptions returns default cache options
func DefaultCacheOptions() *CacheOptions {
	return &CacheOptions{
		DefaultTTL: 5 * time.Minute,
		MaxSize:    1000,
		Prefix:     "cache:",
	}
}

// CacheManager manages multiple cache instances
type CacheManager struct {
	caches map[string]Cache
	opts   *CacheOptions
}

// NewCacheManager creates a new cache manager
func NewCacheManager(opts *CacheOptions) *CacheManager {
	if opts == nil {
		opts = DefaultCacheOptions()
	}

	return &CacheManager{
		caches: make(map[string]Cache),
		opts:   opts,
	}
}

// RegisterCache registers a cache instance with a name
func (cm *CacheManager) RegisterCache(name string, cache Cache) {
	cm.caches[name] = cache
}

// GetCache returns a cache instance by name
func (cm *CacheManager) GetCache(name string) (Cache, error) {
	cache, exists := cm.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache '%s' not found", name)
	}
	return cache, nil
}

// CloseAll closes all registered caches
func (cm *CacheManager) CloseAll() error {
	for name, cache := range cm.caches {
		if err := cache.Close(); err != nil {
			return fmt.Errorf("failed to close cache '%s': %w", name, err)
		}
	}
	return nil
}

// CacheKey generates a cache key with prefix
func (cm *CacheManager) CacheKey(key string) string {
	return cm.opts.Prefix + key
}

// CacheHelper provides convenient methods for common cache operations
type CacheHelper struct {
	cache Cache
	opts  *CacheOptions
}

// NewCacheHelper creates a new cache helper
func NewCacheHelper(cache Cache, opts *CacheOptions) *CacheHelper {
	if opts == nil {
		opts = DefaultCacheOptions()
	}

	return &CacheHelper{
		cache: cache,
		opts:  opts,
	}
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (ch *CacheHelper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := ch.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// SetJSON marshals and stores a JSON value in cache
func (ch *CacheHelper) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if expiration == 0 {
		expiration = ch.opts.DefaultTTL
	}

	return ch.cache.Set(ctx, key, data, expiration)
}

// GetOrSet retrieves a value from cache, or sets it if not found
func (ch *CacheHelper) GetOrSet(ctx context.Context, key string, fn func() ([]byte, time.Duration, error)) ([]byte, error) {
	// Try to get from cache first
	if data, err := ch.cache.Get(ctx, key); err == nil {
		return data, nil
	}

	// If not found, execute the function to get the value
	data, expiration, err := fn()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := ch.cache.Set(ctx, key, data, expiration); err != nil {
		// Log error but don't fail the operation
		// You might want to add proper logging here
	}

	return data, nil
}

// GetOrSetJSON retrieves a JSON value from cache, or sets it if not found
func (ch *CacheHelper) GetOrSetJSON(ctx context.Context, key string, dest interface{}, fn func() (interface{}, time.Duration, error)) error {
	// Try to get from cache first
	if err := ch.GetJSON(ctx, key, dest); err == nil {
		return nil
	}

	// If not found, execute the function to get the value
	value, expiration, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := ch.SetJSON(ctx, key, value, expiration); err != nil {
		// Log error but don't fail the operation
		// You might want to add proper logging here
	}

	// Set the destination
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// InvalidatePattern invalidates all keys matching a pattern (if supported)
func (ch *CacheHelper) InvalidatePattern(ctx context.Context, pattern string) error {
	return ch.cache.InvalidatePattern(ctx, pattern)
}
