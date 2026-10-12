package solana_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	linkchangesets "github.com/smartcontractkit/cld-changesets/tokens/link/changesets"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deploylink"
)

// TestVerifyRemoteDecimalsSolana covers shared.VerifyRemoteDecimals against a real SPL mint.
//
// It lives here because the Solana container harness does, and shared cannot import it
// without an import cycle. Only a Solana container is started; the full CCIP environment is
// not needed to read a mint.
func TestVerifyRemoteDecimalsSolana(t *testing.T) {
	t.Parallel()

	const tokenDecimals = 6

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
	))
	require.NoError(t, err)

	mintKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	require.NoError(t, rt.Exec(
		runtime.ChangesetTask(deploylink.DeployLinkTokenChangeset{}, linkchangesets.DeployLinkTokenInput{
			Solana: map[uint64]linkchangesets.SolanaLinkConfig{
				selector: {TokenPrivKey: mintKey, TokenDecimals: tokenDecimals},
			},
		}),
	))

	env := rt.Environment()

	// The pool stores a remote SVM token as the raw 32-byte mint pubkey.
	mintBytes := mintKey.PublicKey().Bytes()

	t.Run("Success - matching decimals verify", func(t *testing.T) {
		t.Parallel()

		verified, err := shared.VerifyRemoteDecimals(t.Context(), env, selector, mintBytes, tokenDecimals)
		require.NoError(t, err)
		require.True(t, verified, "a mint present in the environment should be verifiable")
	})

	t.Run("Failure - mismatched decimals", func(t *testing.T) {
		t.Parallel()

		verified, err := shared.VerifyRemoteDecimals(t.Context(), env, selector, mintBytes, 18)
		require.Error(t, err)
		require.False(t, verified)
		require.Contains(t, err.Error(), "supplied 18 decimals but the remote token reports 6")
	})

	t.Run("Failure - wrong address length", func(t *testing.T) {
		t.Parallel()

		verified, err := shared.VerifyRemoteDecimals(t.Context(), env, selector, mintBytes[:20], tokenDecimals)
		require.Error(t, err)
		require.False(t, verified)
		require.Contains(t, err.Error(), "expected a 32 byte SVM mint address")
	})
}
