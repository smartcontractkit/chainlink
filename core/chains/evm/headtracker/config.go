package headtracker

import "time"

//go:generate mockery --name Config --output ./mocks/ --case=underscore

// Config represents a subset of options needed by head tracker
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
}
