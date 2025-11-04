package handler

import (
	"fmt"
	"generate-service/internal/model"
	linkSrc "generate-service/internal/service/link"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	linkService linkSrc.Service
	baseURL     string
}

func NewLinkHandler(linkService linkSrc.Service, baseURL string) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
		baseURL:     baseURL,
	}
}

func (h *LinkHandler) CreateShortURL(c *gin.Context) {
	var req model.CreateShortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	link, err := h.linkService.CreateShortURL(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	resp := model.CreateShortResponse{
		ShortURL:  h.buildShortURL(link.ShortCode),
		LongURL:   link.LongURL,
		ShortCode: link.ShortCode,
		CreatedAt: link.CreatedAt,
		ExpiresAt: link.ExpiresAt,
		QRCodeURL: h.buildQRCodeURL(link.ShortCode),
	}
	c.JSON(http.StatusCreated, resp)
}

// GetLinkInfo
// @Router /api/v1/links/{code} [get]
func (h *LinkHandler) GetLinkInfo(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_short_code",
			Message: "Short code is required",
		})
		return
	}

	linkInfo, err := h.linkService.GetLinkInfo(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, linkInfo)
}

// UpdateLink
// @Router /api/v1/links/{code} [put]
func (h *LinkHandler) UpdateLink(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_short_code",
			Message: "Short code is required",
		})
		return
	}
	var req model.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	linkInfo, err := h.linkService.UpdateLink(c.Request.Context(), shortCode, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, linkInfo)
}

// DeleteLink
// @Router /api/v1/links/{code} [delete]
func (h *LinkHandler) DeleteLink(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_short_code",
			Message: "Short code is required",
		})
		return
	}
	err := h.linkService.DeleteLink(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

// ListLinks
// @Router /api/v1/links [get]
func (h *LinkHandler) ListLinks(c *gin.Context) {
	var req model.ListLinksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid query parameters",
		})
		return
	}
	resp, err := h.linkService.ListLinks(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// BatchCreate
// @Router /api/v1/short/batch [post]
func (h *LinkHandler) BatchCreate(c *gin.Context) {
	var req model.BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}
	resp, err := h.linkService.BatchCreate(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// 构建完整的短链URL
func (h *LinkHandler) buildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/%s", h.baseURL, shortCode)
}

// 构建二维码URL
func (h *LinkHandler) buildQRCodeURL(shortCode string) string {
	return fmt.Sprintf("%s/api/v1/links/%s/qrcode", h.baseURL, shortCode)
}
