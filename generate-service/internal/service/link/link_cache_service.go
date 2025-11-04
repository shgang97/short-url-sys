package link

import (
	"generate-service/internal/model"
	"log"
	"time"

	"shared/constants"
)

// 异步发送缓存预热消息
func (s *linkService) sendWarmupAsync(link *model.Link) {
	go func() {
		eventID, _ := s.idGenerator.NextId()
		msg := model.CacheWarmupMessage{
			BaseMessage: model.BaseMessage{
				EventID:   string(eventID),
				EventType: "cache_warmup",
				Timestamp: time.Now(),
				Source:    "generate_service",
			},
			ShortCode:   link.ShortCode,
			OriginalURL: link.LongURL,
			ExpiredAt:   link.ExpiresAt,
			LogID:       link.ID,
		}
		if err := s.kafkaProducer.SendMessage(constants.TopicCacheWarmup, msg); err != nil {
			log.Printf("Failed to send cache warmup message: %v", err)
			// TODO 需要更新缓存预热记录表
		}
	}()
}

// 异步发送缓存更新消息
func (s *linkService) sendCacheUpdateAsync(link *model.Link) {
	go func() {
		eventID, _ := s.idGenerator.NextId()
		msg := model.CacheUpdateMessage{
			BaseMessage: model.BaseMessage{
				EventID:   string(eventID),
				EventType: "cache_warmup",
				Timestamp: time.Now(),
				Source:    "generate_service",
			},
			ShortCode:   link.ShortCode,
			OriginalURL: link.LongURL,
			ExpiredAt:   link.ExpiresAt,
			Status:      link.Status,
		}
		if err := s.kafkaProducer.SendMessage(constants.TopicCacheUpdate, msg); err != nil {
			log.Printf("Failed to send cache warmup message: %v", err)
			// TODO 需要更新缓存预热记录表
		}
	}()
}

// 异步发送缓存更新消息
func (s *linkService) sendCacheDeleteAsync(link *model.Link) {
	go func() {
		eventID, _ := s.idGenerator.NextId()
		msg := model.CacheDeleteMessage{
			BaseMessage: model.BaseMessage{
				EventID:   string(eventID),
				EventType: "cache_warmup",
				Timestamp: time.Now(),
				Source:    "generate_service",
			},
			ShortCode: link.ShortCode,
		}
		if err := s.kafkaProducer.SendMessage(constants.TopicCacheDelete, msg); err != nil {
			log.Printf("Failed to send cache warmup message: %v", err)
			// TODO 需要更新缓存预热记录表
		}
	}()
}
