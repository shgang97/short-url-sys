package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"statistics-service/internal/config"
	"statistics-service/internal/consumer"
	"statistics-service/internal/pkg/database"
	"statistics-service/internal/pkg/idgen"
	"statistics-service/internal/pkg/logger"
	clickRepo "statistics-service/internal/repository/click"
	clickService "statistics-service/internal/service/click"
	detector "statistics-service/internal/service/device_detector"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	config    *config.Config
	server    *http.Server
	router    http.Handler
	mysqlDB   *database.MySQLDB
	clickRepo clickRepo.Repository
	clickSvc  *clickService.Service
	generator idgen.Generator
	detector  detector.DeviceDetector
}

func New(cfg *config.Config) *Server {
	return &Server{config: cfg}
}

func (s *Server) Start() error {
	// 初始化日志
	if err := logger.InitLogger(&s.config.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化 IDGenerator
	generator := idgen.NewSfGenerator(&s.config.Generator.Sonyflake)
	s.generator = generator

	// 初始化设备检测器 DeviceDetector
	s.detector = detector.NewDefaultDeviceDetector()

	// 初始化 MySQL
	sqldb, err := database.NewMySQLDB(s.config.MySQL)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL DB: %v", err)
	}
	s.mysqlDB = sqldb

	// 初始化 Repository
	s.initRepository()

	// 初始化 Service
	s.initService()

	// 初始化 MessageHandler 并注册到 HandlerRouter
	router := consumer.NewHandlerRouter()
	for _, topic := range s.config.Kafka.Topics {
		h, err := consumer.CreatHandler(topic, s.clickSvc)
		if err != nil {
			logger.Logger.Error("Failed to create MessageHandler", zap.String("topic", topic), zap.Error(err))
		} else {
			router.Register(topic, h)
		}
	}

	// 初始化KafkaConsumer
	kafkaConsumer, err := consumer.NewKafkaConsumer(&s.config.Kafka, s.clickSvc, router)
	if err != nil {
		return fmt.Errorf("init kafka consumer failed: %v", err)
	}
	go kafkaConsumer.Start()

	// 设置路由
	setupRouter(s.config, s)

	// 初始化 HTTP 服务器
	httpAddr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr:         httpAddr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动 HTTP 服务器
	go func() {
		logger.Logger.Info("statistics-service server listening on address: " + s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Info("Failed to start statistics-service server", zap.Error(err))
		}
	}()

	// 等待中断信号
	s.waitForShutdown()
	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down statistics-service server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Failed to shut down statistics-service server", zap.Error(err))
	}

	logger.Logger.Info("Shut down statistics-service server successfully")
}

func (s *Server) initRepository() {
	s.clickRepo = clickRepo.NewMySQLRepository(s.mysqlDB.DB)
}

func (s *Server) initService() {
	s.clickSvc = clickService.NewService(
		s.mysqlDB.DB,
		s.clickRepo,
		s.generator,
		s.detector,
	)
}
