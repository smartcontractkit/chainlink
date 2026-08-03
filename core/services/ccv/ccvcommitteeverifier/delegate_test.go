package ccvcommitteeverifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/verifier/pkg/commit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccv/ccvcommitteeverifier"
	clservices "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const secretsCCV = `
[[CCV.AggregatorSecrets]]
VerifierID = "my-committee-verifier-1-v1"
APIKey = "test-api-key"
APISecret = "test-api-secret"

[[CCV.AggregatorSecrets]]
VerifierID = "other-committee-verifier-2-v1"
APIKey = "other-key"
APISecret = "other-secret"
`

func newTestCCVConfig(t *testing.T) clservices.GeneralConfig {
	t.Helper()
	opts := clservices.GeneralConfigOpts{
		SecretsStrings: []string{secretsCCV},
	}
	cfg, err := opts.New()
	require.NoError(t, err)
	return cfg
}

func TestBuildAggregatorSecrets(t *testing.T) {
	t.Parallel()

	const (
		matchingVerifierID = "my-committee-verifier-1-v1"
		aggregatorAddr     = "localhost:9090"
		expectedAPIKey     = "test-api-key"
		expectedAPISecret  = "test-api-secret"
	)

	t.Run("legacy config: aggregator_address + verifier_id, no aggregators list", func(t *testing.T) {
		t.Parallel()
		cfg := commit.Config{
			VerifierID:        matchingVerifierID,
			AggregatorAddress: aggregatorAddr,
		}

		secrets, err := ccvcommitteeverifier.BuildAggregatorSecrets(newTestCCVConfig(t).CCV(), cfg)
		require.NoError(t, err)

		// Legacy path produces one entry keyed by "" (empty SecretName).
		require.Len(t, secrets, 1)
		got, ok := secrets[""]
		require.True(t, ok, "expected credential under empty-string key for legacy config")
		assert.Equal(t, expectedAPIKey, got.APIKey)
		assert.Equal(t, expectedAPISecret, got.Secret)
	})

	t.Run("new config: aggregators list with SecretName", func(t *testing.T) {
		t.Parallel()
		cfg := commit.Config{
			VerifierID: matchingVerifierID,
			Aggregators: []commit.AggregatorConnection{
				{
					Address:    aggregatorAddr,
					SecretName: matchingVerifierID,
				},
			},
		}

		secrets, err := ccvcommitteeverifier.BuildAggregatorSecrets(newTestCCVConfig(t).CCV(), cfg)
		require.NoError(t, err)

		require.Len(t, secrets, 1)
		got, ok := secrets[matchingVerifierID]
		require.True(t, ok, "expected credential keyed by SecretName")
		assert.Equal(t, expectedAPIKey, got.APIKey)
		assert.Equal(t, expectedAPISecret, got.Secret)
	})

	t.Run("error: no matching secret in ccvConfig", func(t *testing.T) {
		t.Parallel()
		cfg := commit.Config{
			VerifierID:        "nonexistent-verifier-id",
			AggregatorAddress: aggregatorAddr,
		}

		_, err := ccvcommitteeverifier.BuildAggregatorSecrets(newTestCCVConfig(t).CCV(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonexistent-verifier-id")
	})

	t.Run("error: no aggregator configured", func(t *testing.T) {
		t.Parallel()
		cfg := commit.Config{
			VerifierID: matchingVerifierID,
			// Neither AggregatorAddress nor Aggregators set.
		}

		_, err := ccvcommitteeverifier.BuildAggregatorSecrets(newTestCCVConfig(t).CCV(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no aggregator configured")
	})
}
