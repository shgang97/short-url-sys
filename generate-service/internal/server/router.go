package server

import (
	"generate-service/internal/config"
	"generate-service/internal/handler"
	"generate-service/internal/model"
	"generate-service/internal/server/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(config *config.Config, srv *Server) {
	// 设置Gin模式
	gin.SetMode(config.Server.HTTP.Mode)
	router := gin.New()

	// 全局中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(middleware.ErrorHandler())

	// 初始化处理器
	linkHandler := handler.NewLinkHandler(srv.linkSvc, config.Server.BaseURL)
	qrcodeHandler := handler.NewQRCodeHandler(srv.linkSvc, config.Server.BaseURL)

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		now := time.Now()
		health := model.HealthResponse{
			Status:    "OK",
			Datetime:  now.Format("2006-01-02 15:04:05.000"),
			Timestamp: now.Unix(),
			Services:  make(map[string]string),
		}

		// 检查MySQL连接
		if srv.mysqlDB != nil {
			if err := srv.mysqlDB.HealthCheck(); err != nil {
				health.Status = "degraded"
				health.Services["mysql"] = "unhealthy"
			} else {
				health.Services["mysql"] = "healthy"
			}
		}

		// 检查Redis连接
		if srv.redisClient != nil {
			if err := srv.redisClient.HealthCheck(); err != nil {
				health.Status = "degraded"
				health.Services["redis"] = "unhealthy"
			} else {
				health.Services["redis"] = "healthy"
			}
		}

		c.JSON(200, health)
	})

	api := router.Group("/api/v1")
	api.Use(middleware.RateLimit(10)) // 每分钟10个请求
	{
		// 短链相关接口
		linkGroup := api.Group("/links")
		linkGroup.POST("/short", linkHandler.CreateShortURL)
		linkGroup.GET("/:code", linkHandler.GetLinkInfo)
		linkGroup.PUT("/:code", linkHandler.UpdateLink)
		linkGroup.DELETE("/:code", linkHandler.DeleteLink)
		linkGroup.GET("", linkHandler.ListLinks)
		linkGroup.POST("/short/batch", linkHandler.BatchCreate)

		// 生成二维码相关接口
		qrcodeGroup := linkGroup.Group("/qrcode")
		{
			qrcodeGroup.GET("/:code", qrcodeHandler.GenerateQRCode)
		}
	}

	api.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "generate-service",
			"version": "1.0.0",
		})
	})

	srv.router = router
}
