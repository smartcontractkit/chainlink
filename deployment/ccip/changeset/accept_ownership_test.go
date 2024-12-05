package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_NewAcceptOwnershipChangeset(t *testing.T) {
	e := NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), memory.MemoryEnvironmentConfig{
		Chains:             2,
		NumOfUsersPerChain: 1,
		Nodes:              4,
		Bootstraps:         1,
	}, &TestConfigs{})
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := maps.Keys(e.Env.Chains)
	source := allChains[0]
	dest := allChains[1]

	timelocks := map[uint64]*gethwrappers.RBACTimelock{
		source: state.EVMState.Chains[source].Timelock,
		dest:   state.EVMState.Chains[dest].Timelock,
	}

	// at this point we have the initial deploys done, now we need to transfer ownership
	// to the timelock contract
	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compose the transfer ownership and accept ownership changesets
	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocks, []commonchangeset.ChangesetApplication{
		// note this doesn't have proposals.
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    genTestTransferOwnershipConfig(e, allChains, state),
		},
	})
	require.NoError(t, err)

	assertTimelockOwnership(t, e, allChains, state)
}

func genTestTransferOwnershipConfig(
	e DeployedEnv,
	chains []uint64,
	state CCIPOnChainState,
) commonchangeset.TransferToMCMSWithTimelockConfig {
	var (
		timelocksPerChain = make(map[uint64]common.Address)
		contracts         = make(map[uint64][]common.Address)
	)

	// chain contracts
	for _, chain := range chains {
		timelocksPerChain[chain] = state.EVMState.Chains[chain].Timelock.Address()
		contracts[chain] = []common.Address{
			state.EVMState.Chains[chain].OnRamp.Address(),
			state.EVMState.Chains[chain].OffRamp.Address(),
			state.EVMState.Chains[chain].FeeQuoter.Address(),
			state.EVMState.Chains[chain].NonceManager.Address(),
			state.EVMState.Chains[chain].RMNRemote.Address(),
		}
	}

	// home chain
	homeChainTimelockAddress := state.EVMState.Chains[e.HomeChainSel].Timelock.Address()
	timelocksPerChain[e.HomeChainSel] = homeChainTimelockAddress
	contracts[e.HomeChainSel] = append(contracts[e.HomeChainSel],
		state.EVMState.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
		state.EVMState.Chains[e.HomeChainSel].CCIPHome.Address(),
		state.EVMState.Chains[e.HomeChainSel].RMNHome.Address(),
	)

	return commonchangeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contracts,
	}
}

// assertTimelockOwnership asserts that the ownership of the contracts has been transferred
// to the appropriate timelock contract on each chain.
func assertTimelockOwnership(
	t *testing.T,
	e DeployedEnv,
	chains []uint64,
	state CCIPOnChainState,
) {
	// check that the ownership has been transferred correctly
	for _, chain := range chains {
		for _, contract := range []common.Address{
			state.EVMState.Chains[chain].OnRamp.Address(),
			state.EVMState.Chains[chain].OffRamp.Address(),
			state.EVMState.Chains[chain].FeeQuoter.Address(),
			state.EVMState.Chains[chain].NonceManager.Address(),
			state.EVMState.Chains[chain].RMNRemote.Address(),
		} {
			owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.Chains[chain].EVMChain.Client)
			require.NoError(t, err)
			require.Equal(t, state.EVMState.Chains[chain].Timelock.Address(), owner)
		}
	}

	// check home chain contracts ownership
	homeChainTimelockAddress := state.EVMState.Chains[e.HomeChainSel].Timelock.Address()
	for _, contract := range []common.Address{
		state.EVMState.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
		state.EVMState.Chains[e.HomeChainSel].CCIPHome.Address(),
		state.EVMState.Chains[e.HomeChainSel].RMNHome.Address(),
	} {
		owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.Chains[e.HomeChainSel].EVMChain.Client)
		require.NoError(t, err)
		require.Equal(t, homeChainTimelockAddress, owner)
	}
}
