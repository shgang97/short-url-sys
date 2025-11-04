package model

import "time"

// CreateShortRequest 创建短链请求
type CreateShortRequest struct {
	LongURL     string     `json:"long_url" binding:"required,url"`
	CustomCode  *string    `json:"custom_code,omitempty" binding:"omitempty,alphanum,min=3,max=10"` // 使用指针类型，区分“未设置”和“设置”，指针为 nil，表示客户端没有提供该字段
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=500"`
	CreatedBy   *string    `json:"created_by,omitempty" binding:"omitempty,max=100"`
}

// BatchCreateRequest 批量创建短链请求
type BatchCreateRequest struct {
	URLs []BatchURLItem `json:"urls" binding:"required,min=1,max=100"`
}

type BatchURLItem struct {
	LongURL    string  `json:"long_url" binding:"required,url"`
	CustomCode *string `json:"custom_code,omitempty" binding:"omitempty,alphanum,min=3,max=10"`
}

// UpdateLinkRequest 更新链接请求
type UpdateLinkRequest struct {
	LongURL     *string    `json:"long_url" binding:"required,url" binding:"omitempty,url"` // 使用指针类型，区分“未设置”和“设置”
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Status      *string    `json:"status,omitempty" binding:"omitempty,oneof=active disabled"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=500"`
}

// ListLinksRequest 列表查询请求
type ListLinksRequest struct {
	Page      int     `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int     `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	CreatedBy *string `form:"created_by,omitempty" binding:"omitempty,max=100"`
	Status    *string `form:"status,omitempty" binding:"omitempty,oneof=active disabled expired"`
}
