package model

type SummaryResponse struct {
	ShortCode   string           `json:"short_code"`
	TotalClicks int64            `json:"total_clicks"`
	DailyStats  []DailyStats     `json:"daily_stats,omitempty"`
	Referrers   map[string]int64 `json:"referrers,omitempty"`
	Countries   map[string]int64 `json:"countries,omitempty"`
	Devices     map[string]int64 `json:"devices,omitempty"`
}
