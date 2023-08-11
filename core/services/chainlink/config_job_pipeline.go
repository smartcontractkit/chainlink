package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ config.JobPipeline = (*jobPipelineConfig)(nil)

type jobPipelineConfig struct {
	c toml.JobPipeline
}

func (j *jobPipelineConfig) DefaultHTTPLimit() int64 {
	return int64(*j.c.HTTPRequest.MaxSize)
}

func (j *jobPipelineConfig) DefaultHTTPTimeout() models.Duration {
	return *j.c.HTTPRequest.DefaultTimeout
}

func (j *jobPipelineConfig) MaxRunDuration() time.Duration {
	return j.c.MaxRunDuration.Duration()
}

func (j *jobPipelineConfig) MaxSuccessfulRuns() uint64 {
	return *j.c.MaxSuccessfulRuns
}

func (j *jobPipelineConfig) ReaperInterval() time.Duration {
	return j.c.ReaperInterval.Duration()
}

func (j *jobPipelineConfig) ReaperThreshold() time.Duration {
	return j.c.ReaperThreshold.Duration()
}

func (j *jobPipelineConfig) ResultWriteQueueDepth() uint64 {
	return uint64(*j.c.ResultWriteQueueDepth)
}

func (j *jobPipelineConfig) ExternalInitiatorsEnabled() bool {
	return *j.c.ExternalInitiatorsEnabled
}
