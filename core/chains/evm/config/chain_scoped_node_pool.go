package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type nodePoolConfig struct {
	c toml.NodePool
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

func (n *nodePoolConfig) LeaseDuration() time.Duration {
	return n.c.LeaseDuration.Duration()
}
