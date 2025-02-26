package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowsConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	w := cfg.Workflows()
	assert.Equal(t, int32(200), w.Limits().Global())
	assert.Equal(t, int32(200), w.Limits().PerOwner())
}
