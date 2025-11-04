package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 一个简单的限流器
// TODO 1. ips 没有清理，会有内存泄漏问题； 2. 在分布式环境下无法准确限流
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (rl *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(rl.r, rl.b)
	rl.ips[ip] = limiter
	return limiter
}

func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exist := rl.ips[ip]
	if !exist {
		return rl.AddIP(ip)
	}
	return limiter
}

// RateLimit 限流中间件
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	rateLimiter := NewIPRateLimiter(rate.Every(time.Minute), requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rateLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
