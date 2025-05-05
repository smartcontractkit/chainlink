package ccip

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	v1_5testhelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

var (
	evm2EVMOnRampABI = abihelpers.MustParseABI(evm_2_evm_onramp.EVM2EVMOnRampABI)
	onRampABI        = abihelpers.MustParseABI(onramp.OnRampABI)
)

// TestV1_5_Message_RMNRemote this test verify that 1.5 lane can send message when using RMNRemote
func TestV1_5_Message_RMNRemote(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					// setting dummy voters, we will permabless this later
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5testhelpers.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			v1_5.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	oldState, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	evmContractParams := make(map[uint64]v1_6.ChainContractParams)
	evmChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.Chains[chain]; ok {
			evmChains = append(evmChains, chain)
		}
	}
	for _, chain := range evmChains {
		evmContractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	var apps []commonchangeset.ConfiguredChangeSet
	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: evmContractParams,
			},
		),
	}...)
	// reload state after adding lanes
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, apps)
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)

	_, err = cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset).Apply(e.Env,
		v1_6.SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
	require.NoError(t, err)

	// send continuous messages in real router until done is closed
	// send a message from the other lane src1 -> dest
	sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, oldState,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(oldState.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	require.NoError(t, err)
	v1_5testhelpers.WaitForCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)
}

// TestV1_5_Message_RMNRemote this test verify that 1.5 lane can be cursed when using RMNRemote
func TestV1_5_Message_RMNRemote_Curse(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					// setting dummy voters, we will permabless this later
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5testhelpers.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			v1_5.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	oldState, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	evmContractParams := make(map[uint64]v1_6.ChainContractParams)
	evmChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.Chains[chain]; ok {
			evmChains = append(evmChains, chain)
		}
	}
	for _, chain := range evmChains {
		evmContractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	var apps []commonchangeset.ConfiguredChangeSet
	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: evmContractParams,
			},
		),
	}...)
	// reload state after adding lanes
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, apps)
	require.NoError(t, err)

	// reload state after adding lanes
	tEnv.UpdateDeployedEnvironment(e)

	_, err = cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset).Apply(e.Env,
		v1_6.SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
	require.NoError(t, err)

	// send continuous messages in real router until done is closed
	// send a message from the other lane src1 -> dest
	sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, oldState,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(oldState.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)

	_, err = cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset).Apply(e.Env, v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{v1_6.CurseChain(e.Env.AllChainSelectors()[0])},
		Reason:       "Curse test",
	})
	require.NoError(t, err)

	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	require.NoError(t, err)
	v1_5testhelpers.WaitForNoCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)
}

