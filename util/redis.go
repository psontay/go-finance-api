package util

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisAddress string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("cannot connect to Redis: " + err.Error())
	}
	return client
}
