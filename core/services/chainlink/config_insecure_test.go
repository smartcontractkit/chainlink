package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsecureConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ins := cfg.Insecure()
	assert.False(t, ins.DevWebServer())
	assert.False(t, ins.DisableRateLimiting())
	assert.False(t, ins.OCRDevelopmentMode())
	assert.False(t, ins.InfiniteDepthQueries())
}