// TestV1_5_Message_RMNRemote this test verify that 1.5 lane can be uncuresed when using RMNRemote
func TestV1_5_Message_RMNRemote_Curse_Uncurse(t *testing.T) {
	t.Skip("Skipping since its flakey, need to fix")
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					// setting dummy voters, we will permabless this later
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	src1, dest := allChains[0], allChains[1]
	pairs := []testhelpers.SourceDestPair{
		{SourceChainSelector: src1, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5testhelpers.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			v1_5.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	oldState, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	evmContractParams := make(map[uint64]v1_6.ChainContractParams)
	evmChains := []uint64{}
	for _, chain := range allChains {
		if _, ok := e.Env.Chains[chain]; ok {
			evmChains = append(evmChains, chain)
		}
	}
	for _, chain := range evmChains {
		evmContractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	var apps []commonchangeset.ConfiguredChangeSet
	apps = append(apps, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     e.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.Chains[e.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      e.HomeChainSel,
				ContractParamsPerChain: evmContractParams,
			},
		),
	}...)
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, apps)
	require.NoError(t, err)
	// reload state after adding lanes

	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)

	_, err = cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset).Apply(e.Env,
		v1_6.SetRMNRemoteOnRMNProxyConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
	require.NoError(t, err)

	// send continuous messages in real router until done is closed
	// send a message from the other lane src1 -> dest
	sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, oldState,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(oldState.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)

	_, err = cldf.CreateLegacyChangeSet(v1_6.RMNCurseChangeset).Apply(e.Env, v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{v1_6.CurseChain(e.Env.AllChainSelectors()[0])},
		Reason:       "Curse test",
	})
	require.NoError(t, err)

	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	v1_5testhelpers.WaitForNoCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)

	commitFound := make(chan struct{})
	go func() {
		v1_5testhelpers.WaitForCommit(t, e.Env.Chains[src1], destChain, oldState.Chains[dest].CommitStore[src1],
			sentEvent.Message.SequenceNumber)
		commitFound <- struct{}{}
	}()

	_, err = cldf.CreateLegacyChangeSet(v1_6.RMNUncurseChangeset).Apply(e.Env, v1_6.RMNCurseConfig{
		CurseActions: []v1_6.CurseAction{v1_6.CurseChain(e.Env.AllChainSelectors()[0])},
		Reason:       "Uncurse test",
	})
	require.NoError(t, err)

	for _, chainSel := range e.Env.AllChainSelectors() {
		subjects, err := state.Chains[chainSel].RMNRemote.GetCursedSubjects(nil)
		require.NoError(t, err)
		require.Empty(t, subjects)
	}

	// We have to restart all chainlink node because it cache the curse status for 30min
	tLocalEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	if !ok {
		t.Fatal("expected tEnv to be a DeployedLocalDevEnvironment")
	}
	err = tLocalEnv.RestartChainlinkNodes(t)
	require.NoError(t, err)

	timeUntilTimeout := 10 * time.Minute
	if deadline, ok := t.Context().Deadline(); ok {
		timeUntilTimeout = time.Until(deadline)
	}

	for {
		select {
		case <-commitFound:
			return
		case <-time.Tick(3 * time.Minute):
			err = tLocalEnv.RestartChainlinkNodes(t)
			require.NoError(t, err)
		case <-time.Tick(timeUntilTimeout):
			t.Fatal("timed out waiting for commit")
		}
	}
}

