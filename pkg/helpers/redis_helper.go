package helpers

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
}

func SetBlacklistedToken(token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	return redisClient.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func IsTokenBlacklisted(token string) (bool, error) {
	val, err := redisClient.Exists(ctx, "blacklist:"+token).Result()
	return val == 1, err
}
