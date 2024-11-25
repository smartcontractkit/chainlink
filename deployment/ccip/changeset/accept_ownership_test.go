package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_NewAcceptOwnershipChangeset(t *testing.T) {
	ctx := tests.Context(t)
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 2, 4)
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := maps.Keys(e.Env.Chains)
	source := allChains[0]
	dest := allChains[1]

	newAddresses := deployment.NewMemoryAddressBook()
	err = deployPrerequisiteChainContracts(e.Env, newAddresses, allChains, nil)
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))

	mcmConfig := commontypes.MCMSWithTimelockConfig{
		Canceller:         commonchangeset.SingleGroupMCMS(t),
		Bypasser:          commonchangeset.SingleGroupMCMS(t),
		Proposer:          commonchangeset.SingleGroupMCMS(t),
		TimelockExecutors: e.Env.AllDeployerKeys(),
		TimelockMinDelay:  big.NewInt(0),
	}
	out, err := commonchangeset.DeployMCMSWithTimelock(e.Env, map[uint64]commontypes.MCMSWithTimelockConfig{
		source: mcmConfig,
		dest:   mcmConfig,
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(out.AddressBook))
	newAddresses = deployment.NewMemoryAddressBook()
	tokenConfig := NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	ocrParams := make(map[uint64]CCIPOCRParams)
	for _, chain := range allChains {
		ocrParams[chain] = DefaultOCRParams(e.FeedChainSel, nil, nil)
	}
	err = deployCCIPContracts(e.Env, newAddresses, NewChainsConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: allChains,
		TokenConfig:    tokenConfig,
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
		OCRParams:      ocrParams,
	})
	require.NoError(t, err)

	// at this point we have the initial deploys done, now we need to transfer ownership
	// to the timelock contract
	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	for _, chain := range allChains {
		deployerBalance, err := e.Env.Chains[chain].Client.BalanceAt(
			ctx, e.Env.Chains[chain].DeployerKey.From, nil)
		require.NoError(t, err)
		t.Log("deployer balance is:", assets.NewWei(deployerBalance).String())

		// OnRamp
		tx, err := state.Chains[chain].OnRamp.TransferOwnership(
			e.Env.Chains[source].DeployerKey,
			state.Chains[source].Timelock.Address(),
		)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)

		// OffRamp
		tx, err = state.Chains[chain].OffRamp.TransferOwnership(
			e.Env.Chains[source].DeployerKey,
			state.Chains[source].Timelock.Address(),
		)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)

		// FeeQuoter
		tx, err = state.Chains[chain].FeeQuoter.TransferOwnership(
			e.Env.Chains[source].DeployerKey,
			state.Chains[source].Timelock.Address(),
		)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)

		// NonceManager
		tx, err = state.Chains[chain].NonceManager.TransferOwnership(
			e.Env.Chains[source].DeployerKey,
			state.Chains[source].Timelock.Address(),
		)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)

		// RMNRemote
		tx, err = state.Chains[chain].RMNRemote.TransferOwnership(
			e.Env.Chains[source].DeployerKey,
			state.Chains[source].Timelock.Address(),
		)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
	}

	// now we can accept ownership via the proposal.
	_, err = commonchangeset.ApplyChangesets(t, e.Env, map[uint64]*gethwrappers.RBACTimelock{
		source: state.Chains[source].Timelock,
		dest:   state.Chains[dest].Timelock,
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewAcceptOwnershipChangeset),
			Config:    AcceptOwnershipConfig{State: state, ChainSelector: source},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(NewAcceptOwnershipChangeset),
			Config:    AcceptOwnershipConfig{State: state, ChainSelector: dest},
		},
	})
	require.NoError(t, err)

	// check that the ownership has been transferred correctly
	for _, chain := range allChains {
		onRampOwner, err := state.Chains[chain].OnRamp.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[source].Timelock.Address(), onRampOwner)

		offRampOwner, err := state.Chains[chain].OffRamp.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[source].Timelock.Address(), offRampOwner)

		feeQuoterOwner, err := state.Chains[chain].FeeQuoter.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[source].Timelock.Address(), feeQuoterOwner)

		nonceManagerOwner, err := state.Chains[chain].NonceManager.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[source].Timelock.Address(), nonceManagerOwner)

		rmnRemoteOwner, err := state.Chains[chain].RMNRemote.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[source].Timelock.Address(), rmnRemoteOwner)
	}
}
