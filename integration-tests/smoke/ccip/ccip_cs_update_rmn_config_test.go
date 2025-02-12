package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
)

var (
	rmnStaging1 = changeset.RMNNopConfig{
		NodeIndex:           0,
		PeerId:              deployment.MustPeerIDFromString("p2p_12D3KooWRXxZq3pd4a3ZGkKj7Nt1SQQrnB8CuvbPnnV9KVeMeWqg"),
		OffchainPublicKey:   [32]byte(common.FromHex("0xb34944857a42444d1b285d7940d6e06682309e0781e43a69676ee9f85c73c2d1")),
		EVMOnChainPublicKey: common.HexToAddress("0x5af8ee32316a6427f169a45fdc1b3a91a85ac459e3c1cb91c69e1c51f0c1fc21"),
	}
	rmnStaging2 = changeset.RMNNopConfig{
		NodeIndex:           1,
		PeerId:              deployment.MustPeerIDFromString("p2p_12D3KooWEmdxYQFsRbD9aFczF32zA3CcUwuSiWCk2CrmACo4v9RL"),
		OffchainPublicKey:   [32]byte(common.FromHex("0x68d9f3f274e3985528a923a9bace3d39c55dd778b187b4120b384cc48c892859")),
		EVMOnChainPublicKey: common.HexToAddress("0x858589216956f482a0f68b282a7050af4cd48ed2"),
	}
	rmnStaging3 = changeset.RMNNopConfig{
		NodeIndex:           2,
		PeerId:              deployment.MustPeerIDFromString("p2p_12D3KooWJS42cNXKJvj6DeZnxEX7aGxhEuap6uNFrz554AbUDw6Q"),
		OffchainPublicKey:   [32]byte(common.FromHex("0x5af8ee32316a6427f169a45fdc1b3a91a85ac459e3c1cb91c69e1c51f0c1fc21")),
		EVMOnChainPublicKey: common.HexToAddress("0x7c5e94162c6fabbdeb3bfe83ae532846e337bfae"),
	}
)

type updateRMNConfigTestCase struct {
	useMCMS bool
	name    string
	nops    []changeset.RMNNopConfig
}

func TestUpdateRMNConfig(t *testing.T) {
	t.Parallel()
	testCases := []updateRMNConfigTestCase{
		{
			useMCMS: true,
			name:    "with MCMS",
			nops:    []changeset.RMNNopConfig{rmnStaging1, rmnStaging2, rmnStaging3},
		},
		{
			useMCMS: false,
			name:    "without MCMS",
			nops:    []changeset.RMNNopConfig{rmnStaging1, rmnStaging2, rmnStaging3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			updateRMNConfig(t, tc)
		})
	}
}

func TestSetDynamicConfig(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	nops := []changeset.RMNNopConfig{rmnStaging1, rmnStaging2, rmnStaging3}
	nodes := make([]rmn_home.RMNHomeNode, 0, len(nops))
	for _, nop := range nops {
		nodes = append(nodes, nop.ToRMNHomeNode())
	}

	setRMNHomeCandidateConfig := changeset.SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
			Nodes:          nodes,
			OffchainConfig: []byte(""),
		},
		RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
			SourceChains:   []rmn_home.RMNHomeSourceChain{},
			OffchainConfig: []byte(""),
		},
	}

	_, err = changeset.SetRMNHomeCandidateConfigChangeset(e.Env, setRMNHomeCandidateConfig)
	require.NoError(t, err)

	candidate, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)

	promoteCandidateConfig := changeset.PromoteRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		DigestToPromote:   candidate,
	}

	_, err = changeset.PromoteRMNHomeCandidateConfigChangeset(e.Env, promoteCandidateConfig)
	require.NoError(t, err)

	active, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	setDynamicConfig := changeset.SetRMNHomeDynamicConfigConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
			SourceChains: []rmn_home.RMNHomeSourceChain{
				{
					ChainSelector:       e.HomeChainSel,
					ObserverNodesBitmap: big.NewInt(1),
				},
			},
			OffchainConfig: []byte(""),
		},
		ActiveDigest: active,
	}

	_, err = changeset.SetRMNHomeDynamicConfigChangeset(e.Env, setDynamicConfig)
	require.NoError(t, err)

	dynamicConfig, err := rmnHome.GetConfig(nil, active)
	require.NoError(t, err)

	require.True(t, dynamicConfig.Ok)
	require.Equal(t, setDynamicConfig.RMNDynamicConfig, dynamicConfig.VersionedConfig.DynamicConfig)
}

func TestRevokeConfig(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	nops := []changeset.RMNNopConfig{rmnStaging1, rmnStaging2, rmnStaging3}
	nodes := make([]rmn_home.RMNHomeNode, 0, len(nops))
	for _, nop := range nops {
		nodes = append(nodes, nop.ToRMNHomeNode())
	}

	setRMNHomeCandidateConfig := changeset.SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
			Nodes:          nodes,
			OffchainConfig: []byte(""),
		},
		RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
			SourceChains:   []rmn_home.RMNHomeSourceChain{},
			OffchainConfig: []byte(""),
		},
	}

	_, err = changeset.SetRMNHomeCandidateConfigChangeset(e.Env, setRMNHomeCandidateConfig)
	require.NoError(t, err)

	candidate, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)

	revokeCandidateConfig := changeset.RevokeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		CandidateDigest:   candidate,
	}

	_, err = changeset.RevokeRMNHomeCandidateConfigChangeset(e.Env, revokeCandidateConfig)
	require.NoError(t, err)

	newCandidate, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	require.NotEqual(t, candidate, newCandidate)
}

