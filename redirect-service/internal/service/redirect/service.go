package redirect

import (
	"context"
	"log"
	"redirect-service/internal/client/grpc/generate"
	"redirect-service/internal/repository/cache"
	"time"
)

type Service struct {
	genClient *generate.Client
	cacheRepo *cache.Repository
}

func NewService(client *generate.Client, cacheRepo *cache.Repository) *Service {
	return &Service{
		genClient: client,
		cacheRepo: cacheRepo,
	}
}

func (s *Service) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	// 从缓存获取长链接
	longUrl, err := s.cacheRepo.GetOriginalURL(ctx, shortCode)
	if err == nil {
		return longUrl, nil
	}
	// 缓存未命中，回溯到generate-service服务
	resp, err := s.genClient.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	// 加入缓存
	var expiredAt *time.Time
	if resp.ExpireTime != nil {
		t := resp.ExpireTime.AsTime()
		expiredAt = &t
	}
	go func() {
		if err = s.cacheRepo.SetShortURL(ctx, shortCode, resp.OriginalUrl, expiredAt); err != nil {
			log.Printf("failed to cache short url: %v", err)
		}

	}()
	return resp.OriginalUrl, nil
}
