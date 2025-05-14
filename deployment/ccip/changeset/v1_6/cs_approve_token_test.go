package v1_6_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestApproveTransferEVMChangeset(t *testing.T) {
	ctx := testcontext.Get(t)
	// Default env just has 2 chains with all contracts
	// deployed but no lanes.
	tenv, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	allChains := maps.Keys(tenv.Env.Chains)
	source := allChains[0]
	dest := allChains[1]
	_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.TokenApproveTransferEVMChangeset),
			v1_6.ApproveTokenEVMConfig{
				ChainSelector: source,
				Amount:        big.NewInt(100),
			},
		),
	)

	require.NoError(t, err)

	// Assert the onramp configuration is as we expect.
	sourceCfg, err := state.Chains[source].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, dest)
	require.NoError(t, err)
	require.Equal(t, state.Chains[source].TestRouter.Address(), sourceCfg.Router)
	require.False(t, sourceCfg.AllowlistEnabled)
	destCfg, err := state.Chains[dest].OnRamp.GetDestChainConfig(&bind.CallOpts{Context: ctx}, source)
	require.NoError(t, err)
	require.Equal(t, state.Chains[dest].Router.Address(), destCfg.Router)
	require.True(t, destCfg.AllowlistEnabled)
}
