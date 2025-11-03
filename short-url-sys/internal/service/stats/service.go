package stats

import (
	"context"
	"short-url-sys/internal/model"
	statsRepo "short-url-sys/internal/repository/stats"
	"time"
)

// Service 数据统计接口
type Service interface {
	GetLinkStats(ctx context.Context, shortCode string, startDate, endDate *time.Time) (*model.StatsResponse, error)
	GetDailyStats(ctx context.Context, shortCode string, days int) ([]model.DailyStats, error)

	// GetSystemStats 系统统计信息
	GetSystemStats(ctx context.Context) (*model.SystemStats, error)

	// 高级统计信息
	GetTopLinks(ctx context.Context, limit int, days int) ([]model.LinkStats, error)
	GetClickTimeline(ctx context.Context, shortCode string, hours int) ([]model.ClickTimeline, error)
	GetGeographicStats(ctx context.Context, shortCode string) ([]model.GeographicStats, error)
	GetPlatformStats(ctx context.Context, shortCode string) ([]model.PlatformStats, error)
}

type statsService struct {
	statRepo statsRepo.Repository
}

func (s *statsService) GetLinkStats(ctx context.Context, shortCode string, startDate, endDate *time.Time) (*model.StatsResponse, error) {
	summary, err := s.statRepo.GetStatsSummary(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, err
	}
	resp := &model.StatsResponse{
		ShortCode:   shortCode,
		TotalClicks: summary.TotalClicks,
		DailyStats:  summary.DailyStats,
		Referrers:   summary.Referrers,
		Countries:   summary.Countries,
		Devices:     summary.Devices,
	}
	return resp, nil
}

func (s *statsService) GetDailyStats(ctx context.Context, shortCode string, days int) ([]model.DailyStats, error) {
	return s.statRepo.GetDailyStats(ctx, shortCode, days)
}

func NewStatsService(statRepo statsRepo.Repository) Service {
	return &statsService{
		statRepo: statRepo,
	}
}
