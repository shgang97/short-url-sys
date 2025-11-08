package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"redirect-service/internal/client/redis"
	"redirect-service/internal/config"
	"shared/constants"
	shrErrors "shared/errors"
	"time"

	redis9 "github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRepository(client *redis.Client, cfg *config.CacheConfig) *Repository {
	return &Repository{
		client: client,
		prefix: cfg.Prefix,
		ttl:    time.Duration(cfg.TTL) * time.Second,
	}
}

func (r *Repository) getKey(typ string, id string) string {
	return fmt.Sprintf("%s:%s:%s", r.prefix, typ, id)
}

// GetOriginalURL 从 Redis 缓存获取短链映射
func (r *Repository) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	key := r.getKey("url", shortCode)
	result, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis9.Nil) {
			return "", shrErrors.ErrLinkNotFound
		}
		return "", &shrErrors.RepositoryError{Operation: "GetShortURL", Err: err}
	}
	return result, nil
}

// SetShortURL 设置短链映射到缓存
func (r *Repository) SetShortURL(ctx context.Context, shortCode string, longURL string, expiredAt *time.Time) error {
	ttl := r.ttl
	if expiredAt != nil {
		now := time.Now()
		ttl = expiredAt.Sub(now)
		// 已过期
		if ttl <= 0 {
			// 已过期，设置为24小时
			log.Printf("⚠️ expiredAt %v is too close (ttl: %v), using ExpirationToleranceCacheTTL: %v",
				*expiredAt, ttl, constants.ExpirationToleranceCacheTTL)
			ttl = constants.ExpirationToleranceCacheTTL
		} else if ttl > constants.MaxCacheTTL {
			// 超过30天，设置为TTL上线30天
			ttl = constants.MaxCacheTTL
		}
	}
	key := r.getKey("url", shortCode)
	err := r.client.Set(ctx, key, longURL, ttl)
	if err != nil {
		return &shrErrors.RepositoryError{Operation: "SetShortURL", Err: err}
	}
	return nil
}

// DeleteShortURL 删除缓存中的短链接
func (r *Repository) DeleteShortURL(ctx context.Context, shortCode string) error {
	key := r.getKey("url", shortCode)

	err := r.client.Del(ctx, key)
	if err != nil {
		return &shrErrors.RepositoryError{Operation: "DeleteShortURL", Err: err}
	}
	return nil
}