// TestMigrateFromV1_5ToV1_6 tests the migration from v1.5 to v1.6
func TestMigrateFromV1_5ToV1_6(t *testing.T) {
	t.Skip("Skipping since its flakey, need to fix")
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS) .
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithPrerequisiteDeploymentOnly(
			&changeset.V1_5DeploymentConfig{
				PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
				RMNConfig: &rmn_contract.RMNConfig{
					BlessWeightThreshold: 2,
					CurseWeightThreshold: 2,
					// setting dummy voters, we will permabless this later
					Voters: []rmn_contract.RMNVoter{
						{
							BlessWeight:   2,
							CurseWeight:   2,
							BlessVoteAddr: utils.RandomAddress(),
							CurseVoteAddr: utils.RandomAddress(),
						},
					},
				},
			}),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfUsersPerChain(2),
		// for in-memory test it is important to set the dest chain id as 1337 otherwise the config digest will not match
		// between nodes' calculated digest and the digest set on the contract
		testhelpers.WithChainIDs([]uint64{chainselectors.GETH_TESTNET.EvmChainID}),
	)
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChainsExcept1337 := e.Env.AllChainSelectorsExcluding([]uint64{chainselectors.GETH_TESTNET.Selector})
	require.Contains(t, e.Env.AllChainSelectors(), chainselectors.GETH_TESTNET.Selector)
	require.Len(t, allChainsExcept1337, 2)
	src1, src2, dest := allChainsExcept1337[0], allChainsExcept1337[1], chainselectors.GETH_TESTNET.Selector
	pairs := []testhelpers.SourceDestPair{
		// as mentioned in the comment above, the dest chain id should be 1337
		{SourceChainSelector: src1, DestChainSelector: dest},
		{SourceChainSelector: src2, DestChainSelector: dest},
	}
	// wire up all lanes
	// deploy onRamp, commit store, offramp , set ocr2config and send corresponding jobs
	e.Env = v1_5testhelpers.AddLanes(t, e.Env, state, pairs)

	// permabless the commit stores
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			v1_5.PermaBlessCommitStoreConfig{
				Configs: map[uint64]v1_5.PermaBlessCommitStoreConfigPerDest{
					dest: {
						Sources: []v1_5.PermaBlessConfigPerSourceChain{
							{
								SourceChainSelector: src1,
								PermaBless:          true,
							},
							{
								SourceChainSelector: src2,
								PermaBless:          true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)
	// ensure that all lanes are functional
	var (
		done                                = make(chan bool) // channel to stop sending messages in real router
		wg                                  sync.WaitGroup    // wait group to wait for all the messages to be delivered
		v1_5Msgs                            = make([]*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested, 0)
		v1_6Msgs                            = make([]*onramp.OnRampCCIPMessageSent, 0)
		initialBlock, lastNonce, firstNonce uint64
	)
	wg.Add(1)
	// send continuous messages in real router until done is closed
	go func() {
		defer wg.Done()
		initialBlock, v1_5Msgs, v1_6Msgs = sendContinuousMessages(
			t, &e, &state, pairs[0].SourceChainSelector, pairs[0].DestChainSelector, done)
	}()
	// send a message from the other lane src2 -> dest
	sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, state,
		testhelpers.WithSourceChain(src2),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, sentEvent)
	destChain := e.Env.Chains[dest]
	destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	v1_5testhelpers.WaitForCommit(t, e.Env.Chains[src2], destChain, state.Chains[dest].CommitStore[src2],
		sentEvent.Message.SequenceNumber)
	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src2], destChain, state.Chains[dest].EVM2EVMOffRamp[src2],
		[]uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())

	// now that all 1.5 lanes work transfer ownership of the contracts to MCMS
	contractsByChain := make(map[uint64][]common.Address)
	for _, chain := range e.Env.AllChainSelectors() {
		contractsByChain[chain] = []common.Address{
			state.Chains[chain].Router.Address(),
			state.Chains[chain].RMNProxy.Address(),
			state.Chains[chain].PriceRegistry.Address(),
			state.Chains[chain].TokenAdminRegistry.Address(),
			state.Chains[chain].RMN.Address(),
		}
		if state.Chains[chain].EVM2EVMOnRamp != nil {
			for _, onRamp := range state.Chains[chain].EVM2EVMOnRamp {
				contractsByChain[chain] = append(contractsByChain[chain], onRamp.Address())
			}
		}
		if state.Chains[chain].EVM2EVMOffRamp != nil {
			for _, offRamp := range state.Chains[chain].EVM2EVMOffRamp {
				contractsByChain[chain] = append(contractsByChain[chain], offRamp.Address())
			}
		}
	}

	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
			},
		),
	)
	require.NoError(t, err)
	// add 1.6 contracts to the environment and send 1.6 jobs
	// First we need to deploy Homechain contracts and restart the nodes with updated cap registry
	// in this test we have already deployed home chain contracts and the nodes are already running with the deployed cap registry.
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, false)
	// Set RMNProxy to point to RMNRemote.
	// nonce manager should point to 1.5 ramps
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			// as we have already transferred ownership for RMNProxy to MCMS, it needs to be done via MCMS proposal
			cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset),
			v1_6.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: e.Env.AllChainSelectors(),
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateNonceManagersChangeset),
			v1_6.UpdateNonceManagerConfig{
				// we only have lanes between src1 --> dest
				UpdatesByChain: map[uint64]v1_6.NonceManagerUpdate{
					src1: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
							{
								RemoteChainSelector: dest,
								AllowEmptyOffRamp:   true,
							},
						},
					},
					src2: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
							{
								RemoteChainSelector: dest,
								AllowEmptyOffRamp:   true,
							},
						},
					},
					dest: {
						PreviousRampsArgs: []v1_6.PreviousRampCfg{
							{
								RemoteChainSelector: src1,
								AllowEmptyOnRamp:    true,
							},
							{
								RemoteChainSelector: src2,
								AllowEmptyOnRamp:    true,
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Enable a single 1.6 lane with test router
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, src1, dest, true)
	require.GreaterOrEqual(t, len(e.Users[src1]), 2)
	testhelpers.ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)
	startBlocks := make(map[uint64]*uint64)
	latesthdr, err := e.Env.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[dest] = &block
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)
	expectedSeqNums := make(map[testhelpers.SourceDestPair]uint64)
	msgSentEvent, err := testhelpers.SendRequest(
		e.Env, state,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(true),
		// Send traffic across single 1.6 lane with a DIFFERENT ( very important to not mess with real sender nonce) sender
		// from test router to ensure 1.6 is working.
		testhelpers.WithSender(e.Users[src1][1]),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}))
	require.NoError(t, err)

	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: src1,
		DestChainSelector:   dest,
	}] = []uint64{msgSentEvent.SequenceNumber}
	expectedSeqNums[testhelpers.SourceDestPair{
		SourceChainSelector: src1,
		DestChainSelector:   dest,
	}] = msgSentEvent.SequenceNumber

	// This sleep is needed so that plugins come up and start indexing logs.
	// Otherwise test will flake.
	time.Sleep(30 * time.Second)
	testhelpers.ReplayLogs(t, e.Env.Offchain, map[uint64]uint64{
		src1: msgSentEvent.Raw.BlockNumber,
	})
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNums, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)

	// now that the 1.6 lane is working, we can enable the real router
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					src1: {
						dest: {
							IsEnabled:        true,
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
			v1_6.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
					dest: {
						src1: {
							IsEnabled:                 true,
							TestRouter:                false,
							IsRMNVerificationDisabled: true,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			// this needs to be MCMS proposal as the router contract is owned by MCMS
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: false,
				MCMS: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// onRamp update on source chain
					src1: {
						OnRampUpdates: map[uint64]bool{
							dest: true,
						},
					},
					// offramp update on dest chain
					dest: {
						OffRampUpdates: map[uint64]bool{
							src1: true,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)
	// confirm that the other lane src2->dest is still working with v1.5
	sentEventOnOtherLane, err := v1_5testhelpers.SendRequest(t, e.Env, state,
		testhelpers.WithSourceChain(src2),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, sentEvent)

	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src2], e.Env.Chains[dest], state.Chains[dest].EVM2EVMOffRamp[src2],
		[]uint64{sentEventOnOtherLane.Message.SequenceNumber}, destStartBlock.Number.Uint64())

	// stop the continuous messages in real router
	close(done) // stop sending messages in real router
	wg.Wait()
	// start validating the messages sent in 1.5 and 1.6
	for _, msg := range v1_5Msgs {
		v1_5testhelpers.WaitForCommit(t, e.Env.Chains[src1], destChain, state.Chains[dest].CommitStore[src1],
			msg.Message.SequenceNumber)
		v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src1], destChain, state.Chains[dest].EVM2EVMOffRamp[src1],
			[]uint64{msg.Message.SequenceNumber}, initialBlock)
		lastNonce = msg.Message.Nonce
	}
	for _, msg := range v1_6Msgs {
		if firstNonce == 0 {
			firstNonce = msg.Message.Header.Nonce
		}
		expectedSeqNumExec[testhelpers.SourceDestPair{
			SourceChainSelector: src1,
			DestChainSelector:   dest,
		}] = []uint64{msg.Message.Header.SequenceNumber}
		expectedSeqNums[testhelpers.SourceDestPair{
			SourceChainSelector: src1,
			DestChainSelector:   dest,
		}] = msg.Message.Header.SequenceNumber
	}
	startBlocks[dest] = &initialBlock
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNums, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
	// this seems to be flakey, also might be incorrect?
	require.Equal(t, lastNonce+1, firstNonce, "sender nonce in 1.6 OnRamp event is not plus one to sender nonce in 1.5 OnRamp")
}

