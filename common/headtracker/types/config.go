package types

import "time"

type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	FinalityDepth() uint32
	FinalityTagEnabled() bool
}

type HeadTrackerConfig interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
}
