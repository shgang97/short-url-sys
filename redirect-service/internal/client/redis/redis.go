package redis

import (
	"context"
	"fmt"
	"log"
	"redirect-service/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Client *redis.Client
}

func NewRedis(cfg *config.RedisConfig) (*Client, error) {
	options := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	client := redis.NewClient(options)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("Connected to redis Successfully\n")
	return &Client{Client: client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Client.Ping(ctx).Err()
}
