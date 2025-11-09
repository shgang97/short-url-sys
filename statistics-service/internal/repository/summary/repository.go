package summary

import (
	"context"
	"statistics-service/internal/model"
	"statistics-service/internal/pkg/logger"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) BatchUpsert(ctx context.Context, tx *gorm.DB, username string, stats []*model.StatTotal) error {
	logger.Logger.Info("Batch upsert in transaction")
	if len(stats) == 0 {
		return nil
	}
	query := `
		INSERT INTO click_stats_summary (id, short_code, stat_date, total_clicks, created_by, created_at, updated_by, updated_at, version) 
		VALUES `

	now := time.Now()
	var values []interface{}
	var valuePlaceholders []string
	for _, stat := range stats {
		valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		values = append(values,
			stat.Id,
			stat.ShortCode,
			stat.StatDate,
			stat.TotalClicks,
			username,
			now,
			username,
			now,
			0,
		)
	}
	query += strings.Join(valuePlaceholders, ",")
	query += ` ON DUPLICATE KEY UPDATE
		total_clicks = total_clicks + VALUES(total_clicks),
		updated_by = ?,
		updated_at = ?,
		version = version + 1 `
	values = append(values, username, now)

	return tx.WithContext(ctx).Exec(query, values...).Error
}
