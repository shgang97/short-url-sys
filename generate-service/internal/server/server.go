package server

import (
	"context"
	"errors"
	"fmt"
	"generate-service/internal/config"
	"generate-service/internal/pkg/database"
	"generate-service/internal/pkg/mq"
	linkRepo "generate-service/internal/repository/link"
	grpcSrv "generate-service/internal/server/grpc"
	"generate-service/internal/service/idgen"
	linkService "generate-service/internal/service/link"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	pb "shared/proto/generate"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config        *config.Config
	router        http.Handler
	httpServer    *http.Server
	grpcServer    *grpc.Server
	mysqlDB       *database.MySQLDB
	redisClient   *database.RedisClient
	linkRepo      linkRepo.Repository
	idGenerator   idgen.Generator
	linkSvc       linkService.Service
	kafkaProducer *mq.KafkaProducer
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	// 初始化数据库
	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	// 初始化mq
	if err := s.initMq(); err != nil {
		return fmt.Errorf("failed to init mq: %w", err)
	}

	// 初始化服务
	if err := s.initServices(); err != nil {
		return fmt.Errorf("failed to init services: %w", err)
	}

	// 设置路由
	setupRouter(s.config, s)

	// 创建gRPC服务器
	grpcConfig := s.config.Server.GRPC
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(grpcConfig.MaxConcurrentStreams),
		grpc.MaxRecvMsgSize(grpcConfig.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(grpcConfig.MaxSendMsgSize),
	}
	gRPCServer := grpc.NewServer(opts...)
	// 注册服务
	generateServer := grpcSrv.NewGenerateServer(s.linkSvc)
	pb.RegisterGenerateServiceServer(gRPCServer, generateServer)
	// 启用反射服务
	reflection.Register(gRPCServer)
	// 启动gRPC服务器
	grpcAddr := fmt.Sprintf("%s:%s", grpcConfig.Host, grpcConfig.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServer = gRPCServer

	// 设置服务器配置
	httpConfig := s.config.Server.HTTP
	addr := fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("HTTPServer listening on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start httpServer: %v", err)
		}
	}()

	go func() {
		log.Printf("GRPCServer listening on %s", grpcAddr)
		if err := gRPCServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
	log.Printf("Shutting down httpServer...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	s.grpcServer.GracefulStop()

	// 关闭数据库连接
	if s.mysqlDB != nil {
		s.mysqlDB.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}

	log.Printf("Server exiting...\n")
}

func (s *Server) initDatabase() error {
	// 初始化MySQL
	mysqlDB, err := database.NewMySQLDB(&s.config.Database.MySQL)
	if err != nil {
		return fmt.Errorf("init mysql failed: %w", err)
	}
	s.mysqlDB = mysqlDB

	// 初始化Redis
	redisClient, err := database.NewRedisClient(&s.config.Redis)
	if err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	s.redisClient = redisClient

	// 初始化Repository
	s.linkRepo = linkRepo.NewMySQLRepository(mysqlDB.DB)

	log.Printf("✅ init database success\n")
	return nil
}

func (s *Server) initServices() error {
	// 厨师话ID生成器
	idGenerator, err := idgen.NewIDGenerator(&s.config.IdGenerator, s.redisClient)
	if err != nil {
		return fmt.Errorf("init ID Generator failed: %w", err)
	}
	s.idGenerator = idGenerator

	// 初始化短链服务
	s.linkSvc = linkService.NewService(
		s.linkRepo,
		s.idGenerator,
		linkService.Config{
			BaseURL: s.config.Server.BaseURL,
		},
		s.kafkaProducer,
	)

	log.Println("✅ Services initialized successfully")
	return nil
}

func (s *Server) initMq() error {
	kafkaConfig := s.config.Kafka
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(kafkaConfig.Producer.RequiredAcks)
	cfg.Producer.Compression = sarama.CompressionCodec(kafkaConfig.Producer.Compression)
	cfg.Producer.Flush.Frequency = kafkaConfig.Producer.Flush.Frequency
	cfg.Producer.Return.Successes = kafkaConfig.Producer.Return.Successes
	cfg.Producer.Return.Errors = kafkaConfig.Producer.Return.Errors
	cfg.Producer.Retry.Max = kafkaConfig.Producer.Retry.Max
	cfg.Producer.Idempotent = kafkaConfig.Producer.Idempotent // 启用幂等性
	cfg.Net.MaxOpenRequests = kafkaConfig.Net.MaxOpenRequests
	// 初始化Kafka
	kp, err := mq.NewKafkaProducer(s.config.Kafka.Brokers, cfg)
	if err != nil {
		return fmt.Errorf("init mq failed: %w", err)
	}
	s.kafkaProducer = kp
	return nil
}
