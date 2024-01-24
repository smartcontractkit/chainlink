package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestJobPipelineConfigTest(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	jp := cfg.JobPipeline()

	assert.Equal(t, int64(100*utils.MB), jp.DefaultHTTPLimit())
	d, err := commonconfig.NewDuration(1 * time.Minute)
	require.NoError(t, err)
	assert.Equal(t, d, jp.DefaultHTTPTimeout())
	assert.Equal(t, 1*time.Hour, jp.MaxRunDuration())
	assert.Equal(t, uint64(123456), jp.MaxSuccessfulRuns())
	assert.Equal(t, 4*time.Hour, jp.ReaperInterval())
	assert.Equal(t, 168*time.Hour, jp.ReaperThreshold())
	assert.Equal(t, uint64(10), jp.ResultWriteQueueDepth())
	assert.True(t, jp.ExternalInitiatorsEnabled())
}
