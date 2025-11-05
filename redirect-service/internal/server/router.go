package server

import (
	"redirect-service/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(config *config.Config, srv *Server) {
	// 设置Gin模式
	gin.SetMode(config.Server.Mode)
	router := gin.New()

	// 全局中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "redirect",
			"timestamp": time.Now().Unix(),
		})
	})

	api := router.Group("/api/v1")

	api.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "redirect-service",
			"version": "1.0.0",
		})
	})

	srv.router = router
}
