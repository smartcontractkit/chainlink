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

	// compose the transfer ownership and accept ownership changesets
	_, err = commonchangeset.ApplyChangesets(t, e.Env, map[uint64]*gethwrappers.RBACTimelock{
		source: state.Chains[source].Timelock,
		dest:   state.Chains[dest].Timelock,
	}, []commonchangeset.ChangesetApplication{
		// note this doesn't have proposals.
		{
			Changeset: commonchangeset.WrapChangeSet(NewTransferOwnershipChangeset),
			Config: TransferOwnershipConfig{
				State:             state,
				ChainSelectors:    allChains,
				HomeChainSelector: e.HomeChainSel,
			},
		},
		// this has proposals, ApplyChangesets will sign & execute them.
		// in practice, signing and executing are separated processes.
		{
			Changeset: commonchangeset.WrapChangeSet(NewAcceptOwnershipChangeset),
			Config: AcceptOwnershipConfig{
				State:             state,
				ChainSelectors:    allChains,
				HomeChainSelector: e.HomeChainSel,
			},
		},
	})
	require.NoError(t, err)

	// check that the ownership has been transferred correctly
	for _, chain := range allChains {
		for _, contract := range []ownershipTransferrer{
			state.Chains[chain].OnRamp,
			state.Chains[chain].OffRamp,
			state.Chains[chain].FeeQuoter,
			state.Chains[chain].NonceManager,
			state.Chains[chain].RMNRemote,
		} {
			owner, err := contract.Owner(&bind.CallOpts{
				Context: ctx,
			})
			require.NoError(t, err)
			require.Equal(t, state.Chains[chain].Timelock.Address(), owner)
		}
	}

	// check home chain contracts ownership
	homeChainTimelockAddress := state.Chains[e.HomeChainSel].Timelock.Address()
	for _, contract := range []ownershipTransferrer{
		state.Chains[e.HomeChainSel].CapabilityRegistry,
		state.Chains[e.HomeChainSel].CCIPHome,
		state.Chains[e.HomeChainSel].RMNHome,
	} {
		owner, err := contract.Owner(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err)
		require.Equal(t, homeChainTimelockAddress, owner)
	}
}
