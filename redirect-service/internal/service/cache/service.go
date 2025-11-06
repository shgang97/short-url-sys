package cache

import (
	"context"
	"redirect-service/internal/repository/cache"
	"time"
)

type Service struct {
	cacheRepo *cache.Repository
}

func NewService(cacheRepo *cache.Repository) *Service {
	return &Service{cacheRepo: cacheRepo}
}

func (s *Service) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	return s.cacheRepo.GetOriginalURL(ctx, shortCode)
}

func (s *Service) SetShortUrl(ctx context.Context, shortCode string, originalUrl string, expiredAt *time.Time) error {
	return s.cacheRepo.SetShortURL(ctx, shortCode, originalUrl, expiredAt)
}

func (s *Service) DelShortUrl(ctx context.Context, shortCode string) error {
	return s.cacheRepo.DeleteShortURL(ctx, shortCode)
}
