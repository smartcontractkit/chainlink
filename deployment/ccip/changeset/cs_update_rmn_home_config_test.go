package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	mcms "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUpdateRMNHomeConfig(t *testing.T) {
	t.Parallel()
	e := NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), memory.MemoryEnvironmentConfig{
		Chains:     2,
		Nodes:      4,
		Bootstraps: 1,
	}, nil)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	// This is required because RMNHome is initially owner by the deployer
	err = transferOwnershipToHomeChainTimelock(t, e, state)
	require.NoError(t, err)

	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	previousCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	previousActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	configInput := SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig:   NewTestRMNStaticConfig(),
		RMNDynamicConfig:  NewTestRMNDynamicConfig(),
	}

	timelocksPerChain := buildTimelockPerChain(e, state)
	commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewSetRMNHomeCandidateConfigChangeset),
			Config:    configInput,
		},
	})

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	currentActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	require.NotEqual(t, previousCandidateDigest, currentCandidateDigest)
	require.Equal(t, previousActiveDigest, currentActiveDigest)

	promoteConfigInput := PromoteRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
	}

	commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewPromoteCandidateConfigChangeset),
			Config:    promoteConfigInput,
		},
	})

	require.NoError(t, err)
	currentActiveDigest, err = rmnHome.GetActiveDigest(nil)

	require.NoError(t, err)
	require.NotEqual(t, previousActiveDigest, currentActiveDigest)
}

func buildTimelockPerChain(e DeployedEnv, state CCIPOnChainState) map[uint64]*mcms.RBACTimelock {
	timelocksPerChain := make(map[uint64]*mcms.RBACTimelock)
	for _, chain := range e.Env.Chains {
		timelocksPerChain[chain.Selector] = state.Chains[chain.Selector].Timelock
	}
	return timelocksPerChain
}

func transferOwnershipToHomeChainTimelock(t *testing.T, e DeployedEnv, state CCIPOnChainState) error {
	rmnHome := state.Chains[e.HomeChainSel].RMNHome
	timelockAddress := state.Chains[e.HomeChainSel].Timelock.Address()
	proposerMcm := state.Chains[e.HomeChainSel].ProposerMcm
	_, err := commonchangeset.NewTransferOwnershipChangeset(e.Env, commonchangeset.TransferOwnershipConfig{
		Contracts: map[uint64][]commonchangeset.OwnershipTransferrer{
			e.HomeChainSel: {rmnHome},
		},
		OwnersPerChain: map[uint64]common.Address{
			e.HomeChainSel: timelockAddress,
		},
	})

	timelocksPerChain := buildTimelockPerChain(e, state)

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.NewAcceptOwnershipChangeset),
			Config: commonchangeset.AcceptOwnershipConfig{
				Contracts: map[uint64][]commonchangeset.OwnershipAcceptor{
					e.HomeChainSel: {rmnHome},
				},
				OwnersPerChain: map[uint64]common.Address{
					e.HomeChainSel: timelockAddress,
				},
				ProposerMCMSes: map[uint64]*mcms.ManyChainMultiSig{
					e.HomeChainSel: proposerMcm,
				},
			},
		},
	})

	return err
}
