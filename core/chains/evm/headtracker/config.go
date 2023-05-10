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

type EvmHtrkConfig types.HtrkConfig

var _ EvmHtrkConfig = (*evmHtrkConfig)(nil)

type evmHtrkConfig struct {
	Config
}

func NewEvmHtrkConfig(c Config) *evmHtrkConfig {
	return &evmHtrkConfig{c}
}

func (c *evmHtrkConfig) FinalityDepth() uint32 {
	return c.EvmFinalityDepth()
}

func (c *evmHtrkConfig) HeadTrackerHistoryDepth() uint32 {
	return c.EvmHeadTrackerHistoryDepth()
}

func (c *evmHtrkConfig) HeadTrackerMaxBufferSize() uint32 {
	return c.EvmHeadTrackerMaxBufferSize()
}

func (c *evmHtrkConfig) HeadTrackerSamplingInterval() time.Duration {
	return c.EvmHeadTrackerSamplingInterval()
}
