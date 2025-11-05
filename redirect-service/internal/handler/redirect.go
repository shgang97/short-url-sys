package handler

import (
	"net/http"
	"redirect-service/internal/service/redirect"
	"shared/model"

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
	c.Redirect(http.StatusFound, originalUrl)

	// TODO 异步记录点击统计
}
