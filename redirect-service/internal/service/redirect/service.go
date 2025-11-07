package redirect

import (
	"context"
	"log"
	"redirect-service/internal/client/grpc/generate"
	"redirect-service/internal/pkg/idgen"
	"redirect-service/internal/producer"
	"redirect-service/internal/repository/cache"
	"redirect-service/internal/service/geoip"
	"shared/constants"
	"shared/message"
	"strconv"
	"time"
)

type Service struct {
	genClient     *generate.Client
	cacheRepo     *cache.Repository
	kafkaProducer *producer.KafkaProducer
	geoIPSvc      geoip.Service
	generator     idgen.Generator
}

func NewService(
	client *generate.Client,
	cacheRepo *cache.Repository,
	kafkaProducer *producer.KafkaProducer,
	geoIpSvc geoip.Service,
	generator idgen.Generator,
) *Service {
	return &Service{
		genClient:     client,
		cacheRepo:     cacheRepo,
		kafkaProducer: kafkaProducer,
		geoIPSvc:      geoIpSvc,
		generator:     generator,
	}
}

func (s *Service) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	// 从缓存获取长链接
	longUrl, err := s.cacheRepo.GetOriginalURL(ctx, shortCode)
	if err == nil {
		return longUrl, nil
	}
	// 缓存未命中，回溯到generate-service服务
	resp, err := s.genClient.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	// 加入缓存
	var expiredAt *time.Time
	if resp.ExpireTime != nil {
		t := resp.ExpireTime.AsTime()
		expiredAt = &t
	}
	go func() {
		if err = s.cacheRepo.SetShortURL(ctx, shortCode, resp.OriginalUrl, expiredAt); err != nil {
			log.Printf("failed to cache short url: %v", err)
		}
	}()
	return resp.OriginalUrl, nil
}

func (s *Service) RecordClick(ctx context.Context, shortCode string, req *RedirectRequest) error {
	now := time.Now()
	msg := &message.ClickEventMessage{
		BaseMessage: message.BaseMessage{
			EventType: "click_event",
			Timestamp: now,
			Source:    "redirect-service",
		},
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		IP:          req.IPAddress,
		UserAgent:   req.UserAgent,
		Referer:     req.Referer,
		ClickTime:   now,
		ClickBy:     req.Username,
	}
	eventId, err := s.generator.NextId()
	if err == nil {
		msg.EventID = strconv.FormatUint(eventId, 10)
	}
	getGeoInfo, err := s.geoIPSvc.GetGeoInfo(req.IPAddress)
	if err == nil {
		msg.Country = getGeoInfo.Country
		msg.Region = getGeoInfo.Region
		msg.City = getGeoInfo.City
	}
	return s.kafkaProducer.SendClickMessage(constants.TopicRecordClickEvent, msg)
}
