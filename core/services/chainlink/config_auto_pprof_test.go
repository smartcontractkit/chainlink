package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestAutoPprofTest(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ap := cfg.AutoPprof()
	assert.True(t, ap.Enabled())
	assert.Equal(t, "prof/root", ap.ProfileRoot())
	assert.Equal(t, 1*time.Minute, ap.PollInterval().Duration())
	assert.Equal(t, 12*time.Second, ap.GatherDuration().Duration())
	assert.Equal(t, 13*time.Second, ap.GatherTraceDuration().Duration())
	assert.Equal(t, utils.FileSize(1*utils.GB), ap.MaxProfileSize())
	assert.Equal(t, 7, ap.CPUProfileRate())
	assert.Equal(t, 9, ap.MemProfileRate())
	assert.Equal(t, 5, ap.BlockProfileRate())
	assert.Equal(t, 2, ap.MutexProfileFraction())
	assert.Equal(t, utils.FileSize(1*utils.GB), ap.MemThreshold())
	assert.Equal(t, 999, ap.GoroutineThreshold())
}
