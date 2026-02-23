package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

func TestLocalSecretsFetcher_GetSecrets(t *testing.T) {
	secrets := map[string]string{
		"api-key":     "sk-abc123",
		"db-password": "hunter2",
	}
	fetcher := NewLocalSecretsFetcher(secrets)

	t.Run("returns known secrets", func(t *testing.T) {
		resp, err := fetcher.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
			Requests: []*sdkpb.SecretRequest{
				{Id: "api-key", Namespace: "default"},
				{Id: "db-password"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp, 2)

		s0 := resp[0].GetSecret()
		require.NotNil(t, s0)
		assert.Equal(t, "api-key", s0.Id)
		assert.Equal(t, "default", s0.Namespace)
		assert.Equal(t, "sk-abc123", s0.Value)

		s1 := resp[1].GetSecret()
		require.NotNil(t, s1)
		assert.Equal(t, "db-password", s1.Id)
		assert.Equal(t, "hunter2", s1.Value)
	})

	t.Run("returns error for unknown secret", func(t *testing.T) {
		resp, err := fetcher.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
			Requests: []*sdkpb.SecretRequest{
				{Id: "nonexistent"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp, 1)

		errResp := resp[0].GetError()
		require.NotNil(t, errResp)
		assert.Equal(t, "nonexistent", errResp.Id)
		assert.Contains(t, errResp.Error, "not found")
	})

	t.Run("handles mixed known and unknown", func(t *testing.T) {
		resp, err := fetcher.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
			Requests: []*sdkpb.SecretRequest{
				{Id: "api-key"},
				{Id: "missing"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp, 2)

		assert.NotNil(t, resp[0].GetSecret())
		assert.NotNil(t, resp[1].GetError())
	})

	t.Run("empty map returns errors for all", func(t *testing.T) {
		emptyFetcher := NewLocalSecretsFetcher(map[string]string{})
		resp, err := emptyFetcher.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{
			Requests: []*sdkpb.SecretRequest{
				{Id: "anything"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp, 1)
		assert.NotNil(t, resp[0].GetError())
	})
}
