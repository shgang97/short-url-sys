package idgen

import (
	"generate-service/internal/config"
	"generate-service/internal/pkg/database"
)

// NewIDGenerator 根据配置创建ID生成器
func NewIDGenerator(cfg *config.IDGeneratorConfig, redisClient *database.RedisClient) (Generator, error) {
	switch cfg.Type {
	case "redis":
		return NewRedisGenerator(redisClient.Client, "short_url:id_counter"), nil
	case "snowflake":
		return NewSnowflake(cfg.Snowflake.NodeID)
	default:
		return NewSnowflake(1)
	}
}
