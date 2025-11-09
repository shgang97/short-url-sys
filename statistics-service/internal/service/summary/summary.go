package summary

import (
	"context"
	"fmt"
	"statistics-service/internal/model"
	"statistics-service/internal/pkg/idgen"
	"statistics-service/internal/pkg/logger"
	sumRepo "statistics-service/internal/repository/summary"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	summaryRepo *sumRepo.Repository
	generator   idgen.Generator
}

func NewService(
	db *gorm.DB,
	sumRepo *sumRepo.Repository,
	generator idgen.Generator,
) *Service {
	return &Service{
		db:          db,
		summaryRepo: sumRepo,
		generator:   generator,
	}
}

func (s *Service) RecordSummary(ctx context.Context, handlerKey string, buffer map[string]int) error {
	logger.Logger.Info("RecordSummary in transaction")
	if len(buffer) == 0 {
		return nil
	}
	stats := s.convertToStats(buffer)
	batches := s.createBatches(stats, 100)

	// 使用单个事务处理所有数据
	fc := func(tx *gorm.DB) error {
		for i, batch := range batches {
			if err := s.summaryRepo.BatchUpsert(ctx, tx, handlerKey, batch); err != nil {
				logger.Logger.Error("Batch upsert failed in transaction",
					zap.Int("batch_index", i),
					zap.Int("batch_size", len(batch)),
					zap.Error(err))
				return fmt.Errorf("batch %d upsert failed in transaction: %w", i, err)
			}
		}
		logger.Logger.Info("All batches processed successfully in single transaction",
			zap.Int("total_batches", len(batches)), zap.Int("total_records", len(stats)))
		return nil
	}
	// 执行事务
	return s.db.Transaction(fc)
}

func (s *Service) convertToStats(buffer map[string]int) []*model.StatTotal {
	var stats []*model.StatTotal

	for key, count := range buffer {
		parts := strings.Split(key, "_")
		if len(parts) != 2 {
			logger.Logger.Warn("invalid key format, skipping", zap.String("key", key))
			continue
		}

		shortCode, date := parts[0], parts[1]
		id, _ := s.generator.NextId()
		stats = append(stats, &model.StatTotal{
			Id:          id,
			ShortCode:   shortCode,
			StatDate:    date,
			TotalClicks: count,
			CreatedBy:   "",
			UpdatedBy:   "",
		})
	}

	return stats
}

func (s *Service) createBatches(stats []*model.StatTotal, batchSize int) [][]*model.StatTotal {
	var batches [][]*model.StatTotal
	for i := 0; i < len(stats); i += batchSize {
		end := i + batchSize
		if end > len(stats) {
			end = len(stats)
		}
		batches = append(batches, stats[i:end])
	}
	return batches
}
