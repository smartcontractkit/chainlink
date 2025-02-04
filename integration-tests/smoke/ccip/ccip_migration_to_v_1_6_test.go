package ccip

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	v1_5testhelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/evm/utils"
)

// TestMigrateFromV1_5ToV1_6 tests the migration from v1.5 to v1.6
// The test steps are as follows:
// 1. Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
// 2. Deploy 1.5 contracts (excluding pools to start, but including MCMS)
// 3. Wire up all lanes
// 4. PermaBless the commit stores
// 5. Send continuous messages from src1 -> dest in real router until stopMsgs is closed
// 6. Send a message from the other lane src2 -> dest
// 6. Transfer ownership of the contracts to MCMS
// 7. Add 1.6 contracts to the environment and send 1.6 jobs
// 8. Set RMNProxy to point to RMNRemote
// 9. Update nonce managers
// 10. Enable a single 1.6 lane with test router
// 11. Send traffic across single 1.6 src1 > dest lane with a DIFFERENT sender from test router to ensure 1.6 is working
// 12. Enable the real router and verify sender nonce in 1.6 OnRamp event is plus one to sender nonce in 1.5 OnRamp
// 13. Confirm that the other lane src2->dest is still working with v1.5
// 14. Stop the continuous messages in real router and wait for all the executions to confirm
func TestMigrateFromV1_5ToV1_6(t *testing.T) {
	// Deploy CCIP 1.5 with 3 chains and 4 nodes + 1 bootstrap
	// Deploy 1.5 contracts (excluding pools to start, but including MCMS)
	const msgInterval = 10 * time.Second
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
	stopMsgs := make(chan bool)     // channel to stop sending messages in real router
	switchTov1_6 := make(chan bool) // channel to switchTov1_6 to switch to 1.6
	var wg sync.WaitGroup           // wait group to wait for all the messages to be delivered
	// send continuous messages from src1 -> dest until stopMsgs is closed
	sendContinuousMessagesInRealRouter(t, &e, &state, pairs[0], &wg, stopMsgs, switchTov1_6, msgInterval)
	// send a message from the other lane src2 -> dest
	sentEvent := sendMsgInV1_5(t, e, state, src2, dest)
	destChain := e.Env.Chains[dest]
	destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	v1_5testhelpers.WaitForCommit(t, e.Env.Chains[src2], destChain, state.Chains[dest].CommitStore[src1],
		sentEvent.Message.SequenceNumber)
	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src2], destChain, state.Chains[dest].EVM2EVMOffRamp[src1],
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
	e = testhelpers.AddCCIPContractsToEnvironment(t, e.Env.AllChainSelectors(), tEnv, false)
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
	// Send traffic across single 1.6 lane with a DIFFERENT ( very important to not mess with real sender nonce) sender
	// from test router to ensure 1.6 is working.
	msgSentEvent := sendMsgInV1_6(t, e, state, src1, dest, e.Users[src1][1], true)
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

	// As the real router are getting continuous messages, after this switch v1.6 OnRamp should start receiving the messages.
	// and the request should get delivered to 1.6 offRamp
	switchTov1_6 <- true // Signal to look for messages in 1.6 OnRamp

	// confirm that the other lane src2->dest is still working with v1.5
	sentEventOnOtherLane := sendMsgInV1_5(t, e, state, src2, dest)
	v1_5testhelpers.WaitForExecute(t, e.Env.Chains[src2], e.Env.Chains[dest], state.Chains[dest].EVM2EVMOffRamp[src2],
		[]uint64{sentEventOnOtherLane.Message.SequenceNumber}, destStartBlock.Number.Uint64())

	// stop the continuous messages in real router and wait for all the executions to confirm
	stopMsgs <- true
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	},
		tests.WaitTimeout(t),
		2*time.Second,
		"all executions did not confirm",
	)
}

type msgEvent struct {
	tx       common.Hash
	blockNum uint64
}

