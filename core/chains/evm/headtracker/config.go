package headtracker

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

// Config represents a subset of options needed by head tracker
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	EvmFinalityDepth() uint32
}

type HeadTrackerConfig interface {
	config.HeadTracker
}

var _ types.Config = (*wrappedConfig)(nil)

// Deprecated - this should be removed once config has been refactored.
type wrappedConfig struct {
	Config
}

func NewWrappedConfig(c Config) *wrappedConfig {
	return &wrappedConfig{c}
}

func (c *wrappedConfig) FinalityDepth() uint32 {
	return c.EvmFinalityDepth()
}
