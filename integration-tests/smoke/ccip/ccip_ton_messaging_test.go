package ccip

import (
	"fmt"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/xssnick/tonutils-go/tlb"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_TON2EVM(t *testing.T) {
	// setup environment with 1 ton chain
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	// load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// get chain selectors
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := allTonChainSelectors[0]
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors) // make evm chains sorted for deterministic test results
	destChain := evmChainSelectors[0]
	t.Log("Chain selectors",
		"TON", allTonChainSelectors,
		"EVM", evmChainSelectors,
		"home", e.HomeChainSel,
		"feed", e.FeedChainSel,
		"source", sourceChain,
		"dest", destChain,
	)

	// setup lane
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// encode sender address(deployer address)
	ac := codec.NewAddressCodec()
	tonChain := e.Env.BlockChains.TonChains()[sourceChain]
	addrBytes, err := ac.AddressStringToBytes(tonChain.WalletAddress.String())
	require.NoError(t, err)

	// wait for event filter registration
	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)
	// ready to test
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
		extraArgs, err := tlb.ToCell(onramp.GenericExtraArgsV2{
			GasLimit:                 big.NewInt(1000000),
			AllowOutOfOrderExecution: true,
		})
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
				ExtraArgs:              extraArgs.ToBOC(),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	_ = out
}

