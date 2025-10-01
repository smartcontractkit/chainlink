package ccip

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	tonlib "github.com/xssnick/tonutils-go/ton"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	test_utils "github.com/smartcontractkit/chainlink-ton/deployment/utils"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/tracetracking"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

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

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Logf("Loaded state: %v", state)
	_ = state

	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]
	t.Log("EVM chain selectors:", evmChainSelectors,
		", TON chain selectors:", allTonChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// -- log TON configurations
	// offRampAddr := state.TonChains[destChain].OffRamp
	// tonChain := e.Env.BlockChains.TonChains()[destChain]

	// 	master, err := tonChain.Client.CurrentMasterchainInfo(t.Context())
	// require.NoError(t, err)

	// result, err := tonChain.Client.RunGetMethod(t.Context(), master, &offRampAddr, "sourceChainConfig", uint64(sourceChain))
	// require.NoError(t, err)
	// var cfg offramp.SourceChainConfig
	// err = cfg.Fr(&cfg, result.)
	// require.NoError(t, err)
	// t.Logf("OffRamp source chain config for sourceChain %d: %+v", sourceChain, cfg)

	// -- log evm onramp
	evmOnramp := state.Chains[sourceChain].OnRamp
	t.Logf("EVM onramp: %+v", evmOnramp.Address())

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

	t.Logf("waiting for filter registration for CCIPMessageSent (onramp), CommitReportAccepted (offramp), and ExecutionStateChanged (offramp), usually takes less than 2 mins")
	// wait for filter registration for CCIPMessageSent (onramp), CommitReportAccepted (offramp), and ExecutionStateChanged (offramp)
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		ccipChainState := state.TonChains[destChain]
		receiver := ccipChainState.ReceiverAddress
		receiverBase64Bytes, err := base64.RawURLEncoding.DecodeString(receiver.String())
		require.NoError(t, err)
		// Prepare 36-byte raw address
		receiver.FlagsToByte()

		t.Logf("Receiver address: %s", receiver.String())
		t.Logf("Receiver base64 bytes: %s", receiverBase64Bytes)

		// activate receiver
		test_utils.FundWallets(t, e.Env.BlockChains.TonChains()[destChain].Client, []*address.Address{&receiver}, []tlb.Coins{tlb.MustFromTON("1")})

		// Subscribe to OffRamp contract transactions in background
		offRampAddr := state.TonChains[destChain].OffRamp
		tonChain := e.Env.BlockChains.TonChains()[destChain]

		// Get current account state to start monitoring from
		SubscribeOfframpTransactions(t, tonChain, offRampAddr)

		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiverBase64Bytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // state would be failed
			},
		)
	})

	_ = out
}

