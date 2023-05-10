package types

import "time"

type HtrkConfig interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	FinalityDepth() uint32
	HeadTrackerHistoryDepth() uint32
	HeadTrackerMaxBufferSize() uint32
	HeadTrackerSamplingInterval() time.Duration
}
