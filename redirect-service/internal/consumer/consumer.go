package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"redirect-service/internal/config"
	"redirect-service/internal/service/cache"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	config       *config.KafkaConfig
	consumer     sarama.ConsumerGroup
	cacheService *cache.Service
	router       *HandlerRouter
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewKafkaConsumer(cfg *config.KafkaConfig, cacheService *cache.Service, router *HandlerRouter) (*KafkaConsumer, error) {
	// 设置kafka客户端配置
	saramaCfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, fmt.Errorf("error parsing Kafka version: %w", err)
	}
	saramaCfg.Version = version

	// 消费者配置
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Fetch.Max = cfg.FetchMaxBytes
	saramaCfg.Consumer.Group.Session.Timeout = cfg.Consumer.SessionTimeout
	saramaCfg.Consumer.Group.Heartbeat.Interval = cfg.Consumer.SessionTimeout / 3 // 心跳间隔通常为会话超时的1/3

	// 自动提交配置
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = cfg.Consumer.AutoCommit
	saramaCfg.Consumer.Offsets.Initial = cfg.Consumer.AutoCommitInterval
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最新位置开始消费

	// 创建消费者组
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		config:       cfg,
		consumer:     consumer,
		cacheService: cacheService,
		router:       router,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// Start 启动消费者
func (k *KafkaConsumer) Start() {
	// 启动错误监听 goroutine
	go func() {
		for err := range k.consumer.Errors() {
			log.Printf("Error from consumer: %v", err)
		}
	}()

	// 循环消费
	handler := &KafkaHandler{
		router: k.router,
	}
	for {
		select {
		case <-k.ctx.Done():
			log.Printf("Consumer shutting down")
			return
		default:
			// 消费指定主题
			if err := k.consumer.Consume(k.ctx, k.config.Topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Printf("Consumer group closed")
				}
				log.Printf("Failed to consume message: %v", err)
			}
		}
	}
}

func (k *KafkaConsumer) Close() error {
	k.cancel()
	return k.consumer.Close()
}
