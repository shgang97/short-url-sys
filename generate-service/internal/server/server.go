package server

import (
	"context"
	"errors"
	"fmt"
	"generate-service/internal/config"
	"generate-service/internal/pkg/database"
	linkRepo "generate-service/internal/repository/link"
	"generate-service/internal/service/idgen"
	linkService "generate-service/internal/service/link"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config      *config.Config
	router      http.Handler
	server      *http.Server
	mysqlDB     *database.MySQLDB
	redisClient *database.RedisClient
	linkRepo    linkRepo.Repository
	idGenerator idgen.Generator
	linkSvc     linkService.Service
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

	// 初始化服务
	if err := s.initServices(); err != nil {
		return fmt.Errorf("failed to init services: %w", err)
	}

	// 设置路由
	setupRouter(s.config, s)

	// 设置服务器配置
	serverCfg := s.config.Server
	addr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)
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

	// TODO 关闭数据库连接

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
	)

	log.Println("✅ Services initialized successfully")
	return nil
}
