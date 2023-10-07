package config

import (
	"github.com/go-redis/redis"
	"os"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
	)

	return &RedisClient{client}
}

func (r *RedisClient) Get(key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RedisClient) Set(key string, value string) error {
	err := r.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
