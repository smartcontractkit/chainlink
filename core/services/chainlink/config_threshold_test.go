package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	thresholdSecrets = `
[Threshold]
ThresholdKeyShare = "something"
`
)

func TestThresholdConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		SecretsStrings: []string{thresholdSecrets},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	th := cfg.Threshold()
	assert.Equal(t, "something", th.ThresholdKeyShare())
}
