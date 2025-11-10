package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"shared/message"
	"statistics-service/internal/pkg/logger"
	"statistics-service/internal/service/click"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// MessageHandler 消息处理器接口（每个Topic对应一个实现）
type MessageHandler interface {
	// Handle 处理消息
	// topic: 消息所属的Topic
	// value: 消息体（原始字节）
	// return: 处理是否成功
	Handle(topic string, msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)
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
	r.handlers[GetHandlerKey(topic, groupId)] = handler
}

func (r *HandlerRouter) GetHandler(groupId, topic string) (MessageHandler, bool) {
	h, exist := r.handlers[GetHandlerKey(topic, groupId)]
	return h, exist
}

func GetHandlerKey(groupId, topic string) string {
	return fmt.Sprintf("%s_%s", groupId, topic)
}

type RecordClickHandler struct {
	clickService *click.Service
}

func NewRecordClickHandler(clickService *click.Service) *RecordClickHandler {
	return &RecordClickHandler{clickService: clickService}
}

func (h *RecordClickHandler) Handle(topic string, msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	var clickMsg message.ClickEventMessage
	if err := json.Unmarshal(msg.Value, &clickMsg); err != nil {
		logger.Logger.Error("failed to unmarshal click message", zap.String("topic", topic), zap.Error(err))
		session.MarkMessage(msg, "")
		return
	}
	clickTime := time.Now()
	if !clickMsg.ClickTime.IsZero() {
		clickTime = clickMsg.ClickTime
	}
	req := &click.RecordClickReq{
		ShortCode:   clickMsg.ShortCode,
		OriginalURL: clickMsg.OriginalURL,
		IP:          clickMsg.IP,
		UserAgent:   clickMsg.UserAgent,
		Referer:     clickMsg.Referer,
		ClickTime:   clickTime,
		Country:     clickMsg.Country,
		Region:      clickMsg.Region,
		City:        clickMsg.City,
		ClickBy:     clickMsg.ClickBy,
	}
	if err := h.clickService.RecordClick(context.Background(), req); err != nil {
		logger.Logger.Error("failed to record click", zap.String("topic", topic), zap.Error(err))
		session.MarkMessage(msg, "")
		return
	}
	session.MarkMessage(msg, "")
}
