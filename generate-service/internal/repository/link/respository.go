package link

import (
	"context"
	"generate-service/internal/model"
)

// Repository 链接数据访问接口
type Repository interface {
	// Create 创建链接
	Create(ctx context.Context, link *model.Link) error

	// FindByShortCode 查询链接
	FindByShortCode(ctx context.Context, shortCode string) (*model.Link, error)
	FindByLongURL(ctx context.Context, longURL string) (*model.Link, error)
	Exists(ctx context.Context, shortCode string) (bool, error)

	// Update 更新链接
	Update(ctx context.Context, link *model.Link) error
	UpdateClickCount(ctx context.Context, shortCode string, increment int64) error

	// Delete 删除链接
	Delete(ctx context.Context, link *model.Link) error

	// List 列表查询
	List(ctx context.Context, filter ListFilter, page, pageSize int) ([]model.Link, int64, error)

	// BatchCreate 批量创建
	BatchCreate(ctx context.Context, links []model.Link) ([]model.Link, error)

	// CleanupExpired 清理过期链接
	CleanupExpired(ctx context.Context) (int64, error)
}

type ListFilter struct {
	CreatedBy string
	Status    string
	Search    string
}
