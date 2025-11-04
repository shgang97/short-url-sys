package link

import (
	"context"
	"generate-service/internal/model"
	"generate-service/internal/pkg/errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MySQLRepository struct {
	db *gorm.DB
}

func (r *MySQLRepository) Create(ctx context.Context, link *model.Link) error {
	result := r.db.WithContext(ctx).Create(link)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return errors.ErrShortCodeExists
		}
		return &errors.RepositoryError{Operation: "Create", Err: result.Error}
	}
	return nil
}

func (r *MySQLRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	var link model.Link
	result := r.db.WithContext(ctx).Where("short_code=? and delete_flag = 'N'", shortCode).First(&link)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errors.ErrLinkNotFound
		}
		return nil, &errors.RepositoryError{Operation: "FindByShortCode", Err: result.Error}
	}
	return &link, nil
}

func (r *MySQLRepository) FindByLongURL(ctx context.Context, longURL string) (*model.Link, error) {
	var link model.Link
	result := r.db.WithContext(ctx).Where("long_url=? and delete_flag = 'N'", longURL).First(&link)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errors.ErrLinkNotFound
		}
		return nil, &errors.RepositoryError{Operation: "FindByLongURL", Err: result.Error}
	}
	return &link, nil
}

func (r *MySQLRepository) Exists(ctx context.Context, shortCode string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Link{}).Where("short_code=?", shortCode).Count(&count)
	if result.Error != nil {
		return false, &errors.RepositoryError{Operation: "Exists", Err: result.Error}
	}
	return count > 0, nil
}

func (r *MySQLRepository) Update(ctx context.Context, link *model.Link) error {
	result := r.db.WithContext(ctx).Save(link)
	if result.Error != nil {
		return &errors.RepositoryError{Operation: "Update", Err: result.Error}
	}
	return nil
}

func (r *MySQLRepository) UpdateClickCount(ctx context.Context, shortCode string, increment int64) error {
	result := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("short_code=?", shortCode).
		Update("click_count", gorm.Expr("click_count + ?", increment))
	if result.Error != nil {
		return &errors.RepositoryError{Operation: "UpdateClickCount", Err: result.Error}
	}
	if result.RowsAffected == 0 {
		return errors.ErrLinkNotFound
	}
	return nil
}

func (r *MySQLRepository) Delete(ctx context.Context, link *model.Link) error {
	result := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("id = ? AND version = ?", link.ID, link.Version).
		Updates(map[string]interface{}{
			"delete_flag": "Y",
			"version":     link.Version + 1,
			"updated_at":  time.Now(),
			"updated_by":  link.UpdatedBy,
		})
	if result.Error != nil {
		return &errors.RepositoryError{Operation: "Delete", Err: result.Error}
	}
	if result.RowsAffected == 0 {
		return errors.ErrLinkNotFound
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context, filter ListFilter, page, pageSize int) ([]model.Link, int64, error) {
	var links []model.Link
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Link{})

	// 应用过滤器
	if filter.CreatedBy != "" {
		query = query.Where("created_by=?", filter.CreatedBy)
	}
	if filter.Status != "" {
		query = query.Where("status=?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("short_code LIKE ? OR long_url LIKE ? ", search, search)
	}
	query = query.Where("delete_flag = 'N'")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, &errors.RepositoryError{Operation: "ListCount", Err: err}
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&links)
	if result.Error != nil {
		return nil, 0, &errors.RepositoryError{Operation: "List", Err: result.Error}
	}
	return links, total, nil
}

func (r *MySQLRepository) BatchCreate(ctx context.Context, links []model.Link) ([]model.Link, error) {
	result := r.db.WithContext(ctx).Create(links)
	if result.Error != nil {
		return nil, &errors.RepositoryError{Operation: "BatchCreate", Err: result.Error}
	}
	return links, nil
}

func (r *MySQLRepository) CleanupExpired(ctx context.Context) (int64, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.Link{}).
		Where("expired_at IS NOT NULL AND expires_at < ? AND status = ?", now, model.LinkStatusActive).
		Update("status", model.LinkStatusExpired)

	if result.Error != nil {
		return 0, &errors.RepositoryError{Operation: "CleanupExpired", Err: result.Error}
	}
	return result.RowsAffected, nil
}

func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}
