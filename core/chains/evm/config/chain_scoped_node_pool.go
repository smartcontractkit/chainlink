package config

import (
	"time"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
)

type nodePoolConfig struct {
	c v2.NodePool
}

func (n *nodePoolConfig) PollFailureThreshold() uint32 {
	return *n.c.PollFailureThreshold
}

func (n *nodePoolConfig) PollInterval() time.Duration {
	return n.c.PollInterval.Duration()
}

func (n *nodePoolConfig) SelectionMode() string {
	return *n.c.SelectionMode
}

func (n *nodePoolConfig) SyncThreshold() uint32 {
	return *n.c.SyncThreshold
}
