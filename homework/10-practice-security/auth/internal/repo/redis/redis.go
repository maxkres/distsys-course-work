package redis

import (
	"fmt"
	"time"

	"auth/config"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDbHost,
		Password: cfg.RedisDbPassword,
		DB:       cfg.RedisDbName,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("client ping failed: %w", err)
	}
	return &Redis{client, cfg}, nil
}

func (r *Redis) SignUp(login, password string) error {
	err := r.client.Set(login, password, time.Duration(0)).Err()
	if err != nil {
		return fmt.Errorf("client set failed: %w", err)
	}
	return nil
}

func (r *Redis) GetPassword(login string) (string, error) {
	val := r.client.Get(login)
	if val.Err() != nil {
		return "", fmt.Errorf("client get failed: %w", val.Err())
	}
	return val.Val(), nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
