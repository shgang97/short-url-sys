package server

import (
	"context"
	"errors"
	"fmt"
	"generate-service/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
	router http.Handler
	server *http.Server
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
	// TODO
	return nil
}

func (s *Server) initServices() error {
	// TODO
	return nil
}
