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
