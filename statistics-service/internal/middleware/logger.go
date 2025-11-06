package middleware

import (
	"net/http"
	"statistics-service/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		//query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		cost := time.Since(start)
		logger.Logger.Info(
			path,
			zap.String("request_id", c.Request.Header.Get("X-Request-Id")),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			//zap.String("path", path),
			//zap.String("query", query),
			//zap.String("ip", c.ClientIP()),
			//zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 恢复中间件
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", c.Request.Header.Get("X-Request-Id")),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("ip", c.ClientIP()),
					zap.Stack("stack"),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
