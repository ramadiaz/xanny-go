package config

import (
	"context"
	"xanny-go/pkg/exceptions"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() *exceptions.Exception {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     GetRedisAddr(),
		Password: GetRedisPass(),
		DB:       0,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return exceptions.NewException(500, "Failed to connect to Redis: "+err.Error())
	}

	return nil
}
