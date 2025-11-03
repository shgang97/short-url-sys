package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin", "Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With",
		},
		ExposedHeaders: []string{
			"Content-Length", "Link",
		},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour),
	})
}