// SendMessages sends messages from src to dest until done is closed
func sendContinuousMessages(
	t *testing.T,
	e *testhelpers.DeployedEnv,
	state *changeset.CCIPOnChainState,
	src, dest uint64,
	done chan bool,
) (uint64, []*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested, []*onramp.OnRampCCIPMessageSent) {
	var (
		// TODO: make this shorter than 10 seconds, maybe 2 seconds?
		ticker           = time.NewTicker(10 * time.Second)
		initialDestBlock uint64
		v1_6Msgs         []*onramp.OnRampCCIPMessageSent
		v1_5Msgs         []*evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested
	)
	for {
		select {
		case <-ticker.C:
			msg := sendMessageInRealRouter(t, e, state, src, dest)
			if msg == nil {
				t.Errorf("failed to send message in real router")
				continue
			}
			switch msg := msg.(type) {
			case *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested:
				v1_5Msgs = append(v1_5Msgs, msg)
				if initialDestBlock == 0 {
					destChain := e.Env.Chains[dest]
					destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
					if err != nil {
						t.Errorf("failed to get block header")
					}
					initialDestBlock = destStartBlock.Number.Uint64()
				}
			case *onramp.OnRampCCIPMessageSent:
				v1_6Msgs = append(v1_6Msgs, msg)
			}
		case <-done:
			return initialDestBlock, v1_5Msgs, v1_6Msgs
		}
	}
}

