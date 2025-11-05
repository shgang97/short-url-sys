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
	"redirect-service/internal/repository/cache"
	redirectService "redirect-service/internal/service/redirect"
	"syscall"
	"time"
)

type Server struct {
	config      *config.Config
	router      http.Handler
	server      *http.Server
	redisClient *redis.Client
	cacheRepo   *cache.Repository
	redirectSvc *redirectService.Service
	genClient   *generate.Client
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Start() error {
	// 初始化Redis
	redisClient, err := redis.NewRedis(&s.config.Redis)
	if err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	s.redisClient = redisClient

	// 初始化KafkaProducer

	// 初始化Repository
	s.cacheRepo = cache.NewRepository(redisClient.Client, &s.config.Cache)

	// 初始化 generate-service 客户端
	genClient, err := generate.NewClient(&s.config.GenerateService)
	if err != nil {
		return fmt.Errorf("init generate client failed: %w", err)
	}
	s.genClient = genClient

	// 初始化重定向服务
	s.redirectSvc = redirectService.NewService(s.genClient, s.cacheRepo)

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
	log.Printf("Server exiting...\n")
}
