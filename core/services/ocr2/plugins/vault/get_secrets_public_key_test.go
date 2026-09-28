package vault

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// TestStateTransitionGetSecrets_IncludesPublicKey verifies the RawVaultPublicKey
// gate: when open, stateTransitionGetSecrets attaches the instance's vault public
// key to the aggregated GetSecrets response; when closed, it does not.
func TestStateTransitionGetSecrets_IncludesPublicKey(t *testing.T) {
	t.Parallel()

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	chosen := []*vaultcommon.Observation{
		{
			Response: &vaultcommon.Observation_GetSecretsResponse{
				GetSecretsResponse: &vaultcommon.GetSecretsResponse{
					Responses: []*vaultcommon.SecretResponse{
						{
							Id: &vaultcommon.SecretIdentifier{Key: "k", Namespace: "n", Owner: "o"},
							Result: &vaultcommon.SecretResponse_Data{
								Data: &vaultcommon.SecretData{EncryptedValue: "deadbeef"},
							},
						},
					},
				},
			},
		},
	}

	t.Run("gate closed omits public key", func(t *testing.T) {
		r.cfg.VaultGetSecretsIncludePublicKey = limits.NewGateLimiter(false)
		o := &vaultcommon.Outcome{}
		r.stateTransitionGetSecrets(t.Context(), chosen, o)
		require.Len(t, o.GetGetSecretsResponse().GetResponses(), 1)
		require.Empty(t, o.GetGetSecretsResponse().GetRawVaultPublicKey())
	})

	t.Run("gate open includes the instance public key", func(t *testing.T) {
		r.cfg.VaultGetSecretsIncludePublicKey = limits.NewGateLimiter(true)
		o := &vaultcommon.Outcome{}
		r.stateTransitionGetSecrets(t.Context(), chosen, o)
		require.Len(t, o.GetGetSecretsResponse().GetResponses(), 1)

		pkb, merr := pk.Marshal()
		require.NoError(t, merr)
		require.Equal(t, hex.EncodeToString(pkb), o.GetGetSecretsResponse().GetRawVaultPublicKey())
	})
}
