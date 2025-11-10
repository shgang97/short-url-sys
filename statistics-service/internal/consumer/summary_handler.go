package consumer

import (
	"context"
	"encoding/json"
	"shared/message"
	trac "statistics-service/internal/consumer/tracker"
	"statistics-service/internal/pkg/logger"
	"statistics-service/internal/service/summary"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// ClickMessageWrapper 包装消息和会话信息
type ClickMessageWrapper struct {
	Message  *sarama.ConsumerMessage
	Session  sarama.ConsumerGroupSession
	EventID  string
	ClickMsg *message.ClickEventMessage
	Count    int32
}

// SummaryHandler 处理点击量的handler
type SummaryHandler struct {
	handlerKey     string
	summaryService *summary.Service
	tracker        trac.ProcessedEventTracker
	buffer         map[string]int
	messageBuffer  []*ClickMessageWrapper // 保存消息用于手动提交
	batchSize      int
	wg             *sync.WaitGroup
	mu             sync.Mutex
	cron           *cron.Cron
	messages       []*sarama.ConsumerMessage
	dateFormat     string
	spec           string
}

func NewSummaryHandler(
	handlerKey string,
	summaryService *summary.Service,
	batchSize int,
	spec string,
) *SummaryHandler {
	h := &SummaryHandler{
		handlerKey:     handlerKey,
		summaryService: summaryService,
		tracker:        trac.NewDefaultTracker(),
		buffer:         make(map[string]int),
		messageBuffer:  make([]*ClickMessageWrapper, 0, batchSize),
		batchSize:      batchSize,
		cron:           cron.New(cron.WithSeconds()),
		messages:       make([]*sarama.ConsumerMessage, 0),
		dateFormat:     "2006-01-02",
		wg:             &sync.WaitGroup{},
		spec:           spec,
	}
	h.scheduleFlush()
	return h
}

// Handle 处理单挑点击消息
func (h *SummaryHandler) Handle(topic string, msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	var baseMsg message.BaseMessage
	if err := json.Unmarshal(msg.Value, &baseMsg); err != nil {
		logger.Logger.Error("failed to unmarshal click message",
			zap.String("topic", topic), zap.Error(err))
		// 标记为已处理，避免重复处理无法解析的消息
		session.MarkMessage(msg, "")
		return
	}

	// 检查事件是否已处理（幂等性检查）
	processed, err := h.tracker.IsProcessed(baseMsg.EventID)
	if err != nil {
		logger.Logger.Error("failed to check event processed status",
			zap.String("eventID", baseMsg.EventID), zap.Error(err))
		// 检查失败，不处理小修，等待重试
		return
	}

	if processed {
		// 标记为已处理，避免重复消费
		session.MarkMessage(msg, "")
		return
	}
	logger.Logger.Info("start handle message",
		zap.String("topic", topic), zap.ByteString("value", msg.Value))
	var clickMsg message.ClickEventMessage
	if err := json.Unmarshal(msg.Value, &clickMsg); err != nil {
		logger.Logger.Error("failed to unmarshal click message",
			zap.String("topic", topic), zap.Error(err))
		// 标记为已处理，避免重复处理无法解析的消息
		session.MarkMessage(msg, "")
		return
	}

	// 生成聚合key: shortCode_date
	if clickMsg.ClickTime.IsZero() {
		clickMsg.ClickTime = time.Now()
	}
	key := h.getKey(clickMsg.ShortCode, clickMsg.ClickTime)

	h.mu.Lock()
	h.buffer[key]++
	h.messageBuffer = append(h.messageBuffer, &ClickMessageWrapper{
		Message:  msg,
		Session:  session,
		EventID:  clickMsg.EventID,
		ClickMsg: &clickMsg,
		Count:    0,
	})
	currentSize := len(h.messageBuffer)
	h.mu.Unlock()

	// 如果达到批量大小，立即刷新
	if currentSize >= h.batchSize {
		h.triggerFlush()
	}
}

func (h *SummaryHandler) triggerFlush() {
	h.mu.Lock()
	if len(h.buffer) == 0 || len(h.messageBuffer) == 0 {
		h.mu.Unlock()
		return
	}
	h.wg.Add(1)
	defer h.wg.Done()

	// 复制并清空缓冲区
	bufferCopy := h.buffer
	h.buffer = make(map[string]int)
	messageBufferCopy := h.messageBuffer
	h.messageBuffer = make([]*ClickMessageWrapper, 0, h.batchSize)
	h.mu.Unlock()

	logger.Logger.Info("Starting flushing buffer",
		zap.Int("message_count", len(messageBufferCopy)),
		zap.Int("batch_size", len(bufferCopy)))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := h.summaryService.RecordSummary(ctx, h.handlerKey, bufferCopy); err != nil {
		logger.Logger.Error("failed to flush buffer",
			zap.Int("record_count", len(bufferCopy)), zap.Error(err))
		// 如果处理失败，将数据重新放回缓冲区
		h.restoreToBuffer(messageBufferCopy)
		return
	}
	// 处理成功，标记所有消息为已处理
	logger.Logger.Info("Successfully flushing buffer", zap.Int("message_count", len(messageBufferCopy)))
	h.markMessageAsProcessed(messageBufferCopy)
}

func (h *SummaryHandler) getKey(shortCode string, clickTime time.Time) string {
	return shortCode + "_" + clickTime.Format(h.dateFormat)
}

func (h *SummaryHandler) restoreToBuffer(msgWrapperBuffer []*ClickMessageWrapper) {
	h.mu.Lock()

	// 将失败的数据重新放回缓冲区
	for _, wrapper := range msgWrapperBuffer {
		if wrapper.Count >= 3 { // 限制最多重试3次
			logger.Logger.Error("message dropped after max retries, committing offset",
				zap.String("eventID", wrapper.EventID), zap.Int("retry_count", int(wrapper.Count)))
			continue
		}
		wrapper.Count++
		key := h.getKey(wrapper.ClickMsg.ShortCode, wrapper.ClickMsg.ClickTime)
		h.buffer[key]++
		h.messageBuffer = append(h.messageBuffer, wrapper)
	}
	h.mu.Unlock()
}

func (h *SummaryHandler) markMessageAsProcessed(messages []*ClickMessageWrapper) {
	// 先标记事件为已处理，再提交Kafka偏移量
	for _, wrapper := range messages {
		if err := h.tracker.MarkProcessed(wrapper.EventID); err != nil {
			// 标记失败，记录错误，依然提交偏移量（因为数据已经被正确处理了）
			logger.Logger.Error("failed to mark message as processed",
				zap.String("eventID", wrapper.EventID), zap.Error(err))
			continue
		}
		wrapper.Session.MarkMessage(wrapper.Message, "")
	}

}

func (h *SummaryHandler) scheduleFlush() {
	_, err := h.cron.AddFunc(h.spec, func() {
		h.mu.Lock()
		if len(h.buffer) == 0 {
			logger.Logger.Info("Buffer is empty, skip flushing")
			h.mu.Unlock()
			return
		}
		h.mu.Unlock()
		h.triggerFlush()
	})
	if err != nil {
		logger.Logger.Error("failed to schedule flush", zap.Error(err))
		return
	}
	h.cron.Start()
	logger.Logger.Info("Successfully start scheduled flush")
}

func (h *SummaryHandler) Close() {
	if h.cron != nil {
		h.cron.Stop()
	}
	// 关闭前处理剩余缓冲区
	h.triggerFlush()

	// 等待处理完成，但设置超时
	done := make(chan struct{})
	go func() {
		h.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logger.Info("SummaryHandler closed gracefully")
	case <-time.After(30 * time.Second):
		logger.Logger.Warn("SummaryHandler close timeout, forcing exit")
	}
}
