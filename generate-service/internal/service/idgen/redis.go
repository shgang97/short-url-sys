package idgen

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisGenerator struct {
	client *redis.Client
	key    string
}

func (r *RedisGenerator) NextId() (uint64, error) {
	ctx := context.Background()
	id, err := r.client.Incr(ctx, r.key).Result()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *RedisGenerator) String() string {
	return "redis"
}

func NewRedisGenerator(client *redis.Client, key string) *RedisGenerator {
	return &RedisGenerator{
		client: client,
		key:    key,
	}
}
