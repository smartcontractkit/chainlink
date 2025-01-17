package ccip

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	v1_5testhelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestMigrateFromV1_5ToV1_6(t *testing.T) {
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
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(v1_5.PermaBlessCommitStoreChangeset),
			Config: v1_5.PermaBlessCommitStoreConfig{
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
		},
	})
	require.NoError(t, err)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)
	// ensure that all lanes are functional
	for _, pair := range pairs {
		sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, state,
			testhelpers.WithSourceChain(pair.SourceChainSelector),
			testhelpers.WithDestChain(pair.DestChainSelector),
			testhelpers.WithTestRouter(false),
			testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[pair.DestChainSelector].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, sentEvent)
		destChain := e.Env.Chains[pair.DestChainSelector]
		destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
		require.NoError(t, err)
		v1_5testhelpers.WaitForCommit(t, e.Env.Chains[pair.SourceChainSelector], destChain, state.Chains[dest].CommitStore[src1], sentEvent.Message.SequenceNumber)
		v1_5testhelpers.WaitForExecute(t, e.Env.Chains[pair.SourceChainSelector], destChain, state.Chains[dest].EVM2EVMOffRamp[src1], []uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())
	}

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

	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config: commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MinDelay:         0,
			},
		},
	})
	require.NoError(t, err)
	// add 1.6 contracts to the environment and send 1.6 jobs
	// First we need to deploy Homechain contracts and restart the nodes with updated cap registry
	// in this test we have already deployed home chain contracts and the nodes are already running with the deployed cap registry.
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, true, true, false)
	// Set RMNProxy to point to RMNRemote.
	// nonce manager should point to 1.5 ramps
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			// as we have already transferred ownership for RMNProxy to MCMS, it needs to be done via MCMS proposal
			Changeset: commonchangeset.WrapChangeSet(changeset.SetRMNRemoteOnRMNProxyChangeset),
			Config: changeset.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: e.Env.AllChainSelectors(),
				MCMSConfig: &changeset.MCMSConfig{
					MinDelay: 0,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateNonceManagersChangeset),
			Config: changeset.UpdateNonceManagerConfig{
				// we only have lanes between src1 --> dest
				UpdatesByChain: map[uint64]changeset.NonceManagerUpdate{
					src1: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: dest,
								EnableOnRamp:        true,
							},
						},
					},
					src2: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: dest,
								EnableOnRamp:        true,
							},
						},
					},
					dest: {
						PreviousRampsArgs: []changeset.PreviousRampCfg{
							{
								RemoteChainSelector: src1,
								EnableOffRamp:       true,
							},
							{
								RemoteChainSelector: src2,
								EnableOffRamp:       true,
							},
						},
					},
				},
			},
		},
	})
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
	msgSentEvent, err := testhelpers.DoSendRequest(
		t, e.Env, state,
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
	time.Sleep(15 * time.Second)
	testhelpers.ReplayLogs(t, e.Env.Offchain, map[uint64]uint64{
		src1: msgSentEvent.Raw.BlockNumber,
	})
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNums, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)

	// send a message from real router, the send requested event should be received in 1.5 onRamp
	// the request should get delivered to 1.5 offRamp
	sentEventBeforeSwitch, err := v1_5testhelpers.SendRequest(t, e.Env, state,
		testhelpers.WithSourceChain(src1),
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
	require.NotNil(t, sentEventBeforeSwitch)

	// now that the 1.6 lane is working, we can enable the real router
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
			Config: changeset.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
					src1: {
						dest: {
							IsEnabled:        true,
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOffRampSourcesChangeset),
			Config: changeset.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
					dest: {
						src1: {
							IsEnabled:  true,
							TestRouter: false,
						},
					},
				},
			},
		},
		{
			// this needs to be MCMS proposal as the router contract is owned by MCMS
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateRouterRampsChangeset),
			Config: changeset.UpdateRouterRampsConfig{
				TestRouter: false,
				MCMS: &changeset.MCMSConfig{
					MinDelay: 0,
				},
				UpdatesByChain: map[uint64]changeset.RouterUpdates{
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
		},
	})
	require.NoError(t, err)

	// send a message from real router the send requested event should be received in 1.6 onRamp
	// the request should get delivered to 1.6 offRamp
	destStartBlock, err := e.Env.Chains[dest].Client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	sentEventAfterSwitch, err := testhelpers.DoSendRequest(
		t, e.Env, state,
		testhelpers.WithSourceChain(src1),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(false),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}))
	require.NoError(t, err)
	// verify that before switch message is received in 1.5 offRamp
	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src1], e.Env.Chains[dest], state.Chains[dest].EVM2EVMOffRamp[src1],
		[]uint64{sentEventBeforeSwitch.Message.SequenceNumber}, destStartBlock.Number.Uint64())

	// verify that after switch message is received in 1.6 offRamp
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: src1,
		DestChainSelector:   dest,
	}] = []uint64{sentEventAfterSwitch.SequenceNumber}
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)

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
	require.NotNil(t, sentEventOnOtherLane)
	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src2], e.Env.Chains[dest], state.Chains[dest].EVM2EVMOffRamp[src2],
		[]uint64{sentEventOnOtherLane.Message.SequenceNumber}, destStartBlock.Number.Uint64())
}
