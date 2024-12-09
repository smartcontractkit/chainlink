package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type updateRMNConfigTestCase struct {
	useMCMS bool
	name    string
}

func TestUpdateRMNConfig(t *testing.T) {
	t.Parallel()
	testCases := []updateRMNConfigTestCase{
		{
			useMCMS: true,
			name:    "with MCMS",
		},
		{
			useMCMS: false,
			name:    "without MCMS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateRMNConfig(t, tc)
		})
	}
}

func updateRMNConfig(t *testing.T, tc updateRMNConfigTestCase) {
	e := NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), memory.MemoryEnvironmentConfig{
		Chains:     2,
		Nodes:      4,
		Bootstraps: 1,
	}, nil)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	contractsByChain := make(map[uint64][]common.Address)
	rmnRemoteAddressesByChain := buildRMNRemoteAddressPerChain(e.Env, state)
	for chainSelector, rmnRemoteAddress := range rmnRemoteAddressesByChain {
		contractsByChain[chainSelector] = []common.Address{rmnRemoteAddress}
	}

	contractsByChain[e.HomeChainSel] = append(contractsByChain[e.HomeChainSel], state.Chains[e.HomeChainSel].RMNHome.Address())

	timelocksPerChain := buildTimelockPerChain(e.Env, state)
	if tc.useMCMS {
		// This is required because RMNHome is initially owned by the deployer
		_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
				Config: commonchangeset.TransferToMCMSWithTimelockConfig{
					ContractsByChain: contractsByChain,
					MinDelay:         0,
				},
			},
		})
	}

	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	previousCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	previousActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	var mcmsConfig *MCMSConfig = nil

	if tc.useMCMS {
		mcmsConfig = &MCMSConfig{
			MinDelay: 0,
		}
	}

	setRMNHomeCandidateConfig := SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
			Nodes:          []rmn_home.RMNHomeNode{},
			OffchainConfig: []byte(""),
		},
		RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
			SourceChains:   []rmn_home.RMNHomeSourceChain{},
			OffchainConfig: []byte(""),
		},
		MCMSConfig: mcmsConfig,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewSetRMNHomeCandidateConfigChangeset),
			Config:    setRMNHomeCandidateConfig,
		},
	})

	require.NoError(t, err)

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
		DigestToPromote:   currentCandidateDigest,
		MCMSConfig:        mcmsConfig,
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
		F:          0,
		MCMSConfig: mcmsConfig,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(NewSetRMNRemoteConfigChangeset),
			Config:    setRemoteConfig,
		},
	})

	require.NoError(t, err)
	rmnRemotePerChain := buildRMNRemotePerChain(e.Env, state)
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

func buildRMNRemoteAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]common.Address {
	rmnRemotePerChain := buildRMNRemotePerChain(e, state)
	rmnRemoteAddressPerChain := make(map[uint64]common.Address)
	for chain, remote := range rmnRemotePerChain {
		if remote == nil {
			continue
		}
		rmnRemoteAddressPerChain[chain] = remote.Address()
	}
	return rmnRemoteAddressPerChain
}
