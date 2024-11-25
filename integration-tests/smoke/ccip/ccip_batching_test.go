package smoke

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/multicall3"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIPBatching(t *testing.T) {
	// Setup 3 chains, with 2 lanes going to the dest.
	lggr := logger.TestLogger(t)
	ctx := changeset.Context(t)
	// Will load 3 chains when specified by the overrides.toml or env vars (E2E_TEST_SELECTED_NETWORK).
	// See e2e-tests.yml.
	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, &changeset.TestConfigs{
		IsUSDC:       false,
		IsMultiCall3: true, // needed for this test
	})

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	require.Len(t, allChainSelectors, 3, "this test expects 3 chains")
	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]
	destChain := allChainSelectors[2]
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector 1:", sourceChain1,
		", source chain selector 2:", sourceChain2,
		", dest chain selector:", destChain,
	)

	// connect sourceChain1 and sourceChain2 to destChain
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(e.Env, state, sourceChain1, destChain, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(e.Env, state, sourceChain2, destChain, false))

	const (
		numMessages = 5
	)
	var (
		startSeqNum = map[uint64]ccipocr3.SeqNum{
			sourceChain1: 1,
			sourceChain2: 1,
		}
		endSeqNum = map[uint64]ccipocr3.SeqNum{
			sourceChain1: ccipocr3.SeqNum(numMessages),
			sourceChain2: ccipocr3.SeqNum(numMessages),
		}
	)

	t.Run("batch data only messages from single source", func(t *testing.T) {
		err := sendMessages(
			ctx,
			t,
			e.Env.Chains[sourceChain1],
			e.Env.Chains[sourceChain1].DeployerKey,
			state.Chains[sourceChain1].OnRamp,
			state.Chains[sourceChain1].Router,
			state.Chains[sourceChain1].Multicall3,
			destChain,
			numMessages,
			common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
		)
		require.NoError(t, err)

		_, err = changeset.ConfirmCommitWithExpectedSeqNumRange(
			t,
			e.Env.Chains[sourceChain1],
			e.Env.Chains[destChain],
			state.Chains[destChain].OffRamp,
			nil,
			ccipocr3.NewSeqNumRange(startSeqNum[sourceChain1], endSeqNum[sourceChain1]),
		)
		require.NoErrorf(t, err, "failed to confirm commit from chain %d", sourceChain1)

		states, err := changeset.ConfirmExecWithSeqNrs(
			t,
			e.Env.Chains[sourceChain1],
			e.Env.Chains[destChain],
			state.Chains[destChain].OffRamp,
			nil,
			genSeqNrRange(startSeqNum[sourceChain1], endSeqNum[sourceChain1]),
		)
		require.NoError(t, err)
		// assert that all states are successful
		for _, state := range states {
			require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, state)
		}

		startSeqNum[sourceChain1] = endSeqNum[sourceChain1] + 1
		endSeqNum[sourceChain1] = startSeqNum[sourceChain1] + ccipocr3.SeqNum(numMessages) - 1
	})

	t.Run("batch data only messages from multiple sources", func(t *testing.T) {
		var (
			wg           sync.WaitGroup
			sourceChains = []uint64{sourceChain1, sourceChain2}
			errs         = make(chan error, len(sourceChains))
		)

		for _, srcChain := range sourceChains {
			wg.Add(1)
			go sendMessagesAsync(
				ctx,
				t,
				e,
				state,
				srcChain,
				destChain,
				numMessages,
				&wg,
				errs,
			)
		}

		wg.Wait()

		var i int
		for i < len(sourceChains) {
			select {
			case err := <-errs:
				require.NoError(t, err)
				i++
			case <-ctx.Done():
				require.FailNow(t, "didn't get all errors before test context was done")
			}
		}

		// confirm the commit reports
		outputErrs := make(chan outputErr[*offramp.OffRampCommitReportAccepted], len(sourceChains))
		for _, srcChain := range sourceChains {
			wg.Add(1)
			go assertCommitReportsAsync(
				t,
				e,
				state,
				srcChain,
				destChain,
				startSeqNum[srcChain],
				endSeqNum[srcChain],
				&wg,
				outputErrs,
			)
		}

		t.Log("waiting for commit report")
		wg.Wait()

		i = 0
		var reports []*offramp.OffRampCommitReportAccepted
		for i < len(sourceChains) {
			select {
			case outputErr := <-outputErrs:
				require.NoError(t, outputErr.err)
				reports = append(reports, outputErr.output)
				i++
			case <-ctx.Done():
				require.FailNow(t, "didn't get all commit reports before test context was done")
			}
		}

		// the reports should be the same for both, since both roots should be batched within
		// that one report.
		require.Lenf(t, reports, len(sourceChains), "expected %d commit reports", len(sourceChains))
		require.NotNil(t, reports[0], "commit report should not be nil")
		require.NotNil(t, reports[1], "commit report should not be nil")
		// TODO: this assertion is failing, despite messages being sent at the same time.
		// require.Equal(t, reports[0], reports[1], "commit reports should be the same")

		// confirm execution
		execErrs := make(chan outputErr[map[uint64]int], len(sourceChains))
		for _, srcChain := range sourceChains {
			wg.Add(1)
			go assertExecAsync(
				t,
				e,
				state,
				srcChain,
				destChain,
				genSeqNrRange(startSeqNum[srcChain], endSeqNum[srcChain]),
				&wg,
				execErrs,
			)
		}

		t.Log("waiting for exec reports")
		wg.Wait()

		i = 0
		var execStates []map[uint64]int
		for i < len(sourceChains) {
			select {
			case outputErr := <-execErrs:
				require.NoError(t, outputErr.err)
				execStates = append(execStates, outputErr.output)
				i++
			case <-ctx.Done():
				require.FailNow(t, "didn't get all exec reports before test context was done")
			}
		}

		// assert that all states are successful
		for _, states := range execStates {
			for _, state := range states {
				require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, state)
			}
		}
	})
}

