package errors

import "fmt"

// 定义错误类型
var (
	ErrLinkNotFound     = NewBusinessError("link not found")
	ErrLinkExpired      = NewBusinessError("link expired")
	ErrLinkDisabled     = NewBusinessError("link disabled")
	ErrInvalidURL       = NewBusinessError("invalid URL")
	ErrShortCodeExists  = NewBusinessError("short code already exists")
	ErrInvalidShortCode = NewBusinessError("invalid short code")
	ErrBreakerOpen      = NewBusinessError("circuit breaker is open")
)

type BusinessError struct {
	message string
}

// NewBusinessError 业务错误
func NewBusinessError(message string) *BusinessError {
	return &BusinessError{message: message}
}

func (e *BusinessError) Error() string {
	return e.message
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on filed %s: %s", e.Field, e.Message)
}

// RepositoryError 数据层错误
type RepositoryError struct {
	Operation string
	Err       error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("repository error on operation %s: %s", e.Operation, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}
