package idgen

import (
	"time"

	"github.com/sony/sonyflake"
)

type SfGenerator struct {
	flake *sonyflake.Sonyflake
}

func NewSfGenerator(cfg *SonyflakeConfig) Generator {
	settings := sonyflake.Settings{
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			return cfg.NodeID, nil
		},
		CheckMachineID: func(id uint16) bool {
			return id <= 1024
		},
	}

	sf := sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("Failed to create Sonyflake generator")
	}
	return &SfGenerator{flake: sf}
}

func (g *SfGenerator) NextId() (uint64, error) {
	return g.flake.NextID()
}

func (g *SfGenerator) String() string {
	return "sonyflake"
}
