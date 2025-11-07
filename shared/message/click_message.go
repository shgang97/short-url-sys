package message

import "time"

// ClickEventMessage 点击事件消息结构
type ClickEventMessage struct {
	BaseMessage
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Referer     string    `json:"referer,omitempty"`
	ClickTime   time.Time `json:"click_time"`
	ClickBy     string    `json:"click_by,omitempty"`
	Country     string    `json:"country,omitempty"`
	Region      string    `json:"region,omitempty"`
	City        string    `json:"city,omitempty"`
}
