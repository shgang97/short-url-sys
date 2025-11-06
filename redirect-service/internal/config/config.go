package config

import "time"

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
	Brokers       []string `mapstructure:"brokers"`
	ClientID      string   `mapstructure:"client_id"`
	Version       string   `mapstructure:"version"`
	GroupID       string   `mapstructure:"group_id"`
	Topics        []string `mapstructure:"topics"`
	FetchMaxBytes int32    `mapstructure:"fetch_max_bytes"`
	Consumer      Consumer `mapstructure:"consumer"`
}

type Consumer struct {
	AutoCommit         bool          `mapstructure:"auto_commit"`
	AutoCommitInterval int64         `mapstructure:"auto_commit_interval"`
	AutoOffset         string        `mapstructure:"auto_offset"`
	SessionTimeout     time.Duration `mapstructure:"session_timeout"`
}

// GenerateService generate-service 客户端配置
type GenerateService struct {
	Address string        `mapstructure:"address"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type Config struct {
	Server          ServerConfig    `mapstructure:"server"`
	Redis           RedisConfig     `mapstructure:"redis"`
	Kafka           KafkaConfig     `mapstructure:"kafka"`
	Cache           CacheConfig     `mapstructure:"cache"`
	GenerateService GenerateService `mapstructure:"generate_service"`
}
