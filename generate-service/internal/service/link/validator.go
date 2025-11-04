package link

import (
	"generate-service/internal/pkg/errors"
	"net/url"
	"strings"
)

// URLValidator URL 验证器
type URLValidator struct {
	allowedSchemes map[string]bool
	maxURLLength   int
}

// NewURLValidator 创建URL验证器
func NewURLValidator() *URLValidator {
	return &URLValidator{
		allowedSchemes: map[string]bool{
			"http":  true,
			"https": true,
		},
		maxURLLength: 2048,
	}
}

// Validate 验证URL格式和安全性
func (v *URLValidator) Validate(rawURL string) error {
	// 检查URL长度
	if len(rawURL) > v.maxURLLength {
		return errors.NewBusinessError("URL too long")
	}

	// 解析URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return errors.ErrInvalidURL
	}

	// 检查协议
	if !v.allowedSchemes[strings.ToLower(u.Scheme)] {
		return errors.NewBusinessError("URL scheme not allowed")
	}

	// 检查主机名
	if u.Hostname() == "" {
		return errors.NewBusinessError("URL hostname is required")
	}

	// 检查保留域名（可选，防止滥用）
	//if v.isReservedDomain(u.Hostname()) {
	//	return errors.NewBusinessError("URL hostname is reserved")
	//}
	return nil
}

// 检查是否为保留域名
func (v *URLValidator) isReservedDomain(hostname string) bool {
	reservedDomains := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
	}
	for _, domain := range reservedDomains {
		if strings.Contains(hostname, domain) {
			return true
		}
	}
	return false
}

// NormalizeURL 标准化URL
func (v *URLValidator) NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// 标准化协议和主机名
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// 移除默认端口
	if u.Scheme == "http" && u.Port() == "80" {
		u.Host = u.Hostname()
	} else if u.Scheme == "https" && u.Port() == "443" {
		u.Host = u.Hostname()
	}

	// 标准化路径（去除末尾斜杠）
	u.Path = strings.TrimPrefix(u.Path, "/")

	// 排序查询参数（可选）
	// ?a=1&b=2 和 ?b=2&a=1 应该被视为同一个URL
	if u.RawQuery != "" {
		queryParams, _ := url.ParseQuery(u.RawQuery)
		u.RawQuery = queryParams.Encode()
	}

	return u.String(), nil
}
