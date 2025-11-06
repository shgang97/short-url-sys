package consumer

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaHandler struct {
	router *HandlerRouter
}

func (k *KafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Setting up session for consumer")
	return nil
}

func (k *KafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("Cleaning up session for consumer")
	return nil
}

// 处理分区消息
func (k *KafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 循环读取消息
	for msg := range claim.Messages() {
		topic := msg.Topic
		handler, exists := k.router.GetHandler(topic)
		if !exists {
			log.Printf("Skipping message due to missing handler for topic: %s", topic)
			continue
		}
		success := handler.Handle(topic, msg.Value)
		if success {
			// 消费成功，提交offset
			session.MarkMessage(msg, "")
		} else {
			// 消费失败
			log.Printf("Failed to consume message for topic: %s", topic)
		}
	}

	return nil
}
