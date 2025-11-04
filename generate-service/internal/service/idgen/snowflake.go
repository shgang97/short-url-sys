package idgen

import (
	"github.com/shgang97/sys-collections/snowflake"
)

type Snowflake struct {
	Node *snowflake.Node
}

func (s *Snowflake) NextId() (uint64, error) {
	id, err := s.Node.NextId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (s *Snowflake) String() string {
	return "snowflake"
}

func NewSnowflake(nodeID int64) (*Snowflake, error) {
	node, err := snowflake.Init(nodeID)
	if err != nil {
		return nil, err
	}
	return &Snowflake{Node: node}, nil
}
