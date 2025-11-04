package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LinkStatus string

const (
	LinkStatusActive   LinkStatus = "active"
	LinkStatusDisabled LinkStatus = "disabled"
	LinkStatusExpired  LinkStatus = "expired"
)

// Scan 实现数据库接口扫描
func (ls *LinkStatus) Scan(value interface{}) error {
	if value == nil {
		*ls = LinkStatusActive
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*ls = LinkStatus(v)
	case string:
		*ls = LinkStatus(v)
	default:
		return fmt.Errorf("unsupported type for LinkStatus: %T", value)
	}
	return nil
}

// Value 实现数据库值接口
func (ls LinkStatus) Value() (driver.Value, error) {
	return string(ls), nil
}

// Link 短链接模型
type Link struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	ShortCode   string     `gorm:"size:10;not null;uniqueIndex" json:"short_code"`
	LongURL     string     `gorm:"type:text;not null" json:"long_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ClickCount  int64      `gorm:"default:0" json:"click_count"`
	Status      LinkStatus `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string     `gorm:"size:100" json:"created_by,omitempty"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   string     `gorm:"size:100" json:"updated_by,omitempty"`
	Description string     `gorm:"size:500" json:"description,omitempty"`
	DeleteFlag  string     `gorm:"size:1" json:"delete_flag,omitempty"`
	Version     uint       `gorm:"default:0" json:"version"`
}

// TableName 指定表名
func (l *Link) TableName() string {
	return "links"
}

// IsActive 检查链接是否有效
func (l *Link) IsActive() bool {
	if l.Status != LinkStatusActive {
		return false
	}

	if l.ExpiresAt != nil && l.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}
