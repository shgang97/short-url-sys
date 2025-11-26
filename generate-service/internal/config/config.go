package config

import (
	"generate-service/internal/service/register"
	"time"
)

type ServerConfig struct {
	GRPC    GRPCConfig `mapstructure:"grpc"`
	HTTP    HTTPConfig `mapstructure:"http"`
	BaseURL string     `mapstructure:"base_url"`
}

// HTTPConfig HTTP 服务器配置
type HTTPConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

// GRPCConfig gRPC 服务器配置
type GRPCConfig struct {
	Host                 string        `mapstructure:"host"`
	Port                 string        `mapstructure:"port"`
	MaxConcurrentStreams uint32        `mapstructure:"max_concurrent_streams"`
	MaxRecvMsgSize       int           `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize       int           `mapstructure:"max_send_msg_size"`
	ConnectionTimeout    time.Duration `mapstructure:"connection_timeout"`
}

type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type IDGeneratorConfig struct {
	Type      string          `mapstructure:"type"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
}

type SnowflakeConfig struct {
	NodeID int64 `mapstructure:"node_id"`
}

type CacheConfig struct {
	TTL    int    `mapstructure:"ttl"`
	Prefix string `mapstructure:"prefix"`
}

type LogConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Version  string   `mapstructure:"version"`
	ClientID string   `mapstructure:"client_id"`
	Producer Producer `mapstructure:"producer"`
	Net      Net      `mapstructure:"net"`
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

type Net struct {
	MaxOpenRequests int `mapstructure:"max_open_requests"` // 最大并发请求数（幂等性需设为1）
}

type RateLimitConfig struct {
	RequestPerMinute int `mapstructure:"requests_per_minute"`
}

type Config struct {
	Server      ServerConfig        `mapstructure:"server"`
	Database    DatabaseConfig      `mapstructure:"database"`
	Redis       RedisConfig         `mapstructure:"redis"`
	IdGenerator IDGeneratorConfig   `mapstructure:"id_generator"`
	Cache       CacheConfig         `mapstructure:"cache"`
	Log         LogConfig           `mapstructure:"log"`
	Kafka       KafkaConfig         `mapstructure:"kafka"`
	Etcd        register.EtcdConfig `mapstructure:"etcd"`
	RateLimit   RateLimitConfig     `mapstructure:"rate_limit"`
}
