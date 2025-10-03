package ccip

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	tonlib "github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink-ton/pkg/bindings"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/tracetracking"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/wrappers"

	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_TON2EVM(t *testing.T) {
	t.Skip("Currently skipping TON2EVM, Debugging EVM2TON")
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithTonChains(1))

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Logf("Loaded state: %v", state)
	_ = state

	// make evm chains sorted for deterministic test results
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)

	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := allTonChainSelectors[0]
	destChain := evmChainSelectors[0]
	t.Log("TON chain selectors:", allTonChainSelectors,
		", EVM chain selectors:", evmChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	tonChain := e.Env.BlockChains.TonChains()[sourceChain]
	ac := codec.NewAddressCodec()
	addrBytes, err := ac.AddressStringToBytes(tonChain.WalletAddress.String())
	require.NoError(t, err)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		sender = addrBytes
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		receiver := common.LeftPadBytes(e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), 32)
		require.NoError(t, err)

		ea := onramp.GenericExtraArgsV2{
			GasLimit:                 big.NewInt(1000000),
			AllowOutOfOrderExecution: true,
		}
		c, err := tlb.ToCell(ea)
		require.NoError(t, err)
		out = mt.Run(
			t,
			mt.TestCase{
				Replayed:               true,
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiver,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              c.ToBOC(),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	_ = out
}

func Test_CCIPMessaging_EVM2TON(t *testing.T) {
	//t.Skip("Test stalls because TON test assertions aren't implemented yet")
	// Setup 2 chains (EVM and Ton) and a single lane.
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]

	t.Logf("=== Test Configuration ===")
	t.Logf("  Source (EVM): %d", sourceChain)
	t.Logf("  Dest (TON):   %d", destChain)
	t.Logf("  OnRamp:       %s", state.Chains[sourceChain].OnRamp.Address())

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	t.Run("message to contract receiver", func(t *testing.T) {
		tonChain := e.Env.BlockChains.TonChains()[destChain]
		offRampAddr := state.TonChains[destChain].OffRamp

		receiver, err := deployReceiverContract(tonChain, &offRampAddr)
		require.NoError(t, err)

		t.Logf("  OffRamp:  %s", offRampAddr.String())
		t.Logf("  Receiver: %s", receiver.String())

		// TODO: should receiver address be saved in state?
		ccipChainState := state.TonChains[destChain]
		ccipChainState.ReceiverAddress = *receiver
		state.TonChains[destChain] = ccipChainState

		ac := codec.NewAddressCodec()
		receiverBytes, err := ac.AddressStringToBytes(receiver.String())
		require.NoError(t, err)
		require.Equal(t, 36, len(receiverBytes), "receiver bytes should be 36 bytes")

		// Subscribe to OffRamp contract transactions in background
		SubscribeOfframpTransactions(t, tonChain, offRampAddr)

		// TODO: receiver.tolk's onInternalMessage asserts that msg.message.data.beginParse().isEmpty()
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType: mt.ValidationTypeExec,
				TestSetup:      setup,
				Nonce:          nil, // TON nonce check is skipped
				Receiver:       receiverBytes,
				MsgData:        []byte{}, // TODO: empty data fails?
				// MsgData:                []byte("hello CCIPReceiver"), // TODO: empty data fails?
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
		// TODO: need a test case with wallet receiver(no reply nor received events)
	})

	_ = out
}

