package consumer

import (
	"context"
	"errors"
	"fmt"
	"statistics-service/internal/pkg/logger"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	config   *KafkaConfig
	group    *GroupConfig
	consumer sarama.ConsumerGroup
	router   *HandlerRouter

	// 独立控制字段
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

type KafkaConfig struct {
	Brokers  []string       `mapstructure:"brokers"`
	ClientID string         `mapstructure:"client_id"`
	Version  string         `mapstructure:"version"`
	interval time.Duration  `mapstructure:"interval"` // 错误恢复时间隔时间
	Groups   []*GroupConfig `mapstructure:"groups"`
}

type GroupConfig struct {
	Id                 string        `mapstructure:"id"`
	Topics             []string      `mapstructure:"topics"`
	FetchMaxBytes      int32         `mapstructure:"fetch_max_bytes"`
	AutoCommit         bool          `mapstructure:"auto_commit"`
	AutoCommitInterval int64         `mapstructure:"auto_commit_interval"`
	AutoOffset         string        `mapstructure:"auto_offset"`
	SessionTimeout     time.Duration `mapstructure:"session_timeout"`
	BatchSize          int           `mapstructure:"batch_size"` // 批次消费控制缓冲区大小
	Spec               string        `mapstructure:"spec"`       // 批次消费控制缓冲区大小
}

func NewKafkaConsumer(cfg *KafkaConfig, groupCfg *GroupConfig, router *HandlerRouter) (*KafkaConsumer, error) {
	// 设置kafka客户端配置
	saramaCfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, fmt.Errorf("error parsing Kafka version: %w", err)
	}
	saramaCfg.Version = version

	// 消费者配置
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Fetch.Max = groupCfg.FetchMaxBytes
	saramaCfg.Consumer.Group.Session.Timeout = groupCfg.SessionTimeout
	saramaCfg.Consumer.Group.Heartbeat.Interval = groupCfg.SessionTimeout / 3 // 心跳间隔通常为会话超时的1/3

	// 自动提交配置
	saramaCfg.Consumer.Offsets.AutoCommit.Enable = groupCfg.AutoCommit
	saramaCfg.Consumer.Offsets.Initial = groupCfg.AutoCommitInterval
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最新位置开始消费

	// 创建消费者组
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, groupCfg.Id, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		config:   cfg,
		group:    groupCfg,
		consumer: consumer,
		router:   router,
		ctx:      ctx,
		cancel:   cancel,
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
	handler := &KafkaHandler{
		router:  k.router,
		groupId: k.group.Id,
	}
	for {
		select {
		case <-k.ctx.Done():
			logger.Logger.Info("Consumer shutting down")
			return
		default:
			// 消费指定主题
			if err := k.consumer.Consume(k.ctx, k.group.Topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					logger.Logger.Info("Consumer group closed")
				}
				logger.Logger.Error("Failed to consume message", zap.Error(err))
				logger.Logger.Info("Retrying...", zap.String("groupId", k.group.Id))
				time.Sleep(k.config.interval * time.Second)
			}
		}
	}
}

func (k *KafkaConsumer) Close() error {
	k.cancel()
	return k.consumer.Close()
}
