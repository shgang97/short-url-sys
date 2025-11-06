package main

import (
	"log"
	"statistics-service/internal/config"
	"statistics-service/internal/server"
)

func main() {
	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建服务器
	srv := server.New(cfg)

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start statistics-service server: %v", err)
	}
}
