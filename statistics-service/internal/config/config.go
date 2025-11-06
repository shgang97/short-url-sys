package config

import (
	"statistics-service/internal/consumer"
	"statistics-service/internal/pkg/database"
	"statistics-service/internal/pkg/idgen"
	"statistics-service/internal/pkg/logger"
)

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type Config struct {
	Server    ServerConfig          `mapstructure:"server"`
	Log       logger.Config         `mapstructure:"log"`
	MySQL     database.MySQLConfig  `mapstructure:"mysql"`
	Kafka     consumer.KafkaConfig  `mapstructure:"kafka"`
	Generator idgen.GeneratorConfig `mapstructure:"id_generator"`
}
