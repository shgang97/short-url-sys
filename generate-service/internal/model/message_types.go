package model

import (
	"fmt"
	"time"
)

type BaseMessage struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

func (m BaseMessage) GetEventType() string {
	return m.EventType
}

// CacheWarmupMessage 缓存预热消息
type CacheWarmupMessage struct {
	BaseMessage
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ExpiredAt   *time.Time `json:"expired_at"`
	LogID       uint64     `json:"log_id"`
}

func (m CacheWarmupMessage) GetKey() string {
	return m.ShortCode
}

func (m CacheWarmupMessage) Validate() error {
	if m.ShortCode == "" || m.OriginalURL == "" {
		return fmt.Errorf("short_code and original_url are required")
	}
	return nil
}

// CacheUpdateMessage 缓存更新消息
type CacheUpdateMessage struct {
	BaseMessage
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ExpiredAt   *time.Time `json:"expired_at,omitempty"`
	Status      LinkStatus `json:"status"`
}

func (m CacheUpdateMessage) GetKey() string {
	return m.ShortCode
}

func (m CacheUpdateMessage) Validate() error {
	if m.ShortCode == "" {
		return fmt.Errorf("short_code is required")
	}
	return nil
}

// CacheDeleteMessage 缓存删除消息
type CacheDeleteMessage struct {
	BaseMessage
	ShortCode string `json:"short_code"`
	Reason    string `json:"reason"`
}

func (m CacheDeleteMessage) GetKey() string {
	return m.ShortCode
}

func (m CacheDeleteMessage) Validate() error {
	if m.ShortCode == "" {
		return fmt.Errorf("short_code is required")
	}
	return nil
}
