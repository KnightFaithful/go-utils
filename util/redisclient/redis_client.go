package redisclient

import "github.com/redis/go-redis/v9"

type RedisClient struct {
	RedisClient *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}
