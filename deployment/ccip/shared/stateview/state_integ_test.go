//go:build integration && db

package stateview_test

import (
	"context"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestLoadChainState_MultipleFeeQuoters(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	fq1 := utils.RandomAddress().Hex()
	fq2 := utils.RandomAddress().Hex()
	state, err := stateview.LoadChainState(context.Background(), tenv.Env.BlockChains.EVMChains()[tenv.HomeChainSel], map[string]cldf.TypeAndVersion{
		fq1: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0),
		fq2: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_2_0),
	})
	require.NoError(t, err)

	require.Equal(t, fq2, state.FeeQuoter.Address().Hex(), "expected latest fee quoter to be selected")
	require.Equal(t, deployment.Version1_2_0, *state.FeeQuoterVersion, "expected latest fee quoter version to be selected")
}

func TestSmokeState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, err = state.View(&tenv.Env, tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)))
	require.NoError(t, err)
}

func TestMCMSState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNoJobsAndContracts())
	addressbook := cldf.NewMemoryAddressBook()
	newTv := cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	newTv.AddLabel(types.BypasserRole.String())
	newTv.AddLabel(types.CancellerRole.String())
	newTv.AddLabel(types.ProposerRole.String())
	addr := utils.RandomAddress()
	require.NoError(t, addressbook.Save(tenv.HomeChainSel, addr.String(), newTv))
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(addressbook))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].BypasserMcm.Address().String())
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].ProposerMcm.Address().String())
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].CancellerMcm.Address().String())
}

// TODO: add solana state test
