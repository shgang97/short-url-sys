package handler

import (
	"fmt"
	"generate-service/internal/model"
	"generate-service/internal/service/link"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	linkService link.Service
	baseURL     string
}

func NewQRCodeHandler(linkService link.Service, baseURL string) *QRCodeHandler {
	return &QRCodeHandler{
		linkService: linkService,
		baseURL:     baseURL,
	}
}

// GenerateQRCode
// @Router /api/v1/links/qrcode/{code} [get]
func (h *QRCodeHandler) GenerateQRCode(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid short code",
			Message: "Short code is required",
		})
		return
	}

	// 验证短码是否存在
	_, err := h.linkService.GetLinkInfo(c.Request.Context(), shortCode)
	if err != nil {
		c.Error(err)
		return
	}

	// 获取二维码尺寸
	size := 256
	if sizeParam := c.Query("size"); sizeParam != "" {
		if s, err := strconv.Atoi(sizeParam); err == nil && s > 0 {
			size = s
		}
	}

	// 限制尺寸
	if size > 1024 {
		size = 1024
	}

	// 生成短链接URL
	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortCode)

	// 生成二维码
	png, err := qrcode.Encode(shortURL, qrcode.Medium, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "qrcode generation failed",
			Message: "Failed to generate QRCode",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(len(png)))
	c.Header("cache-control", "public, max-age=3600") // 缓存1小时
	c.Data(http.StatusOK, "image/png", png)
}
