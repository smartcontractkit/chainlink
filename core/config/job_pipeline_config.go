package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type JobPipeline interface {
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() sqlutil.Duration
	MaxRunDuration() time.Duration
	MaxSuccessfulRuns() uint64
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
	ResultWriteQueueDepth() uint64
	ExternalInitiatorsEnabled() bool
}
