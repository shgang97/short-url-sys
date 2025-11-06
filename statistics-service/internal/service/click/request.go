package click

import "time"

type RecordClickReq struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Referer     string    `json:"referer,omitempty"`
	ClickTime   time.Time `json:"click_time"`
	Country     string    `json:"country,omitempty"`
	Region      string    `json:"region,omitempty"`
	City        string    `json:"city,omitempty"`
	Source      string    `json:"source,omitempty"`
	ClickBy     string    `json:"click_by,omitempty"`
}
