package click

import (
	"context"
	"fmt"
	"shared/errors"
	"statistics-service/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Create 创建点击事件
	Create(ctx context.Context, event *model.ClickEvent) error
	GetStatsSummary(ctx context.Context, shortCode string, startDate, endDate *time.Time) (*model.SummaryStats, error)
}

type repository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, clt *model.ClickEvent) error {
	if clt == nil {
		return fmt.Errorf("click event is nil")
	}

	err := r.db.WithContext(ctx).Create(clt).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.ErrShortCodeExists
		}
		return &errors.RepositoryError{Operation: "create", Err: err}
	}
	return nil
}

func (r *repository) GetStatsSummary(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (*model.SummaryStats, error) {
	var stats model.SummaryStats

	// 基础查询
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Where("short_code = ? AND delete_flag = 'N'", shortCode)

	// 时间范围过滤
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}

	// 获取总点击量和独立IP数
	var totalClicks int64

	err := query.
		Select("COUNT(*) AS total_clicks").
		Scan(&totalClicks).Error
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.TotalClicks = totalClicks

	// 获取每日统计
	dailyStats, err := r.getDailyStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.DailyStats = dailyStats

	// 获取来源统计
	referrerStats, err := r.getReferrerStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.Referrers = referrerStats

	// 获取设备统计
	countryStats, err := r.getCountryStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.Countries = countryStats

	// 获取国家统计
	deviceStats, err := r.getDeviceStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.Devices = deviceStats

	// 获取浏览器统计
	browserStats, err := r.getBrowserStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.Browsers = browserStats

	// 获取操作系统统计
	systemStats, err := r.getOsStats(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, &errors.RepositoryError{Operation: "GetStatsSummary", Err: err}
	}
	stats.Systems = systemStats

	return &stats, nil
}

func (r *repository) getDailyStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) ([]model.DailyStats, error) {
	var stats []model.DailyStats

	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select(`DATE(click_time) as date, COUNT(*) as clicks, COUNT(DISTINCT ip) as unique_ips`).
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}

	err := query.
		Group("DATE(click_time)").
		Order("date DESC").
		Find(&stats).Error
	return stats, err
}

func (r *repository) getReferrerStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (map[string]int64, error) {
	var referrerStats []struct {
		Referer string
		Count   int64
	}
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("COALESCE(referer, 'direct') as referer, COUNT(*) as count").
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}
	err := query.
		Group("COALESCE(referer, 'direct')").
		Order("count DESC").
		Find(&referrerStats).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, stat := range referrerStats {
		stats[stat.Referer] = stat.Count
	}
	return stats, nil
}

// 获取国家统计
func (r *repository) getCountryStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (map[string]int64, error) {
	var countryStats []struct {
		Country string
		Count   int64
	}
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("COALESCE(country, 'unknown') as country, COUNT(*) as count").
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}
	err := query.
		Group("COALESCE(country, 'unknown')").
		Order("count DESC").
		Find(&countryStats).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, stat := range countryStats {
		stats[stat.Country] = stat.Count
	}
	return stats, nil
}

// 获取设备统计
func (r *repository) getDeviceStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (map[string]int64, error) {
	var deviceStats []struct {
		Device string
		Count  int64
	}
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("COALESCE(device_type, 'other') as device, COUNT(*) as count").
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}
	err := query.
		Group("COALESCE(device_type, 'other')").
		Order("count DESC").
		Find(&deviceStats).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, stat := range deviceStats {
		stats[stat.Device] = stat.Count
	}
	return stats, nil
}

// 获取浏览器统计
func (r *repository) getBrowserStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (map[string]int64, error) {
	var osStats []struct {
		Browser string
		Count   int64
	}
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("COALESCE(browser, 'other') as browser, COUNT(*) as count").
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}
	err := query.
		Group("COALESCE(browser, 'other')").
		Order("count DESC").
		Find(&osStats).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, stat := range osStats {
		stats[stat.Browser] = stat.Count
	}
	return stats, nil
}

// 获取操作系统统计
func (r *repository) getOsStats(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (map[string]int64, error) {
	var osStats []struct {
		Os    string
		Count int64
	}
	query := r.db.WithContext(ctx).Model(&model.ClickEvent{}).
		Select("COALESCE(os, 'other') as os, COUNT(*) as count").
		Where("short_code = ? AND delete_flag = 'N'", shortCode)
	if startDate != nil {
		query = query.Where("click_time >= ?", startDate.Format("2006-01-02 15:04:05"))
	}
	if endDate != nil {
		query = query.Where("click_time <= ?", endDate.Format("2006-01-02 15:04:05"))
	}
	err := query.
		Group("COALESCE(os, 'other')").
		Order("count DESC").
		Find(&osStats).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, stat := range osStats {
		stats[stat.Os] = stat.Count
	}
	return stats, nil
}
