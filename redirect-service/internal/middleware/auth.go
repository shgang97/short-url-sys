package middleware

import (
	"redirect-service/internal/pkg/ipgen"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

// AuthMiddleware 模拟的认证中间件，设置用户信息
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(UserIDKey, "1")
		c.Set(UsernameKey, "system")
		c.Request.Header.Set("X-Forwarded-For", ipgen.IpGenerator.GeneratePublicIP())
		c.Next()
	}
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx *gin.Context) (userID, username string) {
	if id, exists := ctx.Get(UserIDKey); exists {
		userID = id.(string)
	}
	if name, exists := ctx.Get(UsernameKey); exists {
		username = name.(string)
	}
	return
}
