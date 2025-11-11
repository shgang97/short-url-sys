package click

import (
	"context"
	"statistics-service/internal/model"
	"time"
)

func (r *repository) GetClickTimeline(
	ctx context.Context,
	shortCode string,
	startTime *time.Time,
	endTime *time.Time,
	groupExpr string,
	periodExpr string,
) ([]model.TimeSeriesStats, error) {
	var data []model.TimeSeriesStats

	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select(periodExpr+"as period, COUNT(*) as clicks, COUNT(DISTINCT ip) as unique_visitors ").
		Where("short_code = ? AND Delete_flag = 'N'", shortCode)
	if startTime != nil {
		query = query.Where("click_time >= ?", startTime.Format(time.DateOnly))
	}
	if endTime != nil {
		query = query.Where("click_time <= ?", endTime.Format(time.DateOnly))
	}
	err := query.
		Group(groupExpr).
		Order("period desc").
		Find(&data).Error
	return data, err
}

// GetGeographicStats 获取地理信息统计
func (r *repository) GetGeographicStats(ctx context.Context, shortCode string) ([]*model.GeographicStats, error) {
	var stats []*model.GeographicStats
	err := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("country, region, city, COUNT(*) as clicks, COUNT(DISTINCT ip) as unique_visitors ").
		Where("short_code = ? AND Delete_flag = 'N' AND country IS NOT NULL", shortCode).
		Group("country, region, city").
		Order("clicks desc").Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// GetPlatformStats 获取设备统计
func (r *repository) GetPlatformStats(ctx context.Context, shortCode string) ([]*model.PlatformStats, error) {
	var stats []*model.PlatformStats
	err := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("device_type, COUNT(*) as clicks, COUNT(DISTINCT ip) as unique_visitors ").
		Where("short_code = ? AND Delete_flag = 'N' AND device_type IS NOT NULL", shortCode).
		Group("device_type, city").
		Order("clicks desc").Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}
