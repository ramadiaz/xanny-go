package helpers

import (
	"xanny-go/pkg/config"
	"time"

	"golang.org/x/net/context"
)

var ctx = context.Background()

func SetBlacklistedToken(token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	return config.RedisClient.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func IsTokenBlacklisted(token string) (bool, error) {
	val, err := config.RedisClient.Exists(ctx, "blacklist:"+token).Result()
	return val == 1, err
}
