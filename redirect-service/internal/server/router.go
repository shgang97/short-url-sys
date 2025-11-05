package server

import (
	"net/http"
	"redirect-service/internal/config"
	"redirect-service/internal/handler"
	"shared/model"
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

	// 初始化处理器
	redirectHandler := handler.NewRedirectHandler(*srv.redirectSvc)

	// 健康检查点
	router.GET("/health", func(c *gin.Context) {
		now := time.Now()
		healthRsp := model.HealthResponse{
			Status:    "OK",
			Datetime:  now.Format("2006-01-02 15:04:05.000"),
			Timestamp: now.Unix(),
			Services:  make(map[string]string),
		}
		// 检查Redis连接
		if srv.redisClient != nil {
			if err := srv.redisClient.HealthCheck(); err != nil {
				healthRsp.Status = "degraded"
				healthRsp.Services["redis"] = "unhealthy"
			} else {
				healthRsp.Services["redis"] = "healthy"
			}
		}
		c.JSON(http.StatusOK, healthRsp)
	})

	// 重定向路由
	router.GET("/:code", redirectHandler.Redirect)

	api := router.Group("/api/v1")

	api.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "redirect-service",
			"version": "1.0.0",
		})
	})

	srv.router = router
}
