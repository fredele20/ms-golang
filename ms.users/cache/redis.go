package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisConnection struct {
	client *redis.Client
	logger *logrus.Logger
}

type RedisStore interface {
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) ([]byte, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

func NewRedisConnection(connectionUri string) (RedisStore, error) {
	fmt.Println(connectionUri)

	client := redis.NewClient(&redis.Options{
		Addr:     connectionUri,
		Password: "",
		DB:       0,
	})

	fmt.Println(client)

	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Redis connection established....")

	return &RedisConnection{
		client: client,
	}, nil
}

func (r *RedisConnection) Get(ctx context.Context, key string) ([]byte, error) {
	fmt.Println(key)

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func (r *RedisConnection) Set(ctx context.Context, key string, value interface{}, duration time.Duration) ([]byte, error) {

	result, err := r.client.Set(ctx, key, value, duration).Result()
	if err != nil {
		r.logger.WithError(err).Error("Something failed")
		return nil, err
	}

	return []byte(result), nil
}

func (r *RedisConnection) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.WithError(err).Error("error deleting key")
		return err
	}

	return nil
}
