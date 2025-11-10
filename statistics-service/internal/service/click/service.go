package click

import (
	"context"
	"statistics-service/internal/model"
	"statistics-service/internal/pkg/idgen"
	"statistics-service/internal/repository/click"
	"statistics-service/internal/service/device_detector"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	clickRepo click.Repository
	generator idgen.Generator
	detector  detector.DeviceDetector
}

func NewService(
	db *gorm.DB,
	repo click.Repository,
	generator idgen.Generator,
	detector detector.DeviceDetector,
) *Service {
	return &Service{
		db:        db,
		clickRepo: repo,
		generator: generator,
		detector:  detector,
	}
}

func (s *Service) RecordClick(ctx context.Context, req *RecordClickReq) error {
	id, _ := s.generator.NextId()
	clc := &model.ClickEvent{
		ID:          id,
		ShortCode:   req.ShortCode,
		OriginalURL: req.OriginalURL,
		IP:          req.IP,
		UserAgent:   req.UserAgent,
		Referer:     req.Referer,
		Country:     req.Country,
		Region:      req.Region,
		City:        req.City,
		ClickTime:   req.ClickTime,
	}
	deviceInfo, err := s.detector.Parse(req.UserAgent)
	if err == nil {
		clc.DeviceType = deviceInfo.DeviceType
		clc.Browser = deviceInfo.Browser
		clc.OS = deviceInfo.OS
	}
	s.setDefaultValues(ctx, clc, time.Now(), req.ClickBy)
	return s.clickRepo.Create(ctx, clc)
}

func (s *Service) setDefaultValues(ctx context.Context, clc *model.ClickEvent, now time.Time, username string) {
	clc.CreatedAt = now
	clc.CreatedBy = username // 当前接口只有Kafka消费者调用，用户名由Kafka传递
	clc.UpdatedAt = now
	clc.UpdatedBy = username
	clc.DeleteFlag = "N"
	clc.Version = 0
}

func (s *Service) GetStatsSummary(
	ctx context.Context,
	shortCode string,
	startDate, endDate *time.Time,
) (*model.SummaryResponse, error) {
	summary, err := s.clickRepo.GetStatsSummary(ctx, shortCode, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return &model.SummaryResponse{
		ShortCode:   shortCode,
		TotalClicks: summary.TotalClicks,
		DailyStats:  summary.DailyStats,
		Referrers:   summary.Referrers,
		Countries:   summary.Countries,
		Devices:     summary.Devices,
		Browsers:    summary.Browsers,
		Systems:     summary.Systems,
	}, nil
}

func (s *Service) GetTimeSeriesSummary(
	ctx context.Context,
	shortCode string,
	startTime *time.Time,
	endTime *time.Time,
	groupExpr string,
	periodExpr string,
) (*model.TimeSeriesResponse, error) {
	seriesStats, err := s.clickRepo.GetClickTimeline(ctx, shortCode, startTime, endTime, groupExpr, periodExpr)
	if err != nil {
		return nil, err
	}
	timeSeries := make([]*model.TimeSeriesData, 0, len(seriesStats))
	var totalClicks, totalUniqueClicks int64
	for _, stat := range seriesStats {
		timeSeries = append(timeSeries, &model.TimeSeriesData{
			Period:         stat.Period,
			Clicks:         stat.Clicks,
			UniqueVisitors: stat.UniqueVisitors,
		})
		totalClicks += stat.Clicks
		totalUniqueClicks += stat.UniqueVisitors
	}
	resp := &model.TimeSeriesResponse{
		ShortCode:  shortCode,
		TimeSeries: timeSeries,
		Summary: &model.SummaryData{
			TotalClicks:         totalClicks,
			TotalUniqueVisitors: totalUniqueClicks,
		},
	}
	return resp, nil
}
