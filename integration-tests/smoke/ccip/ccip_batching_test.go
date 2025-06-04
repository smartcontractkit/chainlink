package ccip

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
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/multicall3"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

const (
	numMessages = 5
)

type batchTestSetup struct {
	e            testhelpers.DeployedEnv
	state        stateview.CCIPOnChainState
	sourceChain1 uint64
	sourceChain2 uint64
	destChain    uint64
}

func newBatchTestSetup(t *testing.T, opts ...testhelpers.TestOps) batchTestSetup {
	// Setup 3 chains, with 2 lanes going to the dest.
	options := []testhelpers.TestOps{
		testhelpers.WithMultiCall3(),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfUsersPerChain(2),
	}
	options = append(options, opts...)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		options...,
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
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
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain1, destChain, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain2, destChain, false)

	return batchTestSetup{e, state, sourceChain1, sourceChain2, destChain}
}

func Test_CCIPBatching_MaxBatchSizeEVM(t *testing.T) {
	ctx := testhelpers.Context(t)
	setup := newBatchTestSetup(t)
	sourceChain1, sourceChain2, destChain, e, state := setup.sourceChain1, setup.sourceChain2, setup.destChain, setup.e, setup.state
	evmChains := e.Env.BlockChains.EVMChains()

	var (
		startSeqNum = map[uint64]ccipocr3.SeqNum{
			sourceChain1: 1,
			sourceChain2: 1,
		}
		sourceChain = sourceChain1
		transactors = []*bind.TransactOpts{
			evmChains[sourceChain].DeployerKey,
			evmChains[sourceChain].Users[0],
		}
		errs = make(chan error, len(transactors))
	)

	for _, transactor := range transactors {
		go func() {
			err := sendMessages(
				ctx,
				t,
				evmChains[sourceChain],
				transactor,
				state.MustGetEVMChainState(sourceChain).OnRamp,
				state.MustGetEVMChainState(sourceChain).Router,
				state.MustGetEVMChainState(sourceChain).Multicall3,
				destChain,
				merklemulti.MaxNumberTreeLeaves/2,
				common.LeftPadBytes(state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(), 32),
			)
			t.Log("sendMessages error:", err, ", writing to channel")
			errs <- err
			t.Log("sent error to channel")
		}()
	}

	var i = 0
	for i < len(transactors) {
		select {
		case err := <-errs:
			require.NoError(t, err)
			i++
		case <-ctx.Done():
			require.FailNow(t, "didn't get all errors before test context was done")
		}
	}

	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceChain,
		evmChains[destChain],
		state.MustGetEVMChainState(destChain).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(
			startSeqNum[sourceChain],
			startSeqNum[sourceChain]+ccipocr3.SeqNum(merklemulti.MaxNumberTreeLeaves)-1,
		),
		true,
	)
	require.NoErrorf(t, err, "failed to confirm commit from chain %d", sourceChain)
}

var multiRootReportOverride = testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
	params.CommitOffChainConfig.MultipleReportsEnabled = true
	params.CommitOffChainConfig.MaxMerkleRootsPerReport = 1
	return params
})

var multiPriceReportOverride = testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
	params.CommitOffChainConfig.MultipleReportsEnabled = true
	params.CommitOffChainConfig.MaxMerkleRootsPerReport = 1
	params.CommitOffChainConfig.MaxPricesPerReport = 1
	return params
})

func Test_CCIPBatching_MultiSource(t *testing.T) {
	ccipBatchingMultiSource(t)
}

func Test_CCIPBatching_MultiSource_MultiRoot(t *testing.T) {
	ccipBatchingMultiSource(t, multiRootReportOverride)
}

func Test_CCIPBatching_MultiSource_MultiPrice(t *testing.T) {
	ccipBatchingMultiSource(t, multiPriceReportOverride)
}

func Test_CCIPBatching_SingleSource(t *testing.T) {
	ccipBatchingSingleSource(t)
}

func Test_CCIPBatching_SingleSource_MultiRoot(t *testing.T) {
	ccipBatchingMultiSource(t, multiRootReportOverride)
}

func Test_CCIPBatching_SingleSource_MultiPrice(t *testing.T) {
	ccipBatchingMultiSource(t, multiPriceReportOverride)
}

func ccipBatchingMultiSource(t *testing.T, opts ...testhelpers.TestOps) {
	// Setup 3 chains, with 2 lanes going to the dest.
	ctx := testhelpers.Context(t)
	setup := newBatchTestSetup(t, opts...)
	sourceChain1, sourceChain2, destChain, e, state := setup.sourceChain1, setup.sourceChain2, setup.destChain, setup.e, setup.state

	var (
		wg           sync.WaitGroup
		sourceChains = []uint64{sourceChain1, sourceChain2}
		errs         = make(chan error, len(sourceChains))
		startSeqNum  = map[uint64]ccipocr3.SeqNum{
			sourceChain1: 1,
			sourceChain2: 1,
		}
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
			startSeqNum[srcChain]+ccipocr3.SeqNum(numMessages)-1,
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
			genSeqNrRange(startSeqNum[srcChain], startSeqNum[srcChain]+ccipocr3.SeqNum(numMessages)-1),
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
			require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, state)
		}
	}
}