// sendContinuousMessagesInRealRouter sends continuous messages from the source chain to the destination chain until stopped.
// The process is as follows:
//  1. Send a message once and then send messages based on the interval in real router until signal to stop the message is received.
//  2. Listen to the message sent event in 1.5 onRamp until the switch signal is received
//     and send the message sequence number and block number to check for commit & execute through performMsgCheckInV1_5 channel.
//  3. Separate Go routines listens to the channel performMsgCheckInV1_5 and performMsgCheckInV1_6 to perform
//     commit and execute in 1.5 and 1.6 version offRamps respectively.
//  4. Start listening to the message sent event in 1.6 onRamp after the switch signal is received and send the message
//     sequence number and block number to check for commit & execute through performMsgCheckInV1_6 channel.
//  5. In the test, wait for all the go routines to complete to confirm all the messages are delivered.
//  6. Assert sender nonce in 1.6 OnRamp event is plus one to sender nonce in 1.5 OnRamp .
func sendContinuousMessagesInRealRouter(
	t *testing.T,
	e *testhelpers.DeployedEnv,
	state *changeset.CCIPOnChainState,
	pair testhelpers.SourceDestPair,
	wg *sync.WaitGroup,
	stopMsgs, signal chan bool,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	type msgSentCheck struct {
		src, dest, seqNr, startBlockNumber uint64
	}
	msgPipeline := make(chan msgEvent, 100) // channel to send the message tx and block number to check msg sent event
	// channel to send the message sent event to check for commit & execute in
	performMsgCheckInV1_5 := make(chan msgSentCheck, 100)
	// channel to send the message sent event to check for commit & execute
	performMsgCheckInV1_6 := make(chan msgSentCheck, 100)
	performNonceCheck := false

	// Go routine to send the messages continuously until stopMsgs is closed
	go func() {
		defer ticker.Stop()
		// send once and further send messages based on the interval
		// the below function sends the message and adds the message sent event to the msgPipeline
		sendMessageWithoutCapturingSentEvent(t, e, state, pair.SourceChainSelector, pair.DestChainSelector, msgPipeline)
		for {
			select {
			case <-stopMsgs: // stop sending messages and close the msgPipeline
				close(msgPipeline)
				return
			case <-ticker.C:
				sendMessageWithoutCapturingSentEvent(t, e, state, pair.SourceChainSelector, pair.DestChainSelector, msgPipeline)
			}
		}
	}()

	// function to listen to the message sent event in v1.6 onRamp
	listenV1_6OnRamp := func(input msgEvent) {
		defer wg.Done()
		// filter the message sent event in 1.6 onRamp
		it, err := state.Chains[pair.SourceChainSelector].OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
			Start:   input.blockNum,
			End:     &input.blockNum,
			Context: context.Background(),
		}, []uint64{pair.DestChainSelector}, []uint64{})
		if err != nil {
			t.Errorf("failed to get msg sent filter iterator")
		}
		if !it.Next() {
			t.Errorf("failed to get next event")
		}
		t.Logf("CCIP message (id %x) sent from chain selector %d to chain selector %d tx %s seqNum %d sender %s",
			it.Event.Message.Header.MessageId[:],
			pair.SourceChainSelector,
			pair.DestChainSelector,
			input.tx.String(),
			it.Event.Message.Header.SequenceNumber,
			it.Event.Message.Sender.String(),
		)
		destChain := e.Env.Chains[pair.DestChainSelector]
		destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Errorf("failed to get header by number")
		}
		if performNonceCheck {
			lastnonce, err := state.Chains[pair.SourceChainSelector].EVM2EVMOnRamp[pair.DestChainSelector].GetSenderNonce(
				nil,
				e.Env.Chains[pair.SourceChainSelector].DeployerKey.From,
			)
			if err != nil {
				t.Errorf("failed to get sender nonce in 1.5 OnRamp")
			}
			// assert sender nonce in 1.6 OnRamp event is plus one to sender nonce in 1.5 OnRamp
			if lastnonce+1 != it.Event.Message.Header.Nonce {
				t.Errorf("sender nonce in 1.6 OnRamp event is not plus one to sender nonce in 1.5 OnRamp")
			}
			performNonceCheck = false
		}
		// channel to send the message sequence number and block number to check for commit & execute
		performMsgCheckInV1_6 <- msgSentCheck{
			src:              pair.SourceChainSelector,
			dest:             pair.DestChainSelector,
			seqNr:            it.Event.Message.Header.SequenceNumber,
			startBlockNumber: destStartBlock.Number.Uint64(),
		}
	}
	// function to listen to the message sent event in v1.5 onRamp
	listenV1_5OnRamp := func(input msgEvent) {
		defer wg.Done()
		// filter the message sent event in 1.5 onRamp
		it, err := state.Chains[pair.SourceChainSelector].EVM2EVMOnRamp[pair.DestChainSelector].FilterCCIPSendRequested(&bind.FilterOpts{
			Start:   input.blockNum,
			End:     &input.blockNum,
			Context: context.Background(),
		})
		if err != nil {
			t.Errorf("failed to get msg sent filter iterator")
		}
		if !it.Next() {
			t.Errorf("failed to get next event")
		}

		t.Logf("CCIP message (id %x) sent from chain selector %d to chain selector %d tx %s seqNum %d sender %s",
			it.Event.Message.MessageId[:],
			pair.SourceChainSelector,
			pair.DestChainSelector,
			input.tx.String(),
			it.Event.Message.SequenceNumber,
			it.Event.Message.Sender.String(),
		)
		destChain := e.Env.Chains[pair.DestChainSelector]
		destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Errorf("failed to get header by number")
		}
		// channel to send the message sequence number and block number to check for commit & execute
		performMsgCheckInV1_5 <- msgSentCheck{
			src:              pair.SourceChainSelector,
			dest:             pair.DestChainSelector,
			seqNr:            it.Event.Message.SequenceNumber,
			startBlockNumber: destStartBlock.Number.Uint64(),
		}
	}

	// Process the msg event in real router
	go func() {
		// switchSignal is used to switch to 1.6 onRamp after the signal is received until then it listens to 1.5 onRamp
		switchSignal := false
		defer close(performMsgCheckInV1_5)
		defer close(performMsgCheckInV1_6)
		defer close(signal)
		for input := range msgPipeline {
			wg.Add(1)
			select {
			case <-signal:
				// switch to 1.6 onRamp
				switchSignal = true
				performNonceCheck = true
				go listenV1_6OnRamp(input)
			default:
				if switchSignal {
					go listenV1_6OnRamp(input)
				} else {
					// listen to 1.5 onRamp until switchSignal is set to true
					go listenV1_5OnRamp(input)
				}
			}
		}
	}()

	// Verify message sent in v1.5
	go func() {
		// range over performMsgCheckInV1_5 channel to verify the commit and execute in 1.5 offRamp
		// this channel will be populated when there is msg sent event in 1.5 onRamp
		for input := range performMsgCheckInV1_5 {
			wg.Add(1)
			go func(input msgSentCheck) {
				defer wg.Done()
				v1_5testhelpers.WaitForCommit(t, e.Env.Chains[input.src], e.Env.Chains[input.dest],
					state.Chains[input.dest].CommitStore[input.src], input.seqNr)
				v1_5testhelpers.WaitForExecute(t, e.Env.Chains[input.src], e.Env.Chains[input.dest],
					state.Chains[pair.DestChainSelector].EVM2EVMOffRamp[pair.SourceChainSelector],
					[]uint64{input.seqNr}, input.startBlockNumber)
			}(input)
		}
	}()

	// verify message sent in v1.6
	go func() {
		// range over performMsgCheckInV1_6 channel to verify the commit and execute in 1.6 offRamp
		// this channel will be populated when there is msg sent event in 1.6 onRamp
		for input := range performMsgCheckInV1_6 {
			wg.Add(1)
			go func(input msgSentCheck) {
				defer wg.Done()
				// confirm commit in 1.6 offRamp
				_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
					t,
					input.src,
					e.Env.Chains[input.dest],
					state.Chains[input.dest].OffRamp,
					&input.startBlockNumber,
					ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(input.seqNr),
						ccipocr3.SeqNum(input.seqNr),
					},
					true,
				)
				if err != nil {
					t.Errorf("failed to confirm commit")
				}
				// confirm execute in 1.6 offRamp
				_, err = testhelpers.ConfirmExecWithSeqNrs(
					t,
					input.src,
					e.Env.Chains[input.dest],
					state.Chains[input.dest].OffRamp,
					&input.startBlockNumber,
					[]uint64{input.seqNr},
				)
				if err != nil {
					t.Errorf("failed to confirm execute")
				}
			}(input)
		}
	}()
}

