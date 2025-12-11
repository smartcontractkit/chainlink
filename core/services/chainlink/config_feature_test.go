package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestFeatureConfig(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	f := cfg.Feature()
	assert.True(t, f.LogPoller())
	assert.True(t, f.FeedsManager())
	assert.True(t, f.UICSAKeys())
	assert.True(t, f.MultiFeedsManagers())
}
