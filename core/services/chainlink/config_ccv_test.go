package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	secretsCCV = `
[[CCV.AggregatorSecrets]]
VerifierID = "default-verifier-1"
APIKey = "default-api-key"
APISecret = "default-api-secret"

[[CCV.AggregatorSecrets]]
VerifierID = "secondary-verifier-1"
APIKey = "secondary-api-key"
APISecret = "secondary-api-secret"

[CCV.IndexerSecret]
APIKey = "indexer-api-key"
APISecret = "indexer-api-secret"
`
)

func TestCCVConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		SecretsStrings: []string{secretsCCV},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	require.Len(t, cfg.CCV().AggregatorSecrets(), 2)
	c := cfg.CCV()
	require.Equal(t, "default-verifier-1", c.AggregatorSecrets()[0].VerifierID())
	require.Equal(t, "default-api-key", c.AggregatorSecrets()[0].APIKey())
	require.Equal(t, "default-api-secret", c.AggregatorSecrets()[0].APISecret())
	require.Equal(t, "secondary-verifier-1", c.AggregatorSecrets()[1].VerifierID())
	require.Equal(t, "secondary-api-key", c.AggregatorSecrets()[1].APIKey())
	require.Equal(t, "secondary-api-secret", c.AggregatorSecrets()[1].APISecret())
	require.Equal(t, "indexer-api-key", c.IndexerSecret().APIKey())
	require.Equal(t, "indexer-api-secret", c.IndexerSecret().APISecret())
}
