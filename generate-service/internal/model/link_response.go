package model

import "time"

// CreateShortResponse 创建短链接响应
type CreateShortResponse struct {
	ShortURL  string     `json:"short_url"`
	LongURL   string     `json:"long_url"`
	ShortCode string     `json:"short_code"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	QRCodeURL string     `json:"qr_code_url,omitempty"`
}

// LinkInfoResponse 链接信息响应
type LinkInfoResponse struct {
	ShortCode    string     `json:"short_code"`
	LongURL      string     `json:"long_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"` // nil值永不过期
	ClickCount   int64      `json:"click_count"`
	LastAccessed *time.Time `json:"last_accessed,omitempty"` // nil值从未被访问
	Status       string     `json:"status"`
	Description  string     `json:"description,omitempty"`
}

// BatchCreateResponse 批量创建响应
type BatchCreateResponse struct {
	Results []BatchResult `json:"results"`
	Failed  []BatchFailed `json:"failed,omitempty"`
}

type BatchResult struct {
	LongURL   string `json:"long_url"`
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}

type BatchFailed struct {
	LongURL string `json:"long_url"`
	Error   string `json:"error"`
}

// ListLinksResponse 链接列表响应
type ListLinksResponse struct {
	Links    []LinkInfoResponse `json:"links"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Pages    int                `json:"pages"`
}
