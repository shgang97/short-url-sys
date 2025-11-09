package consumer

import (
	"statistics-service/internal/pkg/logger"
	"sync"
)

type KafkaConsumerManager struct {
	consumers []*KafkaConsumer
	wg        sync.WaitGroup
}

func NewConsumerManager() *KafkaConsumerManager {
	return &KafkaConsumerManager{
		consumers: make([]*KafkaConsumer, 0),
	}
}

// AddConsumer 添加消费者
func (k *KafkaConsumerManager) AddConsumer(kafkaConsumer *KafkaConsumer) {
	k.consumers = append(k.consumers, kafkaConsumer)
}

func (k *KafkaConsumerManager) CloseAll() {
	logger.Logger.Info("Closing all kafka consumers...")
	for _, consumer := range k.consumers {
		consumer.Close()
	}
}
