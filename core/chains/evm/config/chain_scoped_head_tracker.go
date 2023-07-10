package config

import (
	"time"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
)

type headTrackerConfig struct {
	c v2.HeadTracker
}

func (h *headTrackerConfig) HistoryDepth() uint32 {
	return *h.c.HistoryDepth
}

func (h *headTrackerConfig) MaxBufferSize() uint32 {
	return *h.c.MaxBufferSize
}

func (h *headTrackerConfig) SamplingInterval() time.Duration {
	return h.c.SamplingInterval.Duration()
}
