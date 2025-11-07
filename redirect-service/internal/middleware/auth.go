package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

// AuthMiddleware 模拟的认证中间件，设置用户信息
func AuthMiddleware() gin.HandlerFunc {
	log.Printf("设置上下文信息")
	return func(c *gin.Context) {
		c.Set(UserIDKey, "1")
		c.Set(UsernameKey, "system")
		c.Next()
	}
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx *gin.Context) (userID, username string) {
	log.Printf("获取上下文信息")
	if id, exists := ctx.Get(UserIDKey); exists {
		userID = id.(string)
	}
	if name, exists := ctx.Get(UsernameKey); exists {
		username = name.(string)
	}
	return
}
