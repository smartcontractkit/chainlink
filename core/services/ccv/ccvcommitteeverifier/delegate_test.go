package ccvcommitteeverifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/verifier/pkg/commit"
	chainlinkconfig "github.com/smartcontractkit/chainlink/v2/core/config"
)

// mockCCV implements config.CCV for tests.
type mockCCV struct {
	secrets []chainlinkconfig.AggregatorSecret
}

func (m *mockCCV) AggregatorSecrets() []chainlinkconfig.AggregatorSecret { return m.secrets }
func (m *mockCCV) IndexerSecret() chainlinkconfig.IndexerSecret           { return nil }

// mockAggregatorSecret implements config.AggregatorSecret for tests.
type mockAggregatorSecret struct {
	verifierID string
	apiKey     string
	apiSecret  string
}

func (s *mockAggregatorSecret) VerifierID() string { return s.verifierID }
func (s *mockAggregatorSecret) APIKey() string      { return s.apiKey }
func (s *mockAggregatorSecret) APISecret() string   { return s.apiSecret }

func TestBuildAggregatorSecrets(t *testing.T) {
	const (
		matchingVerifierID  = "my-committee-verifier-1-v1"
		unrelatedVerifierID = "other-committee-verifier-2-v1"
		aggregatorAddr      = "localhost:9090"
		expectedAPIKey      = "test-api-key"
		expectedAPISecret   = "test-api-secret"
	)

	// ccvConfig with two secrets — one matching, one unrelated — to confirm
	// we select by VerifierID rather than taking the first entry.
	ccvConfig := &mockCCV{
		secrets: []chainlinkconfig.AggregatorSecret{
			&mockAggregatorSecret{
				verifierID: unrelatedVerifierID,
				apiKey:     "other-key",
				apiSecret:  "other-secret",
			},
			&mockAggregatorSecret{
				verifierID: matchingVerifierID,
				apiKey:     expectedAPIKey,
				apiSecret:  expectedAPISecret,
			},
		},
	}

	t.Run("legacy config: aggregator_address + verifier_id, no aggregators list", func(t *testing.T) {
		cfg := commit.Config{
			VerifierID:        matchingVerifierID,
			AggregatorAddress: aggregatorAddr,
		}

		secrets, err := buildAggregatorSecrets(ccvConfig, cfg)
		require.NoError(t, err)

		// Legacy path produces one entry keyed by "" (empty SecretName).
		require.Len(t, secrets, 1)
		got, ok := secrets[""]
		require.True(t, ok, "expected credential under empty-string key for legacy config")
		assert.Equal(t, expectedAPIKey, got.APIKey)
		assert.Equal(t, expectedAPISecret, got.Secret)
	})

	t.Run("new config: aggregators list with SecretName", func(t *testing.T) {
		cfg := commit.Config{
			VerifierID: matchingVerifierID,
			Aggregators: []commit.AggregatorConnection{
				{
					Address:    aggregatorAddr,
					SecretName: matchingVerifierID,
				},
			},
		}

		secrets, err := buildAggregatorSecrets(ccvConfig, cfg)
		require.NoError(t, err)

		require.Len(t, secrets, 1)
		got, ok := secrets[matchingVerifierID]
		require.True(t, ok, "expected credential keyed by SecretName")
		assert.Equal(t, expectedAPIKey, got.APIKey)
		assert.Equal(t, expectedAPISecret, got.Secret)
	})

	t.Run("error: no matching secret in ccvConfig", func(t *testing.T) {
		cfg := commit.Config{
			VerifierID:        "nonexistent-verifier-id",
			AggregatorAddress: aggregatorAddr,
		}

		_, err := buildAggregatorSecrets(ccvConfig, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonexistent-verifier-id")
	})

	t.Run("error: no aggregator configured", func(t *testing.T) {
		cfg := commit.Config{
			VerifierID: matchingVerifierID,
			// Neither AggregatorAddress nor Aggregators set.
		}

		_, err := buildAggregatorSecrets(ccvConfig, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no aggregator configured")
	})
}
