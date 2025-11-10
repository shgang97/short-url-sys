package model

type SummaryResponse struct {
	ShortCode   string           `json:"short_code"`
	TotalClicks int64            `json:"total_clicks"`
	DailyStats  []DailyStats     `json:"daily_stats,omitempty"`
	Referrers   map[string]int64 `json:"referrers,omitempty"`
	Countries   map[string]int64 `json:"countries,omitempty"`
	Devices     map[string]int64 `json:"devices,omitempty"`
	Browsers    map[string]int64 `json:"browsers,omitempty"`
	Systems     map[string]int64 `json:"systems,omitempty"`
}

// TimeSeriesData 时间序列数据点
type TimeSeriesData struct {
	Period         string `json:"period"`
	Clicks         int64  `json:"clicks"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

// TimeSeriesResponse 时间序列响应
type TimeSeriesResponse struct {
	ShortCode  string            `json:"short_code"`
	Unit       string            `json:"unit"`
	TimeSeries []*TimeSeriesData `json:"time_series"`
	Summary    *SummaryData      `json:"summary,omitempty"`
}

// SummaryData 汇总数据
type SummaryData struct {
	TotalClicks         int64 `json:"total_clicks"`
	TotalUniqueVisitors int64 `json:"total_unique_visitors"`
}
