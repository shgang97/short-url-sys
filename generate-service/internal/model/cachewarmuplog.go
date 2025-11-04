package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type WarmupStatus string

const (
	WarmupStatusPending WarmupStatus = "pending"
	WarmupStatusSuccess WarmupStatus = "success"
	WarmupStatusFailed  WarmupStatus = "failed"
)

// Scan 实现数据库接口扫描
func (ws *WarmupStatus) Scan(value interface{}) error {
	if value == nil {
		*ws = WarmupStatusPending
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*ws = WarmupStatus(v)
	case string:
		*ws = WarmupStatus(v)
	default:
		return fmt.Errorf("unsupported type for WarmupStatus: %T", value)
	}
	return nil
}

// Value 实现数据库值接口
func (ws WarmupStatus) Value() (driver.Value, error) {
	return string(ws), nil
}

// CacheWarmupLog 缓存预热记录模型
type CacheWarmupLog struct {
	ID          uint64       `gorm:"primaryKey" json:"id"` // 使用雪花算法ID，非自增
	ShortCode   string       `gorm:"size:20;not null" json:"short_code"`
	WarmupTime  time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"warmup_time"`
	Status      WarmupStatus `gorm:"size:20;default:pending" json:"status"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string       `gorm:"size:100" json:"created_by,omitempty"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   string       `gorm:"size:100" json:"updated_by,omitempty"`
	Description string       `gorm:"size:100" json:"description,omitempty"`
	DeleteFlag  string       `gorm:"size:1;default:N" json:"delete_flag,omitempty"`
	Version     uint         `gorm:"default:0" json:"version"`
}

// TableName 指定表名
func (cwl *CacheWarmupLog) TableName() string {
	return "cache_warmup_logs"
}

// IsPending 检查是否处于待处理状态
func (cwl *CacheWarmupLog) IsPending() bool {
	return cwl.Status == WarmupStatusPending
}

// IsCompleted 检查是否已完成（成功或失败）
func (cwl *CacheWarmupLog) IsCompleted() bool {
	return cwl.Status == WarmupStatusSuccess || cwl.Status == WarmupStatusFailed
}
