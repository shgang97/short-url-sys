package server

import (
	"net/http"
	"shared/model"
	"statistics-service/internal/config"
	"statistics-service/internal/handler"
	"statistics-service/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(cfg *config.Config, srv *Server) {
	// 设置路由模式
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()

	// 设置全局中间件
	router.Use(middleware.GinLogger())
	router.Use(middleware.GinRecovery())

	// 初始化处理器
	statsHandler := handler.NewStatsHandler(srv.clickSvc)

	// 健康检查点
	router.GET("/health", func(c *gin.Context) {
		now := time.Now()
		healthRsp := model.HealthResponse{
			Status:    "OK",
			Datetime:  now.Format("2006-01-02 15:04:05.000"),
			Timestamp: now.Unix(),
			Services:  make(map[string]string),
		}
		// 检查MySQL连接
		if srv.mysqlDB != nil {
			if err := srv.mysqlDB.HealthCheck(); err != nil {
				healthRsp.Status = "degraded"
				healthRsp.Services["mysql"] = "unhealthy"
			} else {
				healthRsp.Services["mysql"] = "healthy"
			}
		}
		c.JSON(http.StatusOK, healthRsp)
	})

	api := router.Group("/api/v1")
	statsRouter := api.Group("/stats")
	{
		statsRouter.GET("/:code", statsHandler.GetStatsSummary)
	}

	api.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "statistics-service",
			"version": "1.0.0",
		})
	})

	srv.router = router
}
