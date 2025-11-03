package stats

import (
	"context"
	"short-url-sys/internal/model"
	"time"
)

// CountLinks 统计总链接数
func (r *MySQLRepository) CountLinks(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Link{}).Count(&count)
	return count, result.Error
}

// CountClicks 统计点击总量
func (r *MySQLRepository) CountClicks(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.ClickStats{}).Count(&count)
	return count, result.Error
}

// CountActiveLinks 统计活跃链接数
func (r *MySQLRepository) CountActiveLinks(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("status = ?", model.LinkStatusActive).
		Count(&count)
	return count, result.Error
}

// CountClicksSince 统计指定时间后的点击量
func (r *MySQLRepository) CountClicksSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.ClickStats{}).
		Where("created_at >= ?", since).
		Count(&count)
	return count, result.Error
}
