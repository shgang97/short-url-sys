package producer

import (
	"log"
	"redirect-service/internal/config"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
	config   *sarama.Config
	brokers  []string
	mu       sync.Mutex
	topics   []string
}

func NewKafkaProducer(brokers []string, config *config.KafkaConfig) (*KafkaProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(config.Producer.RequiredAcks)
	cfg.Producer.Compression = sarama.CompressionCodec(config.Producer.Compression)
	cfg.Producer.Flush.Frequency = config.Producer.Flush.Frequency
	cfg.Producer.Return.Successes = config.Producer.Return.Successes
	cfg.Producer.Return.Errors = config.Producer.Return.Errors
	cfg.Producer.Retry.Max = config.Producer.Retry.Max
	cfg.Producer.Idempotent = config.Producer.Idempotent // 启用幂等性
	cfg.Net.MaxOpenRequests = config.NetMaxOpenRequests

	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Printf("brokers: %v", brokers)
		log.Printf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	kp := &KafkaProducer{
		producer: producer,
		config:   cfg,
		brokers:  brokers,
		topics:   make([]string, 0),
	}

	return kp, nil
}

func (kp *KafkaProducer) Start() {
	// 启动 KafkaProducer 处理成功和失败的消息
	for {
		select {
		case success := <-kp.producer.Successes():
			// 真正发送成功，打印详细消息
			log.Printf("Message successfully sent to topic %s: partition=%d, offset=%d, key=%s",
				success.Topic, success.Partition, success.Offset, success.Key)
		case err := <-kp.producer.Errors():
			// TODO 这里可以添加重试逻辑或告警
			log.Printf("Message failed to send: topic=%s, key=%s, error=%v",
				err.Msg.Topic, err.Msg.Key, err.Err)
		}
	}
}

func (kp *KafkaProducer) Close() error {
	if kp.producer != nil {
		return kp.producer.Close()
	}
	return nil
}
