package headtracker

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

// Config represents a subset of options needed by head tracker
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
}

var _ types.Config = (*wrappedConfig)(nil)

type wrappedConfig struct {
	Config
}

func NewWrappedConfig(c Config) *wrappedConfig {
	return &wrappedConfig{c}
}

func (c *wrappedConfig) FinalityDepth() uint32 {
	return c.EvmFinalityDepth()
}

func (c *wrappedConfig) HeadTrackerHistoryDepth() uint32 {
	return c.EvmHeadTrackerHistoryDepth()
}

func (c *wrappedConfig) HeadTrackerMaxBufferSize() uint32 {
	return c.EvmHeadTrackerMaxBufferSize()
}

func (c *wrappedConfig) HeadTrackerSamplingInterval() time.Duration {
	return c.EvmHeadTrackerSamplingInterval()
}
