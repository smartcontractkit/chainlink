package v2

import (
	"time"
)

type nodePoolConfig struct {
	c NodePool
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