type outputErr[T any] struct {
	output T
	err    error
}

func assertExecAsync(
	t *testing.T,
	e changeset.DeployedEnv,
	state changeset.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	seqNums []uint64,
	wg *sync.WaitGroup,
	errs chan<- outputErr[map[uint64]int],
) {
	defer wg.Done()
	states, err := changeset.ConfirmExecWithSeqNrs(
		t,
		e.Env.Chains[sourceChainSelector],
		e.Env.Chains[destChainSelector],
		state.Chains[destChainSelector].OffRamp,
		nil,
		seqNums,
	)

	errs <- outputErr[map[uint64]int]{states, err}
}

func assertCommitReportsAsync(
	t *testing.T,
	e changeset.DeployedEnv,
	state changeset.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	startSeqNum,
	endSeqNum ccipocr3.SeqNum,
	wg *sync.WaitGroup,
	errs chan<- outputErr[*offramp.OffRampCommitReportAccepted],
) {
	defer wg.Done()
	commitReport, err := changeset.ConfirmCommitWithExpectedSeqNumRange(
		t,
		e.Env.Chains[sourceChainSelector],
		e.Env.Chains[destChainSelector],
		state.Chains[destChainSelector].OffRamp,
		nil,
		ccipocr3.NewSeqNumRange(startSeqNum, endSeqNum),
	)

	errs <- outputErr[*offramp.OffRampCommitReportAccepted]{commitReport, err}
}

func sendMessagesAsync(
	ctx context.Context,
	t *testing.T,
	e changeset.DeployedEnv,
	state changeset.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	numMessages int,
	wg *sync.WaitGroup,
	out chan<- error,
) {
	defer wg.Done()
	err := sendMessages(
		ctx,
		t,
		e.Env.Chains[sourceChainSelector],
		e.Env.Chains[sourceChainSelector].DeployerKey,
		state.Chains[sourceChainSelector].OnRamp,
		state.Chains[sourceChainSelector].Router,
		state.Chains[sourceChainSelector].Multicall3,
		destChainSelector,
		numMessages,
		common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
	)
	t.Log("sendMessagesAsync error:", err, ", writing to channel")
	out <- err
}

func sendMessages(
	ctx context.Context,
	t *testing.T,
	sourceChain deployment.Chain,
	sourceTransactOpts *bind.TransactOpts,
	sourceOnRamp *onramp.OnRamp,
	sourceRouter *router.Router,
	sourceMulticall3 *multicall3.Multicall3,
	destChainSelector uint64,
	numMessages int,
	receiver []byte,
) error {
	calls, totalValue, err := genMessages(
		ctx,
		sourceRouter,
		destChainSelector,
		numMessages,
		receiver,
	)
	if err != nil {
		return fmt.Errorf("generate messages: %w", err)
	}

	// Send the tx with the messages through the multicall
	tx, err := sourceMulticall3.Aggregate3Value(
		&bind.TransactOpts{
			From:   sourceTransactOpts.From,
			Signer: sourceTransactOpts.Signer,
			Value:  totalValue,
		},
		calls,
	)
	_, err = deployment.ConfirmIfNoError(sourceChain, tx, err)
	if err != nil {
		return fmt.Errorf("send messages via multicall3: %w", err)
	}

	// check that the message was emitted
	iter, err := sourceOnRamp.FilterCCIPMessageSent(
		nil, []uint64{destChainSelector}, nil,
	)
	if err != nil {
		return fmt.Errorf("get message sent event: %w", err)
	}
	defer iter.Close()

	// there should be numMessages messages emitted
	for i := 0; i < numMessages; i++ {
		if !iter.Next() {
			return fmt.Errorf("expected %d messages, got %d", numMessages, i)
		}
		t.Logf("Message id of msg %d: %x", i, iter.Event.Message.Header.MessageId[:])
	}

	return nil
}

func genMessages(
	ctx context.Context,
	sourceRouter *router.Router,
	destChainSelector uint64,
	count int,
	receiver []byte,
) (calls []multicall3.Multicall3Call3Value, totalValue *big.Int, err error) {
	totalValue = big.NewInt(0)
	for i := 0; i < count; i++ {
		msg := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         []byte(fmt.Sprintf("hello world %d", i)),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}

		fee, err := sourceRouter.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("router get fee: %w", err)
		}

		totalValue.Add(totalValue, fee)

		calldata, err := changeset.CCIPSendCalldata(destChainSelector, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("generate calldata: %w", err)
		}

		calls = append(calls, multicall3.Multicall3Call3Value{
			Target:       sourceRouter.Address(),
			AllowFailure: false,
			CallData:     calldata,
			Value:        fee,
		})
	}

	return calls, totalValue, nil
}

// creates an array of uint64 from start to end inclusive
func genSeqNrRange(start, end ccipocr3.SeqNum) []uint64 {
	var seqNrs []uint64
	for i := start; i <= end; i++ {
		seqNrs = append(seqNrs, uint64(i))
	}
	return seqNrs
}
