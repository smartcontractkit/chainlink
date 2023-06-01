package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type JobPipeline interface {
	DefaultHTTPLimit() int64
	DefaultHTTPTimeout() models.Duration
	JobPipelineMaxRunDuration() time.Duration
	JobPipelineMaxSuccessfulRuns() uint64
	JobPipelineReaperInterval() time.Duration
	JobPipelineReaperThreshold() time.Duration
	JobPipelineResultWriteQueueDepth() uint64
}
