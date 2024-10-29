package redis

import (
	"context"
	"fmt"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr         string
	Password     string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// TODO: build interface to wrap custom methods for mock
func NewRedis(config RedisConfig) (*redisv8.Client, error) {
	client := redisv8.NewClient(&redisv8.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("ping to redis server failed: %w", err)
	}

	return client, nil
}
