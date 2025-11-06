package consumer

import (
	"context"
	"encoding/json"
	"shared/message"
	"statistics-service/internal/pkg/logger"
	"statistics-service/internal/service/click"
	"time"

	"go.uber.org/zap"
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

type RecordClickHandler struct {
	clickService *click.Service
}

func NewRecordClickHandler(clickService *click.Service) *RecordClickHandler {
	return &RecordClickHandler{clickService: clickService}
}

func (h *RecordClickHandler) Handle(topic string, value []byte) bool {
	var msg message.ClickEventMessage
	if err := json.Unmarshal(value, &msg); err != nil {
		logger.Logger.Error("failed to unmarshal click message", zap.String("topic", topic), zap.Error(err))
		return false
	}
	clickTime := time.Now()
	if !msg.ClickTime.IsZero() {
		clickTime = msg.ClickTime
	}
	req := &click.RecordClickReq{
		ShortCode:   msg.ShortCode,
		OriginalURL: msg.OriginalURL,
		IP:          msg.IP,
		UserAgent:   msg.UserAgent,
		Referer:     msg.Referer,
		ClickTime:   clickTime,
		Country:     msg.Country,
		Region:      msg.Region,
		City:        msg.City,
	}
	if err := h.clickService.RecordClick(context.Background(), req); err != nil {
		logger.Logger.Error("failed to record click", zap.String("topic", topic), zap.Error(err))
		return false
	}
	return true
}

func CreatHandler(topic string, clickService *click.Service) (MessageHandler, error) {
	return NewRecordClickHandler(clickService), nil
}
