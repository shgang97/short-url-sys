package idgen

// Generator ID生成器接口
type Generator interface {
	NextId() (uint64, error)
	// String 返回生成器类型
	String() string
}

type GeneratorConfig struct {
	Type      string          `mapstructure:"type"`
	Sonyflake SonyflakeConfig `mapstructure:"snowflake"`
}

type SonyflakeConfig struct {
	NodeID uint16 `mapstructure:"node_id"`
}
