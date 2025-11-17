package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"redirect-service/internal/config"
	"redirect-service/internal/metrics"
	"redirect-service/internal/service/circuitbreaker"
	errors2 "shared/errors"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

type Client struct {
	client *redis.Client
	rcb    *circuitbreaker.RedisCircuitBreaker
	// 用于计算缓存命中率的计数器
	hits   uint64
	misses uint64
}

func NewRedis(cfg *config.RedisConfig, rcbCfg *circuitbreaker.RedisCircuitBreakerConfig) (*Client, error) {
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

	// 创建 Redis 熔断器
	rcb := circuitbreaker.NewRedisCircuitBreaker(rcbCfg)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("Connected to redis Successfully\n")
	return &Client{client: client, rcb: rcb}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.client.Ping(ctx).Err()
	// 更新健康状态指标
	if err != nil {
		metrics.RedisHealthStatus.Set(0)
	} else {
		metrics.RedisHealthStatus.Set(1)
	}
	return err
}

type redisResult struct {
	value string
	exist bool
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {

	operation := func() (interface{}, error) {
		val, err := c.client.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 键不存在不算错误，不触发熔断
				// 更新缓存命中次数
				atomic.AddUint64(&c.misses, 1)
				metrics.RedismissesTotal.Inc()
				c.updateHitRatio()
				return redisResult{value: val, exist: false}, nil
			}
			// 其他错误（连接错误、超时等）触发熔断
			log.Printf("Failed to get value from redis, key: %s, err: %v\n", key, err)
			// 更新缓存异常次数
			metrics.RedisErrorsTotal.WithLabelValues("get").Inc()
			return redisResult{value: val, exist: false}, err
		}
		// 键存在，更新缓存命中次数
		atomic.AddUint64(&c.hits, 1)
		metrics.RedisHitsTotal.Inc()
		c.updateHitRatio()
		return redisResult{value: val, exist: true}, err
	}
	result, err := c.rcb.Execute(operation)
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			log.Printf("%s", err.Error())
			// 触发熔断，记录熔断次数
			metrics.RedisErrorsTotal.WithLabelValues("circuit_breaker").Inc()
			return "", errors2.ErrBreakerOpen
		}
		return "", err
	}
	r := result.(redisResult)
	if !r.exist {
		return r.value, redis.Nil
	}
	return r.value, err
}

func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	operation := func() (interface{}, error) {
		err := c.client.Set(ctx, key, value, ttl).Err()
		if err != nil {
			// 连接异常，更新缓存异常次数
			metrics.RedisErrorsTotal.WithLabelValues("set").Inc()
		}
		return nil, err
	}

	_, err := c.rcb.Execute(operation)
	if err != nil {
		log.Printf("%s", err.Error())
		if errors.Is(err, gobreaker.ErrOpenState) {
			// 出发熔断，记录熔断次数
			metrics.RedisErrorsTotal.WithLabelValues("circuit_breaker").Inc()
			return errors2.ErrBreakerOpen
		}
	}
	return err
}

func (c *Client) Del(ctx context.Context, key string) error {
	operation := func() (interface{}, error) {
		err := c.client.Del(ctx, key).Err()
		metrics.RedisErrorsTotal.WithLabelValues("del").Inc()
		return nil, err
	}
	_, err := c.rcb.Execute(operation)
	if err != nil {
		log.Printf("%s", err.Error())
		if errors.Is(err, gobreaker.ErrOpenState) {
			// 触发熔断，记录熔断次数
			metrics.RedisErrorsTotal.WithLabelValues("circuit_breaker").Inc()
			return errors2.ErrBreakerOpen
		}
	}

	return err
}

// 更新缓存命中率
func (c *Client) updateHitRatio() {
	hits := atomic.LoadUint64(&c.hits)
	misses := atomic.LoadUint64(&c.misses)
	total := hits + misses
	if total > 0 {
		hitRatio := float64(hits) / float64(total)
		metrics.CacheHitRatio.Set(hitRatio)
	}
}