func SubscribeOfframpTransactions(t *testing.T, tonChain ton.Chain, offRampAddr address.Address) {
	t.Logf("=== STARTING OFFRAMP SUBSCRIPTION for %s ===", offRampAddr.String())
	ctx := t.Context() // Use test context to keep subscription alive until test ends

	master, err := tonChain.Client.CurrentMasterchainInfo(ctx)
	require.NoError(t, err)

	acc, err := tonChain.Client.GetAccount(ctx, master, &offRampAddr)
	require.NoError(t, err)
	lastProcessedLT := acc.LastTxLT
	t.Logf("Starting from LT: %d", lastProcessedLT)

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
				t.Logf("=== TRANSACTION #%d RECEIVED from %s ===", transactionCount, offRampAddr.String())
				t.Logf("LT: %d, Hash: %x", tx.LT, tx.Hash)

				// Log input message
				if tx.IO.In != nil {
					if tx.IO.In.MsgType == tlb.MsgTypeInternal { // Internal
						if internal := tx.IO.In.AsInternal(); internal != nil {
							t.Logf("INPUT Internal - From: %s, To: %s, Value: %s TON",
								internal.SrcAddr.String(),
								internal.DstAddr.String(),
								internal.Amount.TON())
							if internal.Body != nil {
								bodyHex := fmt.Sprintf("%x", internal.Body.ToBOC())
								t.Logf("INPUT Body (hex): %s", bodyHex)

								// Try to parse as Commit message
								parser := internal.Body.BeginParse()
								var ev offramp.Commit
								err := tlb.LoadFromCell(&ev, parser)
								if err == nil {
									t.Logf("✅ Parsed Commit msg: %+v", ev)
								} else {
									t.Logf("❌ Failed to parse Commit msg from cell: %v", err)

									// Try to parse as Execute message
									parser2 := internal.Body.BeginParse()
									var ev2 offramp.Execute
									err2 := tlb.LoadFromCell(&ev2, parser2)
									if err2 == nil {
										t.Logf("✅ Parsed Execute msg: %+v", ev2)
									} else {
										t.Logf("❌ Failed to parse Execute msg from cell: %v", err2)

										// Try to parse as generic message
										parser3 := internal.Body.BeginParse()
										if parser3.BitsLeft() >= 32 {
											if opcode, err := parser3.LoadUInt(32); err == nil {
												t.Logf("📋 Message opcode: 0x%x (%d)", opcode, opcode)
											}
										}
									}
								}
							}
						}
					} else if tx.IO.In.MsgType == tlb.MsgTypeExternalIn {
						if extIn := tx.IO.In.AsExternalIn(); extIn != nil {
							t.Logf("INPUT ExternalIn - To: %s", extIn.DstAddr.String())
							if extIn.Body != nil {
								bodyHex := fmt.Sprintf("%x", extIn.Body.ToBOC())
								t.Logf("INPUT Body (hex): %s", bodyHex)
							}
						}
					}
				}

				// Log output messages
				if tx.IO.Out != nil {
					if msgs, err := tx.IO.Out.ToSlice(); err == nil {
						t.Logf("OUTPUT messages: %d", len(msgs))
						for i, msg := range msgs {
							if msg.MsgType == tlb.MsgTypeInternal {
								if internal := msg.AsInternal(); internal != nil {
									t.Logf("OUTPUT[%d] Internal - From: %s, To: %s, Value: %s TON",
										i, internal.SrcAddr.String(),
										internal.DstAddr.String(),
										internal.Amount.TON())
									if internal.Body != nil {
										bodyHex := fmt.Sprintf("%x", internal.Body.ToBOC())
										t.Logf("OUTPUT[%d] Body (hex): %s", i, bodyHex)

										// Try to parse the output message
										parser := internal.Body.BeginParse()
										if parser.BitsLeft() >= 32 {
											if opcode, err := parser.LoadUInt(32); err == nil {
												t.Logf("OUTPUT[%d] Opcode: 0x%x (%d)", i, opcode, opcode)
											}
										}
									}
								}
							} else if msg.MsgType == tlb.MsgTypeExternalOut { // ExternalOut
								if extOut := msg.AsExternalOut(); extOut != nil {
									t.Logf("OUTPUT[%d] ExternalOut - From: %s", i, extOut.SrcAddr.String())
									if extOut.Body != nil {
										bodyHex := fmt.Sprintf("%x", extOut.Body.ToBOC())
										t.Logf("OUTPUT[%d] Body (hex): %s", i, bodyHex)

										// Try to parse opcode
										parser := extOut.Body.BeginParse()
										if parser.BitsLeft() >= 32 {
											if opcode, err := parser.LoadUInt(32); err == nil {
												t.Logf("OUTPUT[%d] Opcode: 0x%x (%d)", i, opcode, opcode)
											}
										}
									}
								}
							}
						}
					}
				}

				receivedMessage, err := tracetracking.MapToReceivedMessage(tx)
				if err != nil {
					t.Logf("Failed to map to ReceivedMessage: %v", err)
					continue
				}

				err = receivedMessage.WaitForTrace(tonChain.Client)
				if err != nil {
					t.Logf("failed to wait for trace: %v", err)
					continue
				}

				if receivedMessage.ExitCode != 0 {
					t.Logf("transaction failed: with exitcode %d: %s", receivedMessage.ExitCode, receivedMessage.ExitCode.Describe())
					continue
				}

				lm, lmerr := waitForReceivedMsgFlatten(t, tonChain.Client, &receivedMessage)
				if lmerr != nil {
					t.Logf("failed to flatten messages: %v", lmerr)
					continue
				}

				// Log the flattened message details
				t.Logf("Flattened message details:")
				t.Logf("  - OutgoingInternalReceivedMessages: %d", len(lm.OutgoingInternalReceivedMessages))
				t.Logf("  - OutgoingExternalMessages: %d", len(lm.OutgoingExternalMessages))

				if len(lm.OutgoingExternalMessages) > 0 {
					ext := lm.OutgoingExternalMessages[0]
					t.Logf("External message: LT %d, Hash %x", ext.LT, ext.Body.ToBOC())
				} else {
					t.Logf("No external messages found")
				}

				// Also log internal messages for debugging
				for i, internal := range lm.OutgoingInternalReceivedMessages {
					t.Logf("Internal message %d: exit code %v, success: %v", i, internal.ExitCode, internal.Success)
				}

				t.Logf("=== END TRANSACTION ===")

			case <-heartbeat.C:
				t.Logf("💓 Heartbeat: Still monitoring %s, received %d transactions so far", offRampAddr.String(), transactionCount)

			case <-ctx.Done():
				t.Logf("🛑 Subscription context cancelled: %v", ctx.Err())
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

		t.Logf("Flattening %d outgoing internal messages", len(currentMsg.OutgoingInternalReceivedMessages))

		for i, outMsg := range currentMsg.OutgoingInternalReceivedMessages {
			t.Logf("Outgoing message %d: exit code %v, success: %v, bounced: %v, status: %v",
				i, outMsg.ExitCode, outMsg.Success, outMsg.EmittedBouncedMessage, outMsg.Status())

			if outMsg.ExitCode != 0 {
				t.Logf("Outgoing message %d failed with exit code %v", i, outMsg.ExitCode)
			}
			if !outMsg.Success {
				t.Logf("Outgoing message %d was not successful", i)
			}
			if outMsg.EmittedBouncedMessage {
				t.Logf("Outgoing message %d was bounced", i)
			}

			err := outMsg.WaitForTrace(clientConn)
			if err != nil {
				t.Logf("failed to wait for trace: %v", err)
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

	// It's okay if there are no external messages - we'll just log what we have
	t.Logf("Final message processing complete:")
	t.Logf("  - OutgoingInternalReceivedMessages: %d", len(lastMsg.OutgoingInternalReceivedMessages))
	t.Logf("  - OutgoingExternalMessages: %d", len(lastMsg.OutgoingExternalMessages))

	// var event any
	// err := tlb.LoadFromCell(&event, lastMsg.OutgoingExternalMessages[0].Body.BeginParse())
	// if err != nil {
	// 	t.Logf("failed to parse CCIPMessageSent from cell: %v", err)
	// 	return nil, err
	// }

	return lastMsg, nil
}
