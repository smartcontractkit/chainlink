package vault

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func testChosenGetSecretsObservations() []*vaultcommon.Observation {
	return []*vaultcommon.Observation{
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
}

// TestStateTransitionGetSecrets_AttachesPublicKey verifies that
// stateTransitionGetSecrets attaches the caller-supplied (quorum-aggregated)
// public key to the response, and attaches nothing when it is empty. The key is
// no longer read from node-local config, so this is a pure function of its args.
func TestStateTransitionGetSecrets_AttachesPublicKey(t *testing.T) {
	t.Parallel()

	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	chosen := testChosenGetSecretsObservations()

	// No aggregated key: nothing attached.
	oEmpty := &vaultcommon.Outcome{}
	r.stateTransitionGetSecrets(chosen, oEmpty, "")
	require.Len(t, oEmpty.GetGetSecretsResponse().GetResponses(), 1)
	require.Empty(t, oEmpty.GetGetSecretsResponse().GetRawVaultPublicKey())

	// Aggregated key present: attached verbatim.
	oKey := &vaultcommon.Outcome{}
	r.stateTransitionGetSecrets(chosen, oKey, "abc123")
	require.Len(t, oKey.GetGetSecretsResponse().GetResponses(), 1)
	require.Equal(t, "abc123", oKey.GetGetSecretsResponse().GetRawVaultPublicKey())
}

// TestObservedVaultPublicKey verifies a node broadcasts its DKG public key only
// when the include-public-key gate is open.
func TestObservedVaultPublicKey(t *testing.T) {
	t.Parallel()

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withOnchainCfg(4, 1))

	pkb, merr := pk.Marshal()
	require.NoError(t, merr)

	r.cfg.VaultGetSecretsIncludePublicKey = limits.NewGateLimiter(false)
	require.Empty(t, r.observedVaultPublicKey(t.Context()))

	r.cfg.VaultGetSecretsIncludePublicKey = limits.NewGateLimiter(true)
	require.Equal(t, pkb, r.observedVaultPublicKey(t.Context()))
}

// TestAggregateVaultPublicKey verifies StateTransition selects the public key
// agreed by at least F+1 observations and returns "" below that quorum.
func TestAggregateVaultPublicKey(t *testing.T) {
	t.Parallel()

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	pkb, merr := pk.Marshal()
	require.NoError(t, merr)
	wantHex := hex.EncodeToString(pkb)

	// N=13, F=4 => threshold F+1 = 5.
	r := newTestReportingPlugin(t, withOnchainCfg(13, 4))

	obsWithKey := func(n int) map[uint8]*vaultcommon.Observations {
		m := map[uint8]*vaultcommon.Observations{}
		for i := range n {
			m[uint8(i)] = &vaultcommon.Observations{RawVaultPublicKey: pkb}
		}
		return m
	}

	// Below quorum: no key.
	require.Empty(t, r.aggregateVaultPublicKey(obsWithKey(4)))

	// At quorum (F+1): key selected.
	require.Equal(t, wantHex, r.aggregateVaultPublicKey(obsWithKey(5)))

	// Nodes that omit the key (gate off) do not count toward the quorum.
	mixed := obsWithKey(4)
	mixed[100] = &vaultcommon.Observations{} // no key
	mixed[101] = &vaultcommon.Observations{} // no key
	require.Empty(t, r.aggregateVaultPublicKey(mixed))
}
