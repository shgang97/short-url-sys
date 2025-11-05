package main

import (
	"log"
	"redirect-service/internal/config"
	"redirect-service/internal/server"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建服务器
	srv := server.New(cfg)

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start redirect-service: %v", err)
	}
}