func Test_CCIPMessaging_EVM2TON(t *testing.T) {
	// setup environment with 1 ton chain
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	// load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// get chain selectors
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]

	t.Log("Chain selectors",
		"TON", allTonChainSelectors,
		"EVM", evmChainSelectors,
		"home", e.HomeChainSel,
		"feed", e.FeedChainSel,
		"source", sourceChain,
		"dest", destChain,
	)
	t.Logf("  OnRamp:       %s", state.Chains[sourceChain].OnRamp.Address())

	// setup lane
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// wait for event filter registration
	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	// ready to test
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

	t.Run("message to contract receiver", func(t *testing.T) {
		offRampAddr := state.TonChains[destChain].OffRamp
		receiverAddr := state.TonChains[destChain].ReceiverAddress

		t.Logf("  TON OffRamp:  %s", offRampAddr.String())
		t.Logf("  TON Receiver: %s", receiverAddr.String())

		ac := codec.NewAddressCodec()
		receiverBytes, err := ac.AddressStringToBytes(receiverAddr.String())
		require.NoError(t, err)
		require.Len(t, receiverBytes, 36, "receiver bytes should be 36 bytes")

		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiverBytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
	_ = out
}

// Test_CCIPMessaging_EVM2TON_BatchCommit sends multiple messages in a single transaction
// using Multicall3 to test batch commit with multiple execution reports on TON.
// This test follows the same pattern as Test_CCIPMessaging_MultiExecReports_EVM2Solana:
// - Sends 5 messages in one Multicall3 transaction
// - All messages committed in one batch (single commit report)
// - Each message executed in a separate execution report (MaxReportMessages=1)
func Test_CCIPMessaging_EVM2TON_BatchCommit(t *testing.T) {
	ctx := testhelpers.Context(t)

	const numMessages = 5

	// setup environment with 1 ton chain and Multicall3
	// Configure execution to create multiple reports (one per message) from batch commit
	// This matches the EVM2Solana test pattern
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithMultiCall3(),
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			// Configure execution plugin for multiple reports (same as Solana test)
			params.ExecuteOffChainConfig.InflightCacheExpiry = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MessageVisibilityInterval = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			params.ExecuteOffChainConfig.MaxReportMessages = 1 // 1 message per execution report
			params.ExecuteOffChainConfig.MaxSingleChainReports = 1
			return params
		}),
	)

	// load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// get chain selectors
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]

	t.Log("Chain selectors",
		"TON", allTonChainSelectors,
		"EVM", evmChainSelectors,
		"home", e.HomeChainSel,
		"feed", e.FeedChainSel,
		"source", sourceChain,
		"dest", destChain,
	)

	// setup lane with EnforceOutOfOrder=true (required when MultipleReportsEnabled=true)
	testhelpers.AddLaneWithEnforceOutOfOrder(t, &e, state, sourceChain, destChain, false)

	// wait for event filter registration
	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	// prepare receiver address
	receiverAddr := state.TonChains[destChain].ReceiverAddress
	t.Logf("  TON Receiver: %s", receiverAddr.String())

	ac := codec.NewAddressCodec()
	receiverBytes, err := ac.AddressStringToBytes(receiverAddr.String())
	require.NoError(t, err)
	require.Len(t, receiverBytes, 36, "receiver bytes should be 36 bytes")

	// send multiple messages using Multicall3
	evmChain := e.Env.BlockChains.EVMChains()[sourceChain]
	onRamp := state.Chains[sourceChain].OnRamp
	routerContract := state.Chains[sourceChain].Router
	multicall3 := state.Chains[sourceChain].Multicall3

	t.Logf("Sending %d messages via Multicall3", numMessages)

	// get next expected sequence number before sending
	nextSeqNum, err := onRamp.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: ctx}, destChain)
	require.NoError(t, err)
	t.Logf("Expected next sequence number: %d", nextSeqNum)

	// generate messages for multicall3
	calls, totalValue, err := testhelpers.GenMessagesForMulticall3(
		ctx,
		routerContract,
		destChain,
		numMessages,
		router.ClientEVM2AnyMessage{
			Receiver:     receiverBytes,
			Data:         fmt.Appendf(nil, "batch message %d", numMessages),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    testhelpers.MakeEVMExtraArgsV2(100000, false),
		},
	)
	require.NoError(t, err)

	// check balance
	currBalance, err := evmChain.Client.BalanceAt(ctx, evmChain.DeployerKey.From, nil)
	require.NoError(t, err)
	t.Logf("Sending %d messages with total value %s, current balance: %s",
		numMessages, totalValue.String(), currBalance.String())
	require.True(t, currBalance.Cmp(totalValue) >= 0, "insufficient balance")

	// send via multicall3
	tx, err := multicall3.Aggregate3Value(
		&bind.TransactOpts{
			From:   evmChain.DeployerKey.From,
			Signer: evmChain.DeployerKey.Signer,
			Value:  totalValue,
		},
		calls,
	)
	require.NoError(t, err)

	_, err = deployment.ConfirmIfNoError(evmChain, tx, err)
	require.NoError(t, err)
	t.Logf("Multicall3 tx sent: %s", tx.Hash().Hex())

	// verify all messages were emitted
	iter, err := onRamp.FilterCCIPMessageSent(
		nil, []uint64{destChain}, nil,
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, iter.Close())
	}()

	// collect all message IDs
	var messageIDs []common.Hash
	for i := 0; i < numMessages; i++ {
		require.True(t, iter.Next(), "expected %d messages, got %d", numMessages, i)
		messageIDs = append(messageIDs, iter.Event.Message.Header.MessageId)
		t.Logf("Message %d: seqNum=%d, msgID=%x",
			i+1, iter.Event.SequenceNumber, iter.Event.Message.Header.MessageId[:8])
	}
	require.Len(t, messageIDs, numMessages, "expected %d messages", numMessages)

	// confirm commit with expected sequence number range
	expectedSeqNumRange := ccipocr3.NewSeqNumRange(
		ccipocr3.SeqNum(nextSeqNum),
		ccipocr3.SeqNum(nextSeqNum+uint64(numMessages)-1),
	)

	t.Logf("Waiting for commit report with seq range [%d, %d]",
		expectedSeqNumRange.Start(), expectedSeqNumRange.End())

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		map[uint64]*uint64{destChain: nil}, // startBlocks
		true,                               // enforceSingleCommit - all messages should be in one commit
		map[testhelpers.SourceDestPair]ccipocr3.SeqNumRange{
			{
				SourceChainSelector: sourceChain,
				DestChainSelector:   destChain,
			}: expectedSeqNumRange,
		},
	)
	require.NoError(t, err)
	t.Logf("✓ Batch commit verified: all %d messages committed in single report", numMessages)

	// confirm execution for all messages
	var expectedSeqNums []uint64
	for i := uint64(0); i < uint64(numMessages); i++ {
		expectedSeqNums = append(expectedSeqNums, nextSeqNum+i)
	}

	t.Logf("Waiting for execution of messages: %v", expectedSeqNums)
	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		map[testhelpers.SourceDestPair][]uint64{
			{
				SourceChainSelector: sourceChain,
				DestChainSelector:   destChain,
			}: expectedSeqNums,
		},
		map[uint64]*uint64{destChain: nil}, // startBlocks
	)

	// verify all executions succeeded
	for pair, states := range execStates {
		t.Logf("Execution states for %d->%d: %v", pair.SourceChainSelector, pair.DestChainSelector, states)
		for seqNum, state := range states {
			require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, state,
				"message %d execution failed", seqNum)
		}
	}
	t.Logf("✓ All %d messages executed successfully", numMessages)
}
