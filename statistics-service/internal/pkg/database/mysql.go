package database

import (
	"fmt"
	"time"

	log "statistics-service/internal/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDB struct {
	DB *gorm.DB
}

type MySQLConfig struct {
	DSN             string          `mapstructure:"dsn"`
	MaxIdleConns    int             `mapstructure:"max_idle_conns"`
	MaxOpenConns    int             `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration   `mapstructure:"conn_max_lifetime"`
	LogLevel        logger.LogLevel `mapstructure:"log_level"`
}

func NewMySQLDB(cfg MySQLConfig) (*MySQLDB, error) {
	// 创建 GROM 配置
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.DSN), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL DB: %w", err)
	}

	log.Logger.Info("Successfully connected MySQL DB")
	return &MySQLDB{DB: db}, nil
}

func (m *MySQLDB) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck 检查数据库连接状态
func (m *MySQLDB) HealthCheck() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