// TODO: do we want to have a changeset for receiver? probably for staging validation
func deployReceiverContract(tonChain ton.Chain, offRampAddr *address.Address) (*address.Address, error) {
	// parse compiled contract
	codeCell, err := wrappers.ParseCompiledContract(bindings.GetBuildDir("Receiver.compiled.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Receiver compiled contract: %w", err)
	}

	// Create initial storage - must match TypeScript: beginCell().storeAddress(offRampAddress).endCell()
	// Note: Unlike other contracts that use tlb.ToCell(storage) which creates empty root + ref structure,
	// receiver.tolk expects a simple cell with address stored directly in root cell.
	receiverStorage := cell.BeginCell().
		MustStoreAddr(offRampAddr).
		EndCell()

	conn := tracetracking.NewSignedAPIClient(tonChain.Client, *tonChain.Wallet)
	contract, _, err := wrappers.Deploy(
		&conn,
		codeCell,
		receiverStorage,
		tlb.MustFromTON("5"), // TODO: Configurable
		cell.BeginCell().EndCell(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Receiver contract: %w", err)
	}
	receiver := contract.Address

	return receiver, nil
}

func SubscribeOfframpTransactions(t *testing.T, tonChain ton.Chain, offRampAddr address.Address) {
	t.Logf("Subscribing to OffRamp transactions...")
	ctx := t.Context() // Use test context to keep subscription alive until test ends

	master, err := tonChain.Client.CurrentMasterchainInfo(ctx)
	require.NoError(t, err)

	acc, err := tonChain.Client.GetAccount(ctx, master, &offRampAddr)
	require.NoError(t, err)
	lastProcessedLT := acc.LastTxLT

	// Create channel for transactions
	transactions := make(chan *tlb.Transaction)

	// Start subscription in background
	go tonChain.Client.SubscribeOnTransactions(ctx, &offRampAddr, lastProcessedLT, transactions)

	// Process transactions in background
	go func() {
		transactionCount := 0
		heartbeat := time.NewTicker(30 * time.Second)
		defer heartbeat.Stop()

		for {
			select {
			case tx := <-transactions:
				transactionCount++
				t.Logf("=== TX #%d (LT: %d) ===", transactionCount, tx.LT)

				// Parse and log input message type
				if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
					if internal := tx.IO.In.AsInternal(); internal != nil && internal.Body != nil {
						parser := internal.Body.BeginParse()
						var commitMsg offramp.Commit
						if err := tlb.LoadFromCell(&commitMsg, parser); err == nil {
							t.Logf("  INPUT: Commit message")
						} else {
							parser2 := internal.Body.BeginParse()
							var execMsg offramp.Execute
							if err2 := tlb.LoadFromCell(&execMsg, parser2); err2 == nil {
								t.Logf("  INPUT: Execute message")
							}
						}
					}
				}

				// Count and summarize output messages
				if tx.IO.Out != nil {
					if msgs, err := tx.IO.Out.ToSlice(); err == nil {
						externalCount := 0
						for _, msg := range msgs {
							if msg.MsgType == tlb.MsgTypeExternalOut {
								externalCount++
							}
						}
						if externalCount > 0 {
							t.Logf("  OUTPUT: %d messages (%d external events)", len(msgs), externalCount)
						}
					}
				}

				receivedMessage, err := tracetracking.MapToReceivedMessage(tx)
				if err != nil {
					t.Logf("  ERROR: Failed to map to ReceivedMessage: %v", err)
					continue
				}

				err = receivedMessage.WaitForTrace(tonChain.Client)
				if err != nil {
					t.Logf("  ERROR: Failed to wait for trace: %v", err)
					continue
				}

				if receivedMessage.ExitCode != 0 {
					t.Logf("  ERROR: Exit code %d - %s", receivedMessage.ExitCode, receivedMessage.ExitCode.Describe())
					continue
				}

				lm, lmerr := waitForReceivedMsgFlatten(t, tonChain.Client, &receivedMessage)
				if lmerr != nil {
					t.Logf("  ERROR: Failed to flatten messages: %v", lmerr)
					continue
				}

				// Only log summary of flattened messages
				if len(lm.OutgoingExternalMessages) > 0 {
					t.Logf("  RESULT: Success (%d internal, %d external)",
						len(lm.OutgoingInternalReceivedMessages), len(lm.OutgoingExternalMessages))
				}

			case <-heartbeat.C:
				t.Logf("Monitoring OffRamp (%d transactions)", transactionCount)

			case <-ctx.Done():
				t.Logf("Subscription ended (%d transactions)", transactionCount)
				return
			}
		}
	}()
}

func waitForReceivedMsgFlatten(t *testing.T, clientConn tonlib.APIClientWrapped, msg *tracetracking.ReceivedMessage) (*tracetracking.ReceivedMessage, error) {
	if msg == nil {
		return nil, errors.New("received message is nil")
	}

	// Collect all messages to process in a queue
	var messagesToProcess []*tracetracking.ReceivedMessage
	messagesToProcess = append(messagesToProcess, msg)

	var lastMsg *tracetracking.ReceivedMessage

	// Process messages iteratively
	for len(messagesToProcess) > 0 {
		// Get the first message from the queue
		currentMsg := messagesToProcess[0]
		messagesToProcess = messagesToProcess[1:]

		if len(currentMsg.OutgoingInternalReceivedMessages) == 0 {
			continue
		}

		for i, outMsg := range currentMsg.OutgoingInternalReceivedMessages {
			// Only log errors and bounces
			if outMsg.ExitCode != 0 {
				t.Logf("    Message %d FAILED: exit code %v - %v", i, outMsg.ExitCode, outMsg.ExitCode.Describe())
			}
			if outMsg.EmittedBouncedMessage {
				t.Logf("    Message %d BOUNCED", i)
			}

			err := outMsg.WaitForTrace(clientConn)
			if err != nil {
				t.Logf("    Message %d: failed to wait for trace: %v", i, err)
				continue
			}

			// Add this message to the queue for further processing
			messagesToProcess = append(messagesToProcess, outMsg)
			lastMsg = outMsg
		}
	}

	if lastMsg == nil {
		return nil, errors.New("no received messages were processed")
	}

	// var event any
	// err := tlb.LoadFromCell(&event, lastMsg.OutgoingExternalMessages[0].Body.BeginParse())
	// if err != nil {
	// 	t.Logf("failed to parse CCIPMessageSent from cell: %v", err)
	// 	return nil, err
	// }

	return lastMsg, nil
}
