package consumer

import (
	"statistics-service/internal/pkg/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaHandler struct {
	router  *HandlerRouter
	groupId string
}

func (k *KafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	logger.Logger.Info("Setting up session for consumer")
	return nil
}

func (k *KafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	logger.Logger.Info("Cleaning up session for consumer")
	return nil
}

// ConsumeClaim 处理分区消息
func (k *KafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 循环读取消息
	for msg := range claim.Messages() {
		topic := msg.Topic
		handler, exists := k.router.GetHandler(k.groupId, topic)
		if !exists {
			logger.Logger.Warn("Skipping message due to missing handler for ",
				zap.String("groupId", k.groupId), zap.String("topic", topic))
			continue
		}
		success := handler.Handle(topic, msg.Value)
		if success {
			// 消费成功，提交offset
			session.MarkMessage(msg, "")
		} else {
			// 消费失败
			logger.Logger.Error("Failed to consume message for ", zap.String("topic", topic))
		}
	}

	return nil
}
