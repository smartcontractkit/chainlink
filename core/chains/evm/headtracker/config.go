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

type EvmConfig types.Config

var _ EvmConfig = (*evmConfig)(nil)

type evmConfig struct {
	Config
}

func NewEvmConfig(c Config) *evmConfig {
	return &evmConfig{c}
}

func (c *evmConfig) FinalityDepth() uint32 {
	return c.EvmFinalityDepth()
}

func (c *evmConfig) HeadTrackerHistoryDepth() uint32 {
	return c.EvmHeadTrackerHistoryDepth()
}

func (c *evmConfig) HeadTrackerMaxBufferSize() uint32 {
	return c.EvmHeadTrackerMaxBufferSize()
}

func (c *evmConfig) HeadTrackerSamplingInterval() time.Duration {
	return c.EvmHeadTrackerSamplingInterval()
}
