package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"redirect-service/internal/service/cache"
	"shared/constants"
	"shared/message"
)

// MessageHandler 消息处理器接口（每个Topic对应一个实现）
type MessageHandler interface {
	// Handle 处理消息
	// topic: 消息所属的Topic
	// value: 消息体（原始字节）
	// return: 处理是否成功
	Handle(topic string, value []byte) bool
}

type HandlerRouter struct {
	handlers map[string]MessageHandler
}

func NewHandlerRouter() *HandlerRouter {
	return &HandlerRouter{
		handlers: make(map[string]MessageHandler),
	}
}

func (r *HandlerRouter) Register(topic string, handler MessageHandler) {
	r.handlers[topic] = handler
}

func (r *HandlerRouter) GetHandler(topic string) (MessageHandler, bool) {
	return r.handlers[topic], true
}

type CacheWarmupHandler struct {
	cacheService *cache.Service
}

func NewCacheWarmupHandler(cacheService *cache.Service) *CacheWarmupHandler {
	return &CacheWarmupHandler{cacheService: cacheService}
}

func (c *CacheWarmupHandler) Handle(topic string, value []byte) bool {
	var msg message.CacheWarmupMessage
	if err := json.Unmarshal(value, &msg); err != nil {
		return false
	}
	err := c.cacheService.SetShortUrl(context.Background(), msg.ShortCode, msg.OriginalURL, msg.ExpiredAt)
	if err != nil {
		return false
	}
	return true
}

type CacheUpdateHandler struct {
	cacheService *cache.Service
}

func NewCacheUpdateHandler(cacheService *cache.Service) *CacheUpdateHandler {
	return &CacheUpdateHandler{cacheService: cacheService}
}

func (c *CacheUpdateHandler) Handle(topic string, value []byte) bool {
	var msg message.CacheUpdateMessage
	if err := json.Unmarshal(value, &msg); err != nil {
		return false
	}
	err := c.cacheService.SetShortUrl(context.Background(), msg.ShortCode, msg.OriginalURL, msg.ExpiredAt)
	if err != nil {
		return false
	}
	return true
}

type CacheDeleteHandler struct {
	cacheService *cache.Service
}

func NewCacheDeleteHandler(cacheService *cache.Service) *CacheDeleteHandler {
	return &CacheDeleteHandler{cacheService: cacheService}
}

func (c *CacheDeleteHandler) Handle(topic string, value []byte) bool {
	var msg message.CacheDeleteMessage
	if err := json.Unmarshal(value, &msg); err != nil {
		return false
	}
	err := c.cacheService.DelShortUrl(context.Background(), msg.ShortCode)
	if err != nil {
		return false
	}
	return true
}

func CreateHandler(topic string, cacheService *cache.Service) (MessageHandler, error) {
	switch topic {
	case constants.TopicCacheWarmup:
		return NewCacheWarmupHandler(cacheService), nil
	case constants.TopicCacheUpdate:
		return NewCacheUpdateHandler(cacheService), nil
	case constants.TopicCacheDelete:
		return NewCacheDeleteHandler(cacheService), nil
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}
}
