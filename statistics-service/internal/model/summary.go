package model

import (
	"time"
)

// ClickStatsSummary 点击统计汇总表(按天)
type ClickStatsSummary struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShortCode      string    `gorm:"size:20;not null;uniqueIndex:uk_short_code_date" json:"shortCode"`
	StatDate       time.Time `gorm:"type:date;not null;uniqueIndex:uk_short_code_date" json:"statDate"`
	TotalClicks    int       `gorm:"not null;default:0" json:"totalClicks"`
	UniqueVisitors int       `gorm:"not null;default:0" json:"uniqueVisitors"`
	MobileClicks   int       `gorm:"not null;default:0" json:"mobileClicks"`
	DesktopClicks  int       `gorm:"not null;default:0" json:"desktopClicks"`
	TabletClicks   int       `gorm:"not null;default:0" json:"tabletClicks"`
	TopCountry     string    `gorm:"size:100" json:"topCountry"`
	TopRegion      string    `gorm:"size:100" json:"topRegion"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreatedBy      string    `gorm:"size:100" json:"createdBy"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdatedBy      string    `gorm:"size:100" json:"updatedBy"`
	Description    string    `gorm:"size:100" json:"description"`
	DeleteFlag     string    `gorm:"size:1;default:N" json:"deleteFlag"`
	Version        uint      `gorm:"default:0" json:"version"`
}

// TableName 指定表名
func (ClickStatsSummary) TableName() string {
	return "click_stats_summary"
}

// StatTotal 实时计算点击总量数据模型
type StatTotal struct {
	Id          uint64
	ShortCode   string
	StatDate    string
	TotalClicks int
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
	Version     int
}

// DailyStats 每日统计摘要
type DailyStats struct {
	Date      string `json:"date"`
	Clicks    int64  `json:"clicks"`
	UniqueIPs int64  `json:"unique_ips"`
}

type SummaryStats struct {
	TotalClicks int64            `json:"total_clicks"`
	DailyStats  []DailyStats     `json:"daily_stats,omitempty"`
	Referrers   map[string]int64 `json:"referrers,omitempty"`
	Countries   map[string]int64 `json:"countries,omitempty"`
	Devices     map[string]int64 `json:"devices,omitempty"`
}
