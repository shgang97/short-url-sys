package detector

import (
	"strings"

	"github.com/ua-parser/uap-go/uaparser"
)

type DeviceDetector interface {
	Parse(userAgent string) (*DeviceInfo, error)
}

type DefaultDeviceDetector struct {
	parser *uaparser.Parser
}

func NewDefaultDeviceDetector() DeviceDetector {
	parser := uaparser.NewFromSaved() // 使用内置规则
	return &DefaultDeviceDetector{parser: parser}
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceType string `json:"device_type"`
	Browser    string `json:"browser"`
	OS         string `json:"os"`
}

func (d *DefaultDeviceDetector) Parse(userAgent string) (*DeviceInfo, error) {
	if userAgent == "" {
		return &DeviceInfo{
			DeviceType: "other",
			Browser:    "unknown",
			OS:         "unknown",
		}, nil
	}

	client := d.parser.Parse(userAgent)

	// 设备类型映射
	deviceType := d.mapDeviceType(client.Device.Family, userAgent)

	return &DeviceInfo{
		DeviceType: deviceType,
		Browser:    d.normalizeBrowser(client.UserAgent.Family),
		OS:         d.normalizeOS(client.Os.Family),
	}, nil
}

// 设备类型映射
func (d *DefaultDeviceDetector) mapDeviceType(deviceFamily, userAgent string) string {
	deviceFamily = strings.ToLower(deviceFamily)
	userAgent = strings.ToLower(userAgent)

	// 机器人检测
	if strings.Contains(userAgent, "bot") ||
		strings.Contains(userAgent, "crawler") ||
		strings.Contains(userAgent, "spider") {
		return "bot"
	}

	// 设备类型判断
	switch {
	case strings.Contains(deviceFamily, "mobile") ||
		strings.Contains(userAgent, "mobile") ||
		strings.Contains(userAgent, "android") ||
		strings.Contains(userAgent, "iphone"):
		return "mobile"
	case strings.Contains(deviceFamily, "tablet") ||
		strings.Contains(userAgent, "tablet") ||
		strings.Contains(userAgent, "ipad"):
		return "tablet"
	case strings.Contains(deviceFamily, "desktop") ||
		strings.Contains(deviceFamily, "pc"):
		return "desktop"
	default:
		// 通过其他特征判断
		if strings.Contains(userAgent, "windows") ||
			strings.Contains(userAgent, "macintosh") ||
			strings.Contains(userAgent, "linux") {
			return "desktop"
		}
		return "other"
	}
}

// 浏览器名称规范化
func (d *DefaultDeviceDetector) normalizeBrowser(browser string) string {
	browser = strings.ToLower(browser)

	switch {
	case strings.Contains(browser, "chrome"):
		return "Chrome"
	case strings.Contains(browser, "firefox"):
		return "Firefox"
	case strings.Contains(browser, "safari") && !strings.Contains(browser, "chrome"):
		return "Safari"
	case strings.Contains(browser, "edge"):
		return "Edge"
	case strings.Contains(browser, "opera"):
		return "Opera"
	default:
		return truncateString(browser, 100) // 截断到数据库字段长度
	}
}

// 操作系统规范化
func (d *DefaultDeviceDetector) normalizeOS(os string) string {
	os = strings.ToLower(os)

	switch {
	case strings.Contains(os, "windows"):
		return "Windows"
	case strings.Contains(os, "mac"):
		return "macOS"
	case strings.Contains(os, "linux"):
		return "Linux"
	case strings.Contains(os, "android"):
		return "Android"
	case strings.Contains(os, "ios"):
		return "iOS"
	default:
		return truncateString(os, 100)
	}
}

// 截断字符串
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}
