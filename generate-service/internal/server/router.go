package server

import (
	"generate-service/internal/config"
	"generate-service/internal/model"
	"generate-service/internal/server/middleware"
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
	router.Use(middleware.Cors())
	router.Use(middleware.ErrorHandler())

	// 初始化处理器
	// TODO

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		now := time.Now()
		health := model.HealthResponse{
			Status:    "OK",
			Datetime:  now.Format("2006-01-02 15:04:05.000"),
			Timestamp: now.Unix(),
			Services:  make(map[string]string),
		}

		// TODO 检查MYSQL连接
		// TODO 检查Redis连接

		c.JSON(200, health)
	})

	srv.router = router
}
