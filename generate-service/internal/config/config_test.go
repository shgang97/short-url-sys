package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  port: 8080
  host: "0.0.0.0"
  mode: "test"
  base_url: "http://127.0.0.1:8080"
database:
  mysql:
    dsn: "test:test@tcp(localhost:3306)/test"
    max_idle_conns: 5
    max_open_conns: 10
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	require.NoError(t, err)
	tmpFile.Close()

	// 测试配置加载
	cfg, err := Load(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证配置值
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "test", cfg.Server.Mode)
	assert.Equal(t, "http://127.0.0.1:8080", cfg.Server.BaseURL)
	assert.Equal(t, "test:test@tcp(localhost:3306)/test", cfg.Database.MySQL.DSN)
}
