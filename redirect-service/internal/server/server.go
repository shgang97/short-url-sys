package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"redirect-service/internal/client/grpc/generate"
	"redirect-service/internal/client/redis"
	"redirect-service/internal/config"
	"redirect-service/internal/consumer"
	"redirect-service/internal/pkg/idgen"
	"redirect-service/internal/producer"
	"redirect-service/internal/repository/cache"
	cacheService "redirect-service/internal/service/cache"
	"redirect-service/internal/service/geoip"
	redirectService "redirect-service/internal/service/redirect"
	"syscall"
	"time"
)

type Server struct {
	config        *config.Config
	router        http.Handler
	server        *http.Server
	redisClient   *redis.Client
	cacheRepo     *cache.Repository
	redirectSvc   *redirectService.Service
	CacheSvc      *cacheService.Service
	genClient     *generate.Client
	kafkaConsumer *consumer.KafkaConsumer
	kafkaProducer *producer.KafkaProducer
	geoIPSvc      geoip.Service
	generator     idgen.Generator
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	// 初始化Redis
	redisClient, err := redis.NewRedis(&s.config.Redis, &s.config.Breaker)
	if err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	s.redisClient = redisClient

	// 初始化 IDGenerator
	generator := idgen.NewSfGenerator(&s.config.Generator.Sonyflake)
	s.generator = generator

	// 初始化Repository
	s.cacheRepo = cache.NewRepository(redisClient, &s.config.Cache)

	// 初始化 GeoIPService
	geoIPSvc, err := geoip.New(s.config.GeoIP.DBPath)
	if err != nil {
		return fmt.Errorf("init geoip failed: %w", err)
	}
	s.geoIPSvc = geoIPSvc

	// 初始化 generate-service 客户端
	genClient, err := generate.NewClient(&s.config.GenerateService, &s.config.Etcd)
	if err != nil {
		return fmt.Errorf("init generate client failed: %w", err)
	}
	s.genClient = genClient

	// 初始化缓存服务
	s.CacheSvc = cacheService.NewService(s.cacheRepo)

	// 初始化 MessageHandler 并注册到 HandlerRouter
	router := consumer.NewHandlerRouter()
	for _, topic := range s.config.Kafka.Topics {
		h, err := consumer.CreateHandler(topic, s.CacheSvc)
		if err != nil {
			log.Printf("create kafka handler failed: %v", err)
		} else {
			router.Register(topic, h)
		}
	}

	// 初始化KafkaConsumer
	kafkaConsumer, err := consumer.NewKafkaConsumer(&s.config.Kafka, s.CacheSvc, router)
	if err != nil {
		return fmt.Errorf("init kafka consumer failed: %w", err)
	}
	s.kafkaConsumer = kafkaConsumer
	go kafkaConsumer.Start()

	// 初始化 KafkaProducer 并启动
	kafkaProducer, err := producer.NewKafkaProducer(s.config.Kafka.Brokers, &s.config.Kafka)
	if err != nil {
		return fmt.Errorf("init kafka producer failed: %w", err)
	}
	s.kafkaProducer = kafkaProducer
	go kafkaProducer.Start()

	// 初始化重定向服务
	s.redirectSvc = redirectService.NewService(
		s.genClient,
		s.cacheRepo,
		s.kafkaProducer,
		s.geoIPSvc,
		s.generator,
	)

	// 设置路由
	setupRouter(s.config, s)

	// 设置服务器
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("server listening on %s", addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Printf("server started")
	}()

	// 等待中断信号
	s.waitForShutdown()

	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭Redis连接
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	// 关闭Kafka
	if s.kafkaConsumer != nil {
		s.redisClient.Close()
	}
	// 关闭KafkaProducer
	if s.kafkaProducer != nil {
		s.kafkaProducer.Close()
	}
	log.Printf("Server exiting...\n")
}
