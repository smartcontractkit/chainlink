package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	mcms "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
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
	err = transferOwnershipForRMNRemote(t, e, state)
	require.NoError(t, err)

	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	previousCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	previousActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	setRMNHomeCandidateConfig := SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig:   NewTestRMNStaticConfig(),
		RMNDynamicConfig:  NewTestRMNDynamicConfig(),
	}

	timelocksPerChain := buildTimelockPerChain(e.Env, state)
	commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewSetRMNHomeCandidateConfigChangeset),
			Config:    setRMNHomeCandidateConfig,
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

	promoteConfig := PromoteRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewPromoteCandidateConfigChangeset),
			Config:    promoteConfig,
		},
	})

	require.NoError(t, err)
	currentActiveDigest, err = rmnHome.GetActiveDigest(nil)

	require.NoError(t, err)
	require.NotEqual(t, previousActiveDigest, currentActiveDigest)

	setRemoteConfig := SetRMNRemoteConfig{
		HomeChainSelector: e.HomeChainSel,
		Signers: []rmn_remote.RMNRemoteSigner{
			{
				OnchainPublicKey: common.Address{},
				NodeIndex:        0,
			},
		},
		F: 0,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewSetRMNRemoteConfigChangeset),
			Config:    setRemoteConfig,
		},
	})

	require.NoError(t, err)
	rmnRemotePerChain := buildRemoteRemotePerChain(e.Env, state)
	for _, rmnRemote := range rmnRemotePerChain {
		remoteConfigSetEvents, err := rmnRemote.FilterConfigSet(nil, nil)
		require.NoError(t, err)
		var lastEvent *rmn_remote.RMNRemoteConfigSet
		for remoteConfigSetEvents.Next() {
			lastEvent = remoteConfigSetEvents.Event
		}
		require.NotNil(t, lastEvent)
		require.Equal(t, lastEvent.Config.RmnHomeContractConfigDigest, currentActiveDigest)
	}
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

	timelocksPerChain := buildTimelockPerChain(e.Env, state)

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

func transferOwnershipForRMNRemote(t *testing.T, e DeployedEnv, state CCIPOnChainState) error {
	rmnRemotePerChain := buildRemoteRemotePerChain(e.Env, state)
	timelockAddressPerChain := buildTimelockAddressPerChain(e.Env, state)
	timelocksPerChain := buildTimelockPerChain(e.Env, state)
	proposers := buildProposerPerChain(e.Env, state)
	for chain, rmnRemote := range rmnRemotePerChain {
		_, err := commonchangeset.NewTransferOwnershipChangeset(e.Env, commonchangeset.TransferOwnershipConfig{
			Contracts: map[uint64][]commonchangeset.OwnershipTransferrer{
				chain: {rmnRemote},
			},
			OwnersPerChain: timelockAddressPerChain,
		})
		if err != nil {
			return err
		}

		_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(commonchangeset.NewAcceptOwnershipChangeset),
				Config: commonchangeset.AcceptOwnershipConfig{
					Contracts: map[uint64][]commonchangeset.OwnershipAcceptor{
						chain: {rmnRemote},
					},
					OwnersPerChain: timelockAddressPerChain,
					ProposerMCMSes: proposers,
				},
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}
