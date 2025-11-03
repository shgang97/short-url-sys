package stats

import (
	"context"
	"short-url-sys/internal/model"
	"time"
)

// GetTopLinks 获取热门链接
func (r *MySQLRepository) GetTopLinks(ctx context.Context, limit int, days int) ([]model.LinkStats, error) {
	var stats []model.LinkStats
	startDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT
			l.short_code,
			l.long_url,
			l.created_at,
			l.click_count,
			COUNT(cs.id) as recent_clicks
		FROM links l
		LEFT JOIN click_stats cs ON l.short_code = cs.short_code AND cs.created_at >= ?
		WHERE l.status = 'active'
		GROUP BY l.short_code, l.long_url
		ORDER BY recent_clicks DESC, l.click_count DESC
		LIMIT ?
	`
	result := r.db.WithContext(ctx).Raw(query, startDate, limit).Scan(&stats)
	if result.Error != nil {
		return nil, result.Error
	}
	return stats, nil
}

// GetClickTimeline 获取点击时间线
func (r *MySQLRepository) GetClickTimeline(ctx context.Context, shortCode string, hours int) ([]model.ClickTimeline, error) {
	var timeline []model.ClickTimeline

	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	query := `
		SELECT
			DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%S') as time_bucket,
			COUNT(*) as clicks,
			COUNT(DISTINCT ip_address) as unique_visitors
		FROM click_stats
		WHERE short_code = ? AND created_at >= ?
		GROUP BY time_bucket
		ORDER BY clicks DESC 
	`
	result := r.db.WithContext(ctx).Raw(query, shortCode, startTime).Scan(&timeline)
	if result.Error != nil {
		return nil, result.Error
	}
	return timeline, nil
}

// GetGeographicStats 获取地理统计
func (r *MySQLRepository) GetGeographicStats(ctx context.Context, shortCode string) ([]model.GeographicStats, error) {
	var stats []model.GeographicStats
	query := `
		SELECT
			country,
			region,
			city,
			COUNT(*) as clicks,
			COUNT(DISTINCT ip_address) as unique_visitors
	FROM click_stats
		WHERE short_code = ? AND country IS NOT NULL AND country != ''
		GROUP BY country, region, city
		ORDER BY clicks DESC
	`

	result := r.db.WithContext(ctx).Raw(query, shortCode).Scan(&stats)
	if result.Error != nil {
		return nil, result.Error
	}
	return stats, nil
}

// GetPlatformStats 获取平台统计
func (r *MySQLRepository) GetPlatformStats(ctx context.Context, shortCode string) ([]model.PlatformStats, error) {
	var stats []model.PlatformStats
	query := `
		SELECT
			device_type, COUNT(*) as clicks, COUNT(DISTINCT) as unique_visitors
		FROM click_stats
		WHERE short_code = ? 
		GROUP BY device_type
		ORDER BY clicks DESC
	`
	result := r.db.WithContext(ctx).Raw(query, shortCode).Scan(&stats)
	if result.Error != nil {
		return nil, result.Error
	}
	return stats, nil
}
