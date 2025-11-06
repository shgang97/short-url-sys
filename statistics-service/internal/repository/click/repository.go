package click

import (
	"context"
	"fmt"
	"shared/errors"
	"statistics-service/internal/model"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	// Create 创建点击事件
	Create(ctx context.Context, event *model.ClickEvent) error
}

type repository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, clt *model.ClickEvent) error {
	if clt == nil {
		return fmt.Errorf("click event is nil")
	}

	err := r.db.WithContext(ctx).Create(clt).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.ErrShortCodeExists
		}
		return &errors.RepositoryError{Operation: "create", Err: err}
	}
	return nil
}
