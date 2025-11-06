package consumer

import (
	"context"
	"errors"
	"fmt"
	"statistics-service/internal/pkg/logger"
	"statistics-service/internal/service/click"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	config       *KafkaConfig
	consumer     sarama.ConsumerGroup
	clickService *click.Service
	router       *HandlerRouter
	ctx          context.Context
	cancel       context.CancelFunc
}

type KafkaConfig struct {
	Brokers       []string `mapstructure:"brokers"`
	ClientID      string   `mapstructure:"client_id"`
	Version       string   `mapstructure:"version"`
	GroupID       string   `mapstructure:"group_id"`
	Topics        []string `mapstructure:"topics"`
	FetchMaxBytes int32    `mapstructure:"fetch_max_bytes"`
	Consumer      Consumer `mapstructure:"consumer"`
}

type Consumer struct {
	AutoCommit         bool          `mapstructure:"auto_commit"`
	AutoCommitInterval int64         `mapstructure:"auto_commit_interval"`
	AutoOffset         string        `mapstructure:"auto_offset"`
	SessionTimeout     time.Duration `mapstructure:"session_timeout"`
}

func NewKafkaConsumer(cfg *KafkaConfig, clickService *click.Service, router *HandlerRouter) (*KafkaConsumer, error) {
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
		clickService: clickService,
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
			logger.Logger.Error("Error from consumer", zap.Error(err))
		}
	}()

	// 循环消费
	handler := &KafkaHandler{router: k.router}
	for {
		select {
		case <-k.ctx.Done():
			logger.Logger.Info("Consumer shutting down")
			return
		default:
			// 消费指定主题
			if err := k.consumer.Consume(k.ctx, k.config.Topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					logger.Logger.Info("Consumer group closed")
				}
				logger.Logger.Info("Failed to consume message", zap.Error(err))
			}
		}
	}
}

func (k *KafkaConsumer) Close() error {
	k.cancel()
	return k.consumer.Close()
}
