package config

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type JobPipeline interface {
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() commonconfig.Duration
	MaxRunDuration() time.Duration
	MaxSuccessfulRuns() uint64
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
	ResultWriteQueueDepth() uint64
	ExternalInitiatorsEnabled() bool
	VerboseLogging() bool
}
