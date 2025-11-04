package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Message interface {
	GetEventType() string
	GetKey() string
	Validate() error
}

type KafkaProducer struct {
	producer sarama.AsyncProducer
	config   *sarama.Config
	brokers  []string
	mu       sync.Mutex
	topics   map[string]bool
}

func NewKafkaProducer(brokers []string, config *sarama.Config) (*KafkaProducer, error) {
	if config == nil {
		config = getDefaultConfig()
	}

	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Printf("brokers: %v", brokers)
		log.Printf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	kp := &KafkaProducer{
		producer: producer,
		config:   config,
		brokers:  brokers,
		topics:   make(map[string]bool),
	}

	// 启动goroutine处理成功和失败的消息
	go func() {
		for {
			select {
			case success := <-producer.Successes():
				// 真正的发送成功，打印详细信息
				log.Printf("Message successfully sent to topic %s: partition=%d, offset=%d, key=%s",
					success.Topic, success.Partition, success.Offset, success.Key)
			case err := <-producer.Errors():
				// TODO 这里可以添加重试逻辑或告警
				log.Printf("Message failed to send: topic=%s, key=%s, error=%v",
					err.Msg.Topic, err.Msg.Key, err.Err)
			}
		}
	}()

	return kp, nil
}

// 默认配置
func getDefaultConfig() *sarama.Config {
	defaultCfg := sarama.NewConfig()
	defaultCfg.Producer.RequiredAcks = sarama.WaitForAll
	defaultCfg.Producer.Compression = sarama.CompressionSnappy
	defaultCfg.Producer.Flush.Frequency = 500 * time.Millisecond
	defaultCfg.Producer.Return.Successes = true
	defaultCfg.Producer.Return.Errors = true
	defaultCfg.Producer.Retry.Max = 3
	defaultCfg.Producer.Idempotent = true // 启用幂等性
	defaultCfg.Net.MaxOpenRequests = 1
	return defaultCfg
}

// SendMessage 发送消息
func (kp *KafkaProducer) SendMessage(topic string, message Message) error {
	if err := message.Validate(); err != nil {
		return fmt.Errorf("message validation failed: %v", err)
	}

	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("message marshalling failed: %v", err)
	}

	// 构造kafka消息
	kafkaMsg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(message.GetKey()),
		Value:     sarama.ByteEncoder(messageBytes),
		Timestamp: time.Now(),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(message.GetEventType()),
			},
		},
	}

	// 发送消息写入Input通道
	select {
	case kp.producer.Input() <- kafkaMsg:
		// 仅表示消息已进入发送队列
		log.Printf("Message added to send queue: topic=%s, key=%s", topic, message.GetKey())
	case <-time.After(5 * time.Second):
		// 防止Input通道阻塞超时
		return fmt.Errorf("timeout sending message to input channel")
	}

	return nil
}

// SendMessageWithRetry 带重试的消息发送
func (kp *KafkaProducer) SendMessageWithRetry(topic string, message Message, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := kp.SendMessage(topic, message)
		if err != nil {
			lastErr = err
			log.Printf("Attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}
		return nil
	}
	return fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}

func (kp *KafkaProducer) Close() error {
	if kp.producer != nil {
		return kp.producer.Close()
	}
	return nil
}
