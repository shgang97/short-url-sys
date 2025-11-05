package constants

import "time"

const (
	// MaxCacheTTL 最大缓存时间 30 天
	MaxCacheTTL = 30 * 24 * time.Hour
	// DefaultCacheTTL 默认缓存时间 7 天
	DefaultCacheTTL = 7 * 24 * time.Hour
	// ExpirationToleranceCacheTTL 如果计算出的时间过期了，依然设置，缓存时间设置为24小时
	ExpirationToleranceCacheTTL = 1 * 24 * time.Hour
)
