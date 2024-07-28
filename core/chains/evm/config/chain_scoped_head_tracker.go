package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type headTrackerConfig struct {
	c toml.HeadTracker
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

func (h *headTrackerConfig) FinalityTagBypass() bool {
	return *h.c.FinalityTagBypass
}

func (h *headTrackerConfig) MaxAllowedFinalityDepth() uint32 {
	return *h.c.MaxAllowedFinalityDepth
}