// sendMessageInRealRouter sends a message and filter the log topic to identify the event type nad parse the event data.
func sendMessageInRealRouter(
	t *testing.T,
	e *testhelpers.DeployedEnv,
	state *changeset.CCIPOnChainState,
	src, dest uint64,
) any {
	cfg := &testhelpers.CCIPSendReqConfig{
		SourceChain:  src,
		DestChain:    dest,
		Sender:       e.Env.Chains[src].DeployerKey,
		IsTestRouter: false,
		Message: router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		},
	}
	t.Logf("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())

	tx, _, err := testhelpers.CCIPSendRequest(e.Env, *state, cfg)
	if err != nil {
		t.Errorf("failed to send message: %v", err)
	}
	receipt, err := e.Env.Chains[src].Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Errorf("failed to get transaction receipt: %v", err)
	}
	// filter the log topic to identify the event type and parse the event data
	for _, lg := range receipt.Logs {
		if lg.Topics[0].Hex() == evm2EVMOnRampABI.Events["CCIPSendRequested"].ID.Hex() {
			unpackedMsg, err := evm2EVMOnRampABI.Events["CCIPSendRequested"].Inputs.Unpack(lg.Data)
			if err != nil {
				t.Errorf("failed to unpack ccip send requested event")
			}
			return &evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested{
				Message: *abi.ConvertType(unpackedMsg[0], new(evm_2_evm_onramp.InternalEVM2EVMMessage)).(*evm_2_evm_onramp.InternalEVM2EVMMessage),
			}
		} else if lg.Topics[0].Hex() == onRampABI.Events["CCIPMessageSent"].ID.Hex() {
			unpackedMsg, err := onRampABI.Events["CCIPMessageSent"].Inputs.Unpack(lg.Data)
			if err != nil {
				t.Errorf("failed to unpack ccip message sent event")
			}
			return &onramp.OnRampCCIPMessageSent{
				Message: *abi.ConvertType(unpackedMsg[0], new(onramp.InternalEVM2AnyRampMessage)).(*onramp.InternalEVM2AnyRampMessage),
			}
		}
	}
	return nil
}
