package stats

import (
	"context"
	"short-url-sys/internal/model"
	"time"
)

type Repository interface {
	// RecordClick 记录点击
	RecordClick(ctx context.Context, stats *model.ClickStats) error

	// GetStatsSummary 查询统计
	GetStatsSummary(ctx context.Context, shortCode string, startDate, endDate *time.Time) (*model.StatsSummary, error)
	GetDailyStats(ctx context.Context, shortCode string, days int) ([]model.DailyStats, error)
	GetLastAccessed(ctx context.Context, shortCode string) (*time.Time, error)

	// 系统统计信息
	CountLinks(ctx context.Context) (int64, error)
	CountClicks(ctx context.Context) (int64, error)
	CountActiveLinks(ctx context.Context) (int64, error)
	CountClicksSince(ctx context.Context, since time.Time) (int64, error)

	// 高级统计信息
	GetTopLinks(ctx context.Context, limit int, days int) ([]model.LinkStats, error)
	GetClickTimeline(ctx context.Context, shortCode string, hours int) ([]model.ClickTimeline, error)
	GetGeographicStats(ctx context.Context, shortCode string) ([]model.GeographicStats, error)
	GetPlatformStats(ctx context.Context, shortCode string) ([]model.PlatformStats, error)
}
