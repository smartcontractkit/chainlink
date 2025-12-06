//go:build integration && db

package v1_6_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 0, 0)
}

func TestDeployStaticLinkToken(t *testing.T) {
	t.Parallel()
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithStaticLink())
	// load onchain state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	for _, chain := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		require.NotNil(t, state.Chains[chain].StaticLinkToken)
	}
}
