package stats

import (
	"context"
	"short-url-sys/internal/model"
	"time"
)

func (s *statsService) GetSystemStats(ctx context.Context) (*model.SystemStats, error) {
	// 获取总链接数
	totalLinks, err := s.statRepo.CountLinks(ctx)
	if err != nil {
		return nil, err
	}
	// 获取总点击量
	totalClicks, err := s.statRepo.CountClicks(ctx)
	if err != nil {
		return nil, err
	}
	// 获取活跃连接数
	totalActiveLinks, err := s.statRepo.CountActiveLinks(ctx)
	if err != nil {
		return nil, err
	}
	// 获取今日点击量
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayClicks, err := s.statRepo.CountClicksSince(ctx, todayStart)
	if err != nil {
		return nil, err
	}
	// 获取本周点击量
	weekStart := getWeekStart()
	weekClicks, err := s.statRepo.CountClicksSince(ctx, weekStart)
	if err != nil {
		return nil, err
	}
	// 获取本月点击量
	monthStart := getMonthStart()
	monthClicks, err := s.statRepo.CountClicksSince(ctx, monthStart)
	if err != nil {
		return nil, err
	}

	stats := model.SystemStats{
		TotalLinks:  totalLinks,
		TotalClicks: totalClicks,
		ActiveLinks: totalActiveLinks,
		TodayClicks: todayClicks,
		WeekClicks:  weekClicks,
		MonthClicks: monthClicks,
	}
	return &stats, nil
}

func (s *statsService) GetTopLinks(ctx context.Context, limit int, days int) ([]model.LinkStats, error) {
	return s.statRepo.GetTopLinks(ctx, limit, days)
}

func (s *statsService) GetClickTimeline(ctx context.Context, shortCode string, hours int) ([]model.ClickTimeline, error) {
	return s.statRepo.GetClickTimeline(ctx, shortCode, hours)
}

func (s *statsService) GetGeographicStats(ctx context.Context, shortCode string) ([]model.GeographicStats, error) {
	return s.statRepo.GetGeographicStats(ctx, shortCode)
}

func (s *statsService) GetPlatformStats(ctx context.Context, shortCode string) ([]model.PlatformStats, error) {
	return s.statRepo.GetPlatformStats(ctx, shortCode)
}

func getWeekStart() time.Time {
	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return now.AddDate(0, 0, -int(weekday)+1).Truncate(24 * time.Hour)
}

func getMonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}
