package model

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
