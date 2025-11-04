package main

import (
	"generate-service/internal/config"
	"generate-service/internal/server"
	"log"
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
		log.Fatalf("Failed to start api-server: %v", err)
	}
}