func ccipBatchingSingleSource(t *testing.T, opts ...testhelpers.TestOps) {
	// Setup 3 chains, with 2 lanes going to the dest.
	ctx := testhelpers.Context(t)
	setup := newBatchTestSetup(t, opts...)
	sourceChain1, sourceChain2, destChain, e, state := setup.sourceChain1, setup.sourceChain2, setup.destChain, setup.e, setup.state
	evmChains := e.Env.BlockChains.EVMChains()

	var (
		startSeqNum = map[uint64]ccipocr3.SeqNum{
			sourceChain1: 1,
			sourceChain2: 1,
		}
	)

	var (
		sourceChain = sourceChain1
	)
	err := sendMessages(
		ctx,
		t,
		evmChains[sourceChain],
		evmChains[sourceChain].DeployerKey,
		state.MustGetEVMChainState(sourceChain).OnRamp,
		state.MustGetEVMChainState(sourceChain).Router,
		state.MustGetEVMChainState(sourceChain).Multicall3,
		destChain,
		numMessages,
		common.LeftPadBytes(state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(), 32),
	)
	require.NoError(t, err)

	_, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceChain,
		evmChains[destChain],
		state.MustGetEVMChainState(destChain).OffRamp,
		nil,
		ccipocr3.NewSeqNumRange(startSeqNum[sourceChain], startSeqNum[sourceChain]+numMessages-1),
		true,
	)
	require.NoErrorf(t, err, "failed to confirm commit from chain %d", sourceChain)

	states, err := testhelpers.ConfirmExecWithSeqNrs(
		t,
		sourceChain,
		evmChains[destChain],
		state.MustGetEVMChainState(destChain).OffRamp,
		nil,
		genSeqNrRange(startSeqNum[sourceChain], startSeqNum[sourceChain]+numMessages-1),
	)
	require.NoError(t, err)
	// assert that all states are successful
	for _, state := range states {
		require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, state)
	}
}

type outputErr[T any] struct {
	output T
	err    error
}

func assertExecAsync(
	t *testing.T,
	e testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	seqNums []uint64,
	wg *sync.WaitGroup,
	errs chan<- outputErr[map[uint64]int],
) {
	defer wg.Done()
	states, err := testhelpers.ConfirmExecWithSeqNrs(
		t,
		sourceChainSelector,
		e.Env.BlockChains.EVMChains()[destChainSelector],
		state.MustGetEVMChainState(destChainSelector).OffRamp,
		nil,
		seqNums,
	)

	errs <- outputErr[map[uint64]int]{states, err}
}

func assertCommitReportsAsync(
	t *testing.T,
	e testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	startSeqNum,
	endSeqNum ccipocr3.SeqNum,
	wg *sync.WaitGroup,
	errs chan<- outputErr[*offramp.OffRampCommitReportAccepted],
) {
	defer wg.Done()
	commitReport, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceChainSelector,
		e.Env.BlockChains.EVMChains()[destChainSelector],
		state.MustGetEVMChainState(destChainSelector).OffRamp,
		nil,
		ccipocr3.NewSeqNumRange(startSeqNum, endSeqNum),
		true,
	)

	errs <- outputErr[*offramp.OffRampCommitReportAccepted]{commitReport, err}
}

func sendMessagesAsync(
	ctx context.Context,
	t *testing.T,
	e testhelpers.DeployedEnv,
	state stateview.CCIPOnChainState,
	sourceChainSelector,
	destChainSelector uint64,
	numMessages int,
	wg *sync.WaitGroup,
	out chan<- error,
) {
	defer wg.Done()
	var err error

	const (
		numRetries = 3
	)

	// we retry a bunch of times just in case there is a race b/w the prices being
	// posted and the messages being sent.
	for i := 0; i < numRetries; i++ {
		err = sendMessages(
			ctx,
			t,
			e.Env.BlockChains.EVMChains()[sourceChainSelector],
			e.Env.BlockChains.EVMChains()[sourceChainSelector].DeployerKey,
			state.MustGetEVMChainState(sourceChainSelector).OnRamp,
			state.MustGetEVMChainState(sourceChainSelector).Router,
			state.MustGetEVMChainState(sourceChainSelector).Multicall3,
			destChainSelector,
			numMessages,
			common.LeftPadBytes(state.MustGetEVMChainState(destChainSelector).Receiver.Address().Bytes(), 32),
		)
		if err == nil {
			break
		}

		t.Log("sendMessagesAsync error is non-nil:", err, ", retrying")
	}

	t.Log("sendMessagesAsync error:", err, ", writing to channel")
	out <- err
}

func sendMessages(
	ctx context.Context,
	t *testing.T,
	sourceChain cldf_evm.Chain,
	sourceTransactOpts *bind.TransactOpts,
	sourceOnRamp onramp.OnRampInterface,
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

	currBalance, err := sourceChain.Client.BalanceAt(ctx, sourceTransactOpts.From, nil)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	// Send the tx with the messages through the multicall
	t.Logf("Sending %d messages with total value %s, current balance: %s", numMessages, totalValue.String(), currBalance.String())
	tx, err := sourceMulticall3.Aggregate3Value(
		&bind.TransactOpts{
			From:   sourceTransactOpts.From,
			Signer: sourceTransactOpts.Signer,
			Value:  totalValue,
		},
		calls,
	)
	_, err = cldf.ConfirmIfNoError(sourceChain, tx, err)
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
	defer func() {
		require.NoError(t, iter.Close())
	}()

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
			ExtraArgs:    testhelpers.MakeEVMExtraArgsV2(50_000, false),
		}

		fee, err := sourceRouter.GetFee(&bind.CallOpts{Context: ctx}, destChainSelector, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("router get fee: %w", err)
		}

		totalValue.Add(totalValue, fee)

		calldata, err := testhelpers.CCIPSendCalldata(destChainSelector, msg)
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
