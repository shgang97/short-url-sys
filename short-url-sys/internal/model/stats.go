package model

import "time"

// LinkStats 链接统计信息
type LinkStats struct {
	ShortCode    string    `json:"short_code"`
	LongUrl      string    `json:"long_url"`
	CreatedAt    time.Time `json:"created_at"`
	ClickCount   int64     `json:"click_count"`
	RecentClicks int64     `json:"recent_clicks"`
}

// ClickTimeline 点击时间线
type ClickTimeline struct {
	TimeBucket     string `json:"time_bucket"`
	Clicks         int64  `json:"clicks"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

// GeographicStats 地理统计
type GeographicStats struct {
	Country        string `json:"country"`
	Region         string `json:"region"`
	City           string `json:"city"`
	Clicks         int64  `json:"clicks"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

type PlatformStats struct {
	DeviceType     string `json:"device_type"`
	Clicks         int64  `json:"clicks"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

// SystemStats 系统统计
type SystemStats struct {
	TotalLinks  int64 `json:"total_links"`
	TotalClicks int64 `json:"total_clicks"`
	ActiveLinks int64 `json:"active_links"`
	TodayClicks int64 `json:"today_clicks"`
	WeekClicks  int64 `json:"week_clicks"`
	MonthClicks int64 `json:"month_clicks"`
}

type ReferrerStats struct {
	Referrer       string `json:"referrer"`
	Clicks         int64  `json:"clicks"`
	UniqueVisitors int64  `json:"unique_visitors"`
}
