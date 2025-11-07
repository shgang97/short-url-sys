package config

import (
	"redirect-service/internal/pkg/idgen"
	"time"
)

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type CacheConfig struct {
	TTL    int    `mapstructure:"ttl"`
	Prefix string `mapstructure:"prefix"`
}

type KafkaConfig struct {
	Brokers            []string `mapstructure:"brokers"`
	ClientID           string   `mapstructure:"client_id"`
	Version            string   `mapstructure:"version"`
	GroupID            string   `mapstructure:"group_id"`
	Topics             []string `mapstructure:"topics"`
	FetchMaxBytes      int32    `mapstructure:"fetch_max_bytes"`
	Consumer           Consumer `mapstructure:"consumer"`
	Producer           Producer `mapstructure:"producer"`
	NetMaxOpenRequests int      `mapstructure:"net_max_open_requests"`
}

type Consumer struct {
	AutoCommit         bool          `mapstructure:"auto_commit"`
	AutoCommitInterval int64         `mapstructure:"auto_commit_interval"`
	AutoOffset         string        `mapstructure:"auto_offset"`
	SessionTimeout     time.Duration `mapstructure:"session_timeout"`
}

type Producer struct {
	RequiredAcks int `mapstructure:"required_acks"` // 对应 sarama.WaitForAll（-1）
	Compression  int `mapstructure:"compression"`   // 对应 sarama.CompressionSnappy（2）
	Flush        struct {
		Frequency time.Duration `mapstructure:"frequency"` // 时间字符串，如 "500ms"
	} `mapstructure:"flush"`
	Return struct {
		Successes bool `mapstructure:"successes"` // 是否返回成功信息
		Errors    bool `mapstructure:"errors"`    // 是否返回错误信息
	} `mapstructure:"return"`
	Retry struct {
		Max int `mapstructure:"max"` // 最大重试次数
	} `mapstructure:"retry"`
	Idempotent bool `mapstructure:"idempotent"` // 是否启用幂等性
}

// GenerateService generate-service 客户端配置
type GenerateService struct {
	Address string        `mapstructure:"address"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type GeoIPConfig struct {
	DBPath string `mapstructure:"db_path"`
}

type Config struct {
	Server          ServerConfig          `mapstructure:"server"`
	Redis           RedisConfig           `mapstructure:"redis"`
	Kafka           KafkaConfig           `mapstructure:"kafka"`
	Cache           CacheConfig           `mapstructure:"cache"`
	GenerateService GenerateService       `mapstructure:"generate_service"`
	GeoIP           GeoIPConfig           `mapstructure:"geo_ip"`
	Generator       idgen.GeneratorConfig `mapstructure:"id_generator"`
}
