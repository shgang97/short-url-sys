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
	Brokers  []string `mapstructure:"brokers"`
	ClientID string   `mapstructure:"client_id"`
	version  string   `mapstructure:"version"`
	Consumer Consumer `mapstructure:"consumer"`
}

type Consumer struct {
	GroupID    string   `mapstructure:"group_id"`
	Topics     []string `mapstructure:"topics"`
	AutoOffset string   `mapstructure:"auto_offset"`
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
