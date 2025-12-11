package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

const (
	thresholdSecrets = `
[Threshold]
ThresholdKeyShare = "something"
`
)

func TestThresholdConfig(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	opts := GeneralConfigOpts{
		SecretsStrings: []string{thresholdSecrets},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	th := cfg.Threshold()
	assert.Equal(t, "something", th.ThresholdKeyShare())
}
