package cache

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConnection struct {
	client *redis.Client
}

func NewRedisConnection() *RedisConnection {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	return &RedisConnection{
		client: client,
	}
}

func (r *RedisConnection) Get(ctx context.Context, key string) ([]byte, error) {

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func (r *RedisConnection) Set(ctx context.Context, key string, value []byte, duration time.Duration) ([]byte, error) {

	result, err := r.client.Set(ctx, key, bytes.NewBuffer(value).Bytes(), duration).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}
