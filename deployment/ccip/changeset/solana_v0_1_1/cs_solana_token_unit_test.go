//go:build !integration

package solana_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestDeployLinkToken(t *testing.T) {
	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
	))
	require.NoError(t, err)

	// solana test
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken), commonchangeset.DeploySolanaLinkTokenConfig{
			ChainSelector: selector,
			TokenPrivKey:  solLinkTokenPrivKey,
			TokenDecimals: 9,
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.NotEmpty(t, addrs)
}
