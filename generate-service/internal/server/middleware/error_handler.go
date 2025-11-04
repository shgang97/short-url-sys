package middleware

import (
	"generate-service/internal/model"
	"generate-service/internal/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// 记录错误日志
				gin.DefaultErrorWriter.Write([]byte(err.Error() + "\n"))
			}

			// 处理最后一个错误
			lastError := c.Errors.Last().Err

			var statusCode int
			var errorResp model.ErrorResponse

			switch err := lastError.(type) {
			case *errors.BusinessError:
				statusCode = http.StatusBadRequest
				errorResp = model.ErrorResponse{
					Error:   "business_error",
					Message: err.Error(),
				}
			case *errors.ValidationError:
				statusCode = http.StatusBadRequest
				errorResp = model.ErrorResponse{
					Error:   "validation_error",
					Message: err.Error(),
					Code:    err.Field,
				}
			case *errors.RepositoryError:
				statusCode = http.StatusInternalServerError
				errorResp = model.ErrorResponse{
					Error:   "internal_error",
					Message: "Internal server error",
				}
			default:
				// 处理已知的业务错误
				switch lastError {
				case errors.ErrLinkNotFound:
					statusCode = http.StatusNotFound
					errorResp = model.ErrorResponse{
						Error:   "link_not_found",
						Message: "Short link not found",
					}
				case errors.ErrLinkExpired:
					statusCode = http.StatusGone
					errorResp = model.ErrorResponse{
						Error:   "link_expired",
						Message: "Short link has expired",
					}
				case errors.ErrLinkDisabled:
					statusCode = http.StatusForbidden
					errorResp = model.ErrorResponse{
						Error:   "link_disabled",
						Message: "Short link is disabled",
					}
				case errors.ErrInvalidURL:
					statusCode = http.StatusBadRequest
					errorResp = model.ErrorResponse{
						Error:   "invalid_url",
						Message: "Invalid URL format",
					}
				case errors.ErrShortCodeExists:
					statusCode = http.StatusConflict
					errorResp = model.ErrorResponse{
						Error:   "short_code_exists",
						Message: "Short code already exists",
					}
				case errors.ErrInvalidShortCode:
					statusCode = http.StatusBadRequest
					errorResp = model.ErrorResponse{
						Error:   "invalid_short_code",
						Message: "Invalid short code format",
					}
				default:
					statusCode = http.StatusInternalServerError
					errorResp = model.ErrorResponse{
						Error:   "internal_error",
						Message: "Internal server error",
					}
				}
			}

			c.JSON(statusCode, errorResp)
			c.Abort()
			return
		}
	}
}