// sendMessageWithoutCapturingSentEvent sends a message from the source chain to the destination chain without capturing the message sent event.
func sendMessageWithoutCapturingSentEvent(
	t *testing.T,
	e *testhelpers.DeployedEnv,
	state *changeset.CCIPOnChainState,
	src, dest uint64,
	msgPipeline chan msgEvent,
) {
	cfg := &testhelpers.CCIPSendReqConfig{
		SourceChain:  src,
		DestChain:    dest,
		Sender:       e.Env.Chains[src].DeployerKey,
		IsTestRouter: false,
		Evm2AnyMessage: router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		},
	}
	t.Logf("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())

	tx, blockNum, err := testhelpers.CCIPSendRequest(e.Env, *state, cfg)
	if err != nil {
		t.Errorf("failed to send message: %v", err)
	}
	// send the message tx and block number to msgPipeline channel which looks for the message sent event in the onRamp.
	// initially it will look for the message sent event in 1.5 onRamp and
	// after switchSignal is set to true, it will look for the message sent event in 1.6 onRamp
	msgPipeline <- msgEvent{
		tx:       tx.Hash(),
		blockNum: blockNum,
	}
}

// sendMsgInV1_5 sends a message and filter the message sent event in 1.5 onRamp.
func sendMsgInV1_5(
	t *testing.T,
	e testhelpers.DeployedEnv,
	state changeset.CCIPOnChainState,
	src, dest uint64) *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested {
	sentEvent, err := v1_5testhelpers.SendRequest(t, e.Env, state,
		testhelpers.WithSourceChain(src),
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
	return sentEvent
}

// sendMsgInV1_6 sends a message and filter the message sent event in 1.6 onRamp.
func sendMsgInV1_6(
	t *testing.T,
	e testhelpers.DeployedEnv,
	state changeset.CCIPOnChainState,
	src, dest uint64,
	sender *bind.TransactOpts,
	isTestRouter bool,
) *onramp.OnRampCCIPMessageSent {
	sentEvent, err := testhelpers.DoSendRequest(
		t, e.Env, state,
		testhelpers.WithSourceChain(src),
		testhelpers.WithDestChain(dest),
		testhelpers.WithTestRouter(isTestRouter),
		testhelpers.WithSender(sender),
		testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}))
	require.NoError(t, err)
	return sentEvent
}
