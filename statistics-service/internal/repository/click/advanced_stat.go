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
