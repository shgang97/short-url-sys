package handler

import (
	"net/http"
	"short-url-sys/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSystemStats 获取系统统计
// @Router /api/v1/links/stats/system [get]
func (h *StatsHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.statsService.GetSystemStats(c)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetTopLinks 获取热门链接
// @Router /api/v1/links/stats/top-links [get]
func (h *StatsHandler) GetTopLinks(c *gin.Context) {
	limit := 10
	days := 7
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}
	links, err := h.statsService.GetTopLinks(c.Request.Context(), days, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, links)
}

// GetClickTimeline 获取点击时间线
// @Router /api/v1/links/stats/timeline/{code} [get]
func (h *StatsHandler) GetClickTimeline(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "Short code is required",
		})
		return
	}

	hours := 24
	if hoursParam := c.Query("hours"); hoursParam != "" {
		if h, err := strconv.Atoi(hoursParam); err == nil && h > 0 {
			hours = h
		}
	}

	timeline, err := h.statsService.GetClickTimeline(c.Request.Context(), shortCode, hours)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, timeline)
}

// GetGeographicStats 获取地理统计
// @Router /api/v1/links/stats/geographic/{code} [get]
func (h *StatsHandler) GetGeographicStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "Short code is required",
		})
		return
	}
	stats, err := h.statsService.GetGeographicStats(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetPlatformStats 获取平台统计
// @Router /api/v1/links/{code}/stats/platform [get]
func (h *StatsHandler) GetPlatformStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "Short code is required",
		})
		return
	}
	stats, err := h.statsService.GetPlatformStats(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
