package handler

import (
	"net/http"
	"shared/model"
	"statistics-service/internal/service/click"
	"time"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	clickService *click.Service
}

func NewStatsHandler(clickService *click.Service) *StatsHandler {
	return &StatsHandler{clickService: clickService}
}

func (h *StatsHandler) GetStatsSummary(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "short code is required",
		})
		return
	}
	var req struct {
		StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
		EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid request",
			Message: "invalid query parameters",
		})
		return
	}

	resp, err := h.clickService.GetStatsSummary(c.Request.Context(), shortCode, req.StartDate, req.EndDate)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *StatsHandler) GetTimeSeries(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "short code is required",
		})
		return
	}
	unit := TimeUnit(c.Param("unit"))
	if unit == "" {
		unit = Daily
	}
	if !unit.IsValid() {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid unit",
			Message: "unit is must be one of: hourly, daily, weekly, monthly",
		})
		return
	}

	var req struct {
		StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
		EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid request",
			Message: "invalid query parameters",
		})
	}
	startDate, endDate := unit.getDefaultDateRange()
	if req.StartDate != nil {
		startDate = req.StartDate
	}
	if req.EndDate != nil {
		endDate = req.EndDate
	}
	resp, err := h.clickService.GetTimeSeriesSummary(
		c.Request.Context(),
		shortCode,
		startDate,
		endDate,
		unit.getGroupExpr(),
		unit.getPeriodExpr())
	if err != nil {
		c.Error(err)
		return
	}
	resp.Unit = string(unit)
	c.JSON(http.StatusOK, resp)
}

func (h *StatsHandler) GetGeographicStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "short code is required",
		})
		return
	}
	resp, err := h.clickService.GetGeographicStats(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *StatsHandler) GetPlatformStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "short code is required",
		})
		return
	}
	resp, err := h.clickService.GetPlatformStats(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