func updateRMNConfig(t *testing.T, tc updateRMNConfigTestCase) {
	e, _ := testhelpers.NewMemoryEnvironment(t)

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	contractsByChain := make(map[uint64][]common.Address)
	rmnRemoteAddressesByChain := buildRMNRemoteAddressPerChain(e.Env, state)
	for chainSelector, rmnRemoteAddress := range rmnRemoteAddressesByChain {
		contractsByChain[chainSelector] = []common.Address{rmnRemoteAddress}
	}

	contractsByChain[e.HomeChainSel] = append(contractsByChain[e.HomeChainSel], state.Chains[e.HomeChainSel].RMNHome.Address())

	timelocksPerChain := changeset.BuildTimelockPerChain(e.Env, state)
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
		require.NoError(t, err)
	}

	rmnHome := state.Chains[e.HomeChainSel].RMNHome

	previousCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	previousActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	var mcmsConfig *changeset.MCMSConfig

	if tc.useMCMS {
		mcmsConfig = &changeset.MCMSConfig{
			MinDelay: 0,
		}
	}

	nodes := make([]rmn_home.RMNHomeNode, 0, len(tc.nops))
	for _, nop := range tc.nops {
		nodes = append(nodes, nop.ToRMNHomeNode())
	}

	setRMNHomeCandidateConfig := changeset.SetRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
			Nodes:          nodes,
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
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNHomeCandidateConfigChangeset),
			Config:    setRMNHomeCandidateConfig,
		},
	})

	require.NoError(t, err)

	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	require.NoError(t, err)
	currentActiveDigest, err := rmnHome.GetActiveDigest(nil)
	require.NoError(t, err)

	require.NotEqual(t, previousCandidateDigest, currentCandidateDigest)
	require.Equal(t, previousActiveDigest, currentActiveDigest)

	promoteConfig := changeset.PromoteRMNHomeCandidateConfig{
		HomeChainSelector: e.HomeChainSel,
		DigestToPromote:   currentCandidateDigest,
		MCMSConfig:        mcmsConfig,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.PromoteRMNHomeCandidateConfigChangeset),
			Config:    promoteConfig,
		},
	})

	require.NoError(t, err)
	currentActiveDigest, err = rmnHome.GetActiveDigest(nil)

	require.NoError(t, err)
	require.NotEqual(t, previousActiveDigest, currentActiveDigest)

	signers := make([]rmn_remote.RMNRemoteSigner, 0, len(tc.nops))
	for _, nop := range tc.nops {
		signers = append(signers, nop.ToRMNRemoteSigner())
	}

	remoteConfigs := make(map[uint64]changeset.RMNRemoteConfig, len(e.Env.Chains))
	for _, chain := range e.Env.Chains {
		remoteConfig := changeset.RMNRemoteConfig{
			Signers: signers,
			F:       0,
		}

		remoteConfigs[chain.Selector] = remoteConfig
	}

	setRemoteConfig := changeset.SetRMNRemoteConfig{
		HomeChainSelector: e.HomeChainSel,
		RMNRemoteConfigs:  remoteConfigs,
		MCMSConfig:        mcmsConfig,
	}

	_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNRemoteConfigChangeset),
			Config:    setRemoteConfig,
		},
	})

	require.NoError(t, err)
	rmnRemotePerChain := changeset.BuildRMNRemotePerChain(e.Env, state)
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

func buildRMNRemoteAddressPerChain(e deployment.Environment, state changeset.CCIPOnChainState) map[uint64]common.Address {
	rmnRemotePerChain := changeset.BuildRMNRemotePerChain(e, state)
	rmnRemoteAddressPerChain := make(map[uint64]common.Address)
	for chain, remote := range rmnRemotePerChain {
		if remote == nil {
			continue
		}
		rmnRemoteAddressPerChain[chain] = remote.Address()
	}
	return rmnRemoteAddressPerChain
}

func TestSetRMNRemoteOnRMNProxy(t *testing.T) {
	t.Parallel()
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNoJobsAndContracts())
	allChains := e.Env.AllChainSelectors()
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	var err error
	prereqCfgs := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, c := range e.Env.AllChainSelectors() {
		mcmsCfg[c] = proposalutils.SingleGroupTimelockConfig(t)
		prereqCfgs = append(prereqCfgs, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: c,
		})
	}
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    allChains,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfgs,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)
	contractsByChain := make(map[uint64][]common.Address)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	for _, chain := range allChains {
		rmnProxy := state.Chains[chain].RMNProxy
		require.NotNil(t, rmnProxy)
		contractsByChain[chain] = []common.Address{rmnProxy.Address()}
	}
	timelockContractsPerChain := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	allContractParams := make(map[uint64]changeset.ChainContractParams)
	for _, chain := range allChains {
		timelockContractsPerChain[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain].Timelock,
			CallProxy: state.Chains[chain].CallProxy,
		}
		allContractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, timelockContractsPerChain, []commonchangeset.ChangesetApplication{
		// transfer ownership of RMNProxy to timelock
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config: commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MinDelay:         0,
			},
		},

		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: allContractParams,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNRemoteOnRMNProxyChangeset),
			Config: changeset.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: allChains,
				MCMSConfig: &changeset.MCMSConfig{
					MinDelay: 0,
				},
			},
		},
	})
	require.NoError(t, err)
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	for _, chain := range allChains {
		rmnProxy := state.Chains[chain].RMNProxy
		proxyOwner, err := rmnProxy.Owner(nil)
		require.NoError(t, err)
		require.Equal(t, state.Chains[chain].Timelock.Address(), proxyOwner)
		rmnAddr, err := rmnProxy.GetARM(nil)
		require.NoError(t, err)
		require.Equal(t, rmnAddr, state.Chains[chain].RMNRemote.Address())
	}
}
