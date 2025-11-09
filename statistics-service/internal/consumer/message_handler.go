package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"shared/constants"
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

func (r *HandlerRouter) Register(groupId, topic string, handler MessageHandler) {
	r.handlers[getHandlerKey(topic, groupId)] = handler
}

func (r *HandlerRouter) GetHandler(groupId, topic string) (MessageHandler, bool) {
	h, exist := r.handlers[getHandlerKey(topic, groupId)]
	return h, exist
}

func getHandlerKey(groupId, topic string) string {
	return fmt.Sprintf("%s_%s", groupId, topic)
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
		ClickBy:     msg.ClickBy,
	}
	if err := h.clickService.RecordClick(context.Background(), req); err != nil {
		logger.Logger.Error("failed to record click", zap.String("topic", topic), zap.Error(err))
		return false
	}
	return true
}

func CreatHandler(groupId, topic string, other interface{}) (MessageHandler, error) {
	key := getHandlerKey(groupId, topic)
	switch key {
	case getHandlerKey(constants.StatsGroupDetail, constants.TopicRecordClickEvent):
		service := other.(*click.Service)
		return NewRecordClickHandler(service), nil
	default:
		return nil, fmt.Errorf("unknown topic '%s'", topic)
	}
}
