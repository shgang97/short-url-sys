package handler

import (
	"log"
	"net/http"
	"redirect-service/internal/middleware"
	"redirect-service/internal/service/redirect"
	"shared/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	redirectService redirect.Service
}

func NewRedirectHandler(redirectService redirect.Service) *RedirectHandler {
	return &RedirectHandler{
		redirectService: redirectService,
	}
}

// Redirect 重定向
// @Router /{code} [get]
func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "Short code is required",
		})
		return
	}
	// 获取原始URL
	originalUrl, err := h.redirectService.GetOriginalUrl(c, shortCode)
	if err != nil {
		c.Error(err)
		return
	}

	// 异步记录点击事件
	ip := getClientIP(c)
	_, username := middleware.GetUserFromContext(c)
	req := &redirect.RedirectRequest{
		OriginalURL: originalUrl,
		IPAddress:   ip,
		UserAgent:   c.Request.UserAgent(),
		Referer:     c.Request.Referer(),
		Username:    username,
	}
	go func() {
		err = h.redirectService.RecordClick(c.Request.Context(), shortCode, req)
		if err != nil {
			log.Printf("failed to record click on short code: %v", err)
		}
	}()

	// 302 重定向
	c.Redirect(http.StatusFound, originalUrl)
}

// 获取客户端IP
func getClientIP(c *gin.Context) string {
	// 尝试从 X-Forwarded-For 获取
	// X-Forwarded-For：包含客户端和所有代理服务器的 IP 链
	if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
		if ips := strings.Split(forwarded, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从 X-Real-IP 获取
	// X-Real-IP：通常由第一个代理设置，直接包含客户端真实 IP
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 使用远程地址
	return c.ClientIP()
}
