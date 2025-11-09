package consumer

import (
	"context"
	"encoding/json"
	"shared/message"
	"statistics-service/internal/pkg/logger"
	"statistics-service/internal/service/summary"
	"sync"
	"time"

	"go.uber.org/zap"
)

type SummaryHandler struct {
	handlerKey     string
	summaryService *summary.Service
	buffer         map[string]int
	bufferSize     int
	wg             *sync.WaitGroup
	mu             sync.Mutex
}

func NewSummaryHandler(
	handlerKey string,
	summaryService *summary.Service,
	bufferSize int,
) *SummaryHandler {
	return &SummaryHandler{
		handlerKey:     handlerKey,
		summaryService: summaryService,
		buffer:         make(map[string]int),
		bufferSize:     bufferSize,
	}
}

func (h *SummaryHandler) Handle(topic string, value []byte) bool {
	var msg message.ClickEventMessage
	if err := json.Unmarshal(value, &msg); err != nil {
		logger.Logger.Error("failed to unmarshal click message",
			zap.String("topic", topic), zap.Error(err))
		return false
	}

	// 生成聚合key: shortCode_date
	clickTime := time.Now()
	if !msg.ClickTime.IsZero() {
		clickTime = msg.ClickTime
	}
	date := clickTime.Format("2006-01-02")
	key := msg.ShortCode + "_" + date

	h.mu.Lock()
	h.buffer[key]++
	currentSize := len(h.buffer)
	h.mu.Unlock()

	// 如果达到批量大小，立即刷新
	if currentSize >= h.bufferSize {
		h.triggerFlush()
	}
	return true
}

func (h *SummaryHandler) triggerFlush() {
	h.mu.Lock()
	if len(h.buffer) == 0 {
		h.bufferSize = 0
		return
	}

	// 复制并清空缓冲区
	bufferCopy := h.buffer
	h.buffer = make(map[string]int)
	h.mu.Unlock()

	logger.Logger.Info("Starting flushing buffer", zap.Int("buffer_size", len(bufferCopy)))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := h.summaryService.RecordSummary(ctx, h.handlerKey, bufferCopy); err != nil {
		logger.Logger.Error("failed to flush buffer",
			zap.Int("record_count", len(bufferCopy)), zap.Error(err))
	}
}
