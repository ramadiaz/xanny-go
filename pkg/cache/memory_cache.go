package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// cacheItem represents an item stored in memory cache
type cacheItem struct {
	value      []byte
	expiration time.Time
}

// MemoryCache implements the Cache interface using in-memory storage
type MemoryCache struct {
	data sync.Map
	opts *CacheOptions
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache(opts *CacheOptions) *MemoryCache {
	if opts == nil {
		opts = DefaultCacheOptions()
	}

	cache := &MemoryCache{
		opts: opts,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a value from memory cache
func (mc *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	key = mc.opts.Prefix + key

	value, exists := mc.data.Load(key)
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	item, ok := value.(*cacheItem)
	if !ok {
		return nil, fmt.Errorf("invalid cache item type for key: %s", key)
	}

	// Check if item has expired
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		mc.data.Delete(key)
		return nil, fmt.Errorf("key expired: %s", key)
	}

	return item.value, nil
}

// Set stores a value in memory cache with expiration
func (mc *MemoryCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	key = mc.opts.Prefix + key

	var expirationTime time.Time
	if expiration > 0 {
		expirationTime = time.Now().Add(expiration)
	} else if mc.opts.DefaultTTL > 0 {
		expirationTime = time.Now().Add(mc.opts.DefaultTTL)
	}

	item := &cacheItem{
		value:      value,
		expiration: expirationTime,
	}

	mc.data.Store(key, item)
	return nil
}

// Delete removes a value from memory cache
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	key = mc.opts.Prefix + key
	mc.data.Delete(key)
	return nil
}

// Exists checks if a key exists in memory cache
func (mc *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	key = mc.opts.Prefix + key

	value, exists := mc.data.Load(key)
	if !exists {
		return false, nil
	}

	item, ok := value.(*cacheItem)
	if !ok {
		return false, nil
	}

	// Check if item has expired
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		mc.data.Delete(key)
		return false, nil
	}

	return true, nil
}

// Flush removes all keys from memory cache
func (mc *MemoryCache) Flush(ctx context.Context) error {
	mc.data = sync.Map{}
	return nil
}

// Close closes the memory cache (no-op for memory cache)
func (mc *MemoryCache) Close() error {
	return nil
}

// cleanup periodically removes expired items
func (mc *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		mc.data.Range(func(key, value interface{}) bool {
			item, ok := value.(*cacheItem)
			if !ok {
				mc.data.Delete(key)
				return true
			}

			if !item.expiration.IsZero() && now.After(item.expiration) {
				mc.data.Delete(key)
			}

			return true
		})
	}
}

// GetStats returns cache statistics
func (mc *MemoryCache) GetStats() map[string]interface{} {
	count := 0
	mc.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return map[string]interface{}{
		"size":     count,
		"max_size": mc.opts.MaxSize,
	}
}

// GetKeys returns all keys in the cache
func (mc *MemoryCache) GetKeys() []string {
	var keys []string
	mc.data.Range(func(key, value interface{}) bool {
		if keyStr, ok := key.(string); ok {
			// Remove prefix for consistency
			if len(keyStr) > len(mc.opts.Prefix) {
				keys = append(keys, keyStr[len(mc.opts.Prefix):])
			}
		}
		return true
	})
	return keys
}

// InvalidatePattern removes all keys matching a pattern
func (mc *MemoryCache) InvalidatePattern(ctx context.Context, pattern string) error {
	// Simple pattern matching for memory cache
	// This is a basic implementation - you might want to use regex for more complex patterns
	pattern = mc.opts.Prefix + pattern

	var keysToDelete []interface{}
	mc.data.Range(func(key, value interface{}) bool {
		if keyStr, ok := key.(string); ok {
			// Simple contains check - you can implement more sophisticated pattern matching
			if contains(keyStr, pattern) {
				keysToDelete = append(keysToDelete, key)
			}
		}
		return true
	})

	for _, key := range keysToDelete {
		mc.data.Delete(key)
	}

	return nil
}

// contains checks if a string contains another string (simple pattern matching)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
