package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type headTrackerConfig struct {
	c toml.HeadTracker
	blockEmissionIdleWarningThreshold time.Duration
	finalityDepth uint32
	finalityTagEnabled bool
}

func (h *headTrackerConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	return h.blockEmissionIdleWarningThreshold
}

func (h *headTrackerConfig) FinalityDepth() uint32 {
	return h.finalityDepth
}

func (h *headTrackerConfig) FinalityTagEnabled() bool {
	return h.finalityTagEnabled
}

func (h *headTrackerConfig) HistoryDepth() uint32 {
	return *h.c.HistoryDepth
}

func (h *headTrackerConfig) SamplingInterval() time.Duration {
	return h.c.SamplingInterval.Duration()
}
