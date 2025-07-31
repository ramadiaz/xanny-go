package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	opts   *CacheOptions
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client, opts *CacheOptions) *RedisCache {
	if opts == nil {
		opts = DefaultCacheOptions()
	}

	return &RedisCache{
		client: client,
		opts:   opts,
	}
}

// Get retrieves a value from Redis cache
func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	key = rc.opts.Prefix + key
	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, err
	}
	return []byte(result), nil
}

// Set stores a value in Redis cache with expiration
func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	key = rc.opts.Prefix + key
	if expiration == 0 {
		expiration = rc.opts.DefaultTTL
	}
	return rc.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from Redis cache
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	key = rc.opts.Prefix + key
	return rc.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis cache
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	key = rc.opts.Prefix + key
	result, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Flush removes all keys from Redis cache (use with caution!)
func (rc *RedisCache) Flush(ctx context.Context) error {
	pattern := rc.opts.Prefix + "*"
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...).Err()
	}
	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// InvalidatePattern removes all keys matching a pattern
func (rc *RedisCache) InvalidatePattern(ctx context.Context, pattern string) error {
	pattern = rc.opts.Prefix + pattern
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...).Err()
	}
	return nil
}

// GetTTL returns the remaining TTL for a key
func (rc *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	key = rc.opts.Prefix + key
	return rc.client.TTL(ctx, key).Result()
}

// SetNX sets a value only if the key doesn't exist
func (rc *RedisCache) SetNX(ctx context.Context, key string, value []byte, expiration time.Duration) (bool, error) {
	key = rc.opts.Prefix + key
	if expiration == 0 {
		expiration = rc.opts.DefaultTTL
	}
	return rc.client.SetNX(ctx, key, value, expiration).Result()
}

// Increment increments a numeric value
func (rc *RedisCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	key = rc.opts.Prefix + key
	return rc.client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements a numeric value
func (rc *RedisCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	key = rc.opts.Prefix + key
	return rc.client.DecrBy(ctx, key, value).Result()
}

// GetMultiple retrieves multiple values at once
func (rc *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string][]byte, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = rc.opts.Prefix + key
	}

	results, err := rc.client.MGet(ctx, prefixedKeys...).Result()
	if err != nil {
		return nil, err
	}

	data := make(map[string][]byte)
	for i, result := range results {
		if result != nil {
			if str, ok := result.(string); ok {
				data[keys[i]] = []byte(str)
			}
		}
	}

	return data, nil
}

// SetMultiple sets multiple values at once
func (rc *RedisCache) SetMultiple(ctx context.Context, data map[string][]byte, expiration time.Duration) error {
	if expiration == 0 {
		expiration = rc.opts.DefaultTTL
	}

	pipe := rc.client.Pipeline()
	for key, value := range data {
		prefixedKey := rc.opts.Prefix + key
		pipe.Set(ctx, prefixedKey, value, expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}
