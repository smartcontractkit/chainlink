package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

func ConfirmGasPriceUpdatedForAll(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	startBlocks map[uint64]*uint64,
	gasPrice *big.Int,
) {
	var wg errgroup.Group
	for src, srcChain := range e.Chains {
		for dest, dstChain := range e.Chains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Go(func() error {
				var startBlock *uint64
				if startBlocks != nil {
					startBlock = startBlocks[srcChain.Selector]
				}
				return ConfirmGasPriceUpdated(
					t,
					dstChain,
					state.Chains[srcChain.Selector].FeeQuoter,
					*startBlock,
					gasPrice,
				)
			})
		}
	}
	require.NoError(t, wg.Wait())
}

func ConfirmGasPriceUpdated(
	t *testing.T,
	dest deployment.Chain,
	srcFeeQuoter *fee_quoter.FeeQuoter,
	startBlock uint64,
	gasPrice *big.Int,
) error {
	it, err := srcFeeQuoter.FilterUsdPerUnitGasUpdated(&bind.FilterOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, []uint64{dest.Selector})

	require.NoError(t, err)
	require.Truef(t, it.Next(), "No gas price update event found on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	require.NotEqualf(t, gasPrice, it.Event.Value, "Gas price not updated on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	return nil
}

func ConfirmTokenPriceUpdatedForAll(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	startBlocks map[uint64]*uint64,
	linkPrice *big.Int,
	wethPrice *big.Int,
) {
	var wg errgroup.Group
	for _, chain := range e.Chains {
		chain := chain
		wg.Go(func() error {
			var startBlock *uint64
			if startBlocks != nil {
				startBlock = startBlocks[chain.Selector]
			}
			linkAddress := state.Chains[chain.Selector].LinkToken.Address()
			wethAddress := state.Chains[chain.Selector].Weth9.Address()
			tokenToPrice := make(map[common.Address]*big.Int)
			tokenToPrice[linkAddress] = linkPrice
			tokenToPrice[wethAddress] = wethPrice
			return ConfirmTokenPriceUpdated(
				t,
				chain,
				state.Chains[chain.Selector].FeeQuoter,
				*startBlock,
				tokenToPrice,
			)
		})
	}
	require.NoError(t, wg.Wait())
}

func ConfirmTokenPriceUpdated(
	t *testing.T,
	chain deployment.Chain,
	feeQuoter *fee_quoter.FeeQuoter,
	startBlock uint64,
	tokenToInitialPrice map[common.Address]*big.Int,
) error {
	tokens := make([]common.Address, 0, len(tokenToInitialPrice))
	for token := range tokenToInitialPrice {
		tokens = append(tokens, token)
	}
	it, err := feeQuoter.FilterUsdPerTokenUpdated(&bind.FilterOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, tokens)
	require.NoError(t, err)
	for it.Next() {
		token := it.Event.Token
		initialValue, ok := tokenToInitialPrice[token]
		if ok {
			require.Contains(t, tokens, token)
			// Initial Value should be changed
			require.NotEqual(t, initialValue, it.Event.Value)
		}

		// Remove the token from the map until we assert all tokens are updated
		delete(tokenToInitialPrice, token)
		if len(tokenToInitialPrice) == 0 {
			return nil
		}
	}

	if len(tokenToInitialPrice) > 0 {
		return fmt.Errorf("not all tokens updated on chain  %d", chain.Selector)
	}

	return nil
}

// SourceDestPair is represents a pair of source and destination chain selectors.
// Use this as a key in maps that need to identify sequence numbers, nonces, or
// other things that require identification.
type SourceDestPair struct {
	SourceChainSelector uint64
	DestChainSelector   uint64
}

// ConfirmCommitForAllWithExpectedSeqNums waits for all chains in the environment to commit the given expectedSeqNums.
// expectedSeqNums is a map that maps a (source, dest) selector pair to the expected sequence number
// to confirm the commit for.
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmCommitForAllWithExpectedSeqNums(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	expectedSeqNums map[SourceDestPair]uint64,
	startBlocks map[uint64]*uint64,
) {
	var wg errgroup.Group
	for src, srcChain := range e.Chains {
		for dest, dstChain := range e.Chains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Go(func() error {
				var startBlock *uint64
				if startBlocks != nil {
					startBlock = startBlocks[dstChain.Selector]
				}

				expectedSeqNum, ok := expectedSeqNums[SourceDestPair{
					SourceChainSelector: srcChain.Selector,
					DestChainSelector:   dstChain.Selector,
				}]
				if !ok || expectedSeqNum == 0 {
					return nil
				}

				return commonutils.JustError(ConfirmCommitWithExpectedSeqNumRange(
					t,
					srcChain,
					dstChain,
					state.Chains[dstChain.Selector].OffRamp,
					startBlock,
					ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(expectedSeqNum),
						ccipocr3.SeqNum(expectedSeqNum),
					},
					true,
				))
			})
		}
	}

	done := make(chan struct{})
	go func() {
		require.NoError(t, wg.Wait())
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
		"all commitments did not confirm",
	)
}

type commitReportTracker struct {
	seenMessages map[uint64]map[uint64]bool
}

func newCommitReportTracker(sourceChainSelector uint64, seqNrs ccipocr3.SeqNumRange) commitReportTracker {
	seenMessages := make(map[uint64]map[uint64]bool)
	seenMessages[sourceChainSelector] = make(map[uint64]bool)

	for i := seqNrs.Start(); i <= seqNrs.End(); i++ {
		seenMessages[sourceChainSelector][uint64(i)] = false
	}
	return commitReportTracker{seenMessages: seenMessages}
}

func (c *commitReportTracker) visitCommitReport(sourceChainSelector uint64, minSeqNr uint64, maxSeqNr uint64) {
	if _, ok := c.seenMessages[sourceChainSelector]; !ok {
		return
	}

	for i := minSeqNr; i <= maxSeqNr; i++ {
		c.seenMessages[sourceChainSelector][i] = true
	}
}

func (c *commitReportTracker) allCommited(sourceChainSelector uint64) bool {
	for _, v := range c.seenMessages[sourceChainSelector] {
		if !v {
			return false
		}
	}
	return true
}

// ConfirmMultipleCommits waits for multiple ccipocr3.SeqNumRange to be committed by the Offramp.
// Waiting is done in parallel per every sourceChain/destChain (lane) passed as argument.
func ConfirmMultipleCommits(
	t *testing.T,
	chains map[uint64]deployment.Chain,
	state map[uint64]changeset.CCIPChainState,
	startBlocks map[uint64]*uint64,
	enforceSingleCommit bool,
	expectedSeqNums map[SourceDestPair]ccipocr3.SeqNumRange,
) error {
	errGrp := &errgroup.Group{}

	for sourceDest, seqRange := range expectedSeqNums {
		seqRange := seqRange
		srcChain := sourceDest.SourceChainSelector
		destChain := sourceDest.DestChainSelector

		errGrp.Go(func() error {
			_, err := ConfirmCommitWithExpectedSeqNumRange(
				t,
				chains[srcChain],
				chains[destChain],
				state[destChain].OffRamp,
				startBlocks[destChain],
				seqRange,
				enforceSingleCommit,
			)
			return err
		})
	}

	return errGrp.Wait()
}

// ConfirmCommitWithExpectedSeqNumRange waits for a commit report on the destination chain with the expected sequence number range.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
func ConfirmCommitWithExpectedSeqNumRange(
	t *testing.T,
	src deployment.Chain,
	dest deployment.Chain,
	offRamp *offramp.OffRamp,
	startBlock *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (*offramp.OffRampCommitReportAccepted, error) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe CommitReportAccepted : %w", err)
	}

	seenMessages := newCommitReportTracker(src.Selector, expectedSeqNumRange)

	defer subscription.Unsubscribe()
	var duration time.Duration
	deadline, ok := t.Deadline()
	if ok {
		// make this timer end a minute before so that we don't hit the deadline
		duration = deadline.Sub(time.Now().Add(-1 * time.Minute))
	} else {
		duration = 5 * time.Minute
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// if it's simulated backend, commit to ensure mining
			if backend, ok := src.Client.(*memory.Backend); ok {
				backend.Commit()
			}
			if backend, ok := dest.Client.(*memory.Backend); ok {
				backend.Commit()
			}
			t.Logf("Waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, src.Selector, expectedSeqNumRange.String())

			// Need to do this because the subscription sometimes fails to get the event.
			iter, err := offRamp.FilterCommitReportAccepted(&bind.FilterOpts{
				Context: tests.Context(t),
			})
			require.NoError(t, err)
			for iter.Next() {
				event := iter.Event
				if len(event.MerkleRoots) > 0 {
					for _, mr := range event.MerkleRoots {
						t.Logf("Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v, tx hash: %s",
							mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, src.Selector, expectedSeqNumRange.String(), event.PriceUpdates.TokenPriceUpdates, event.Raw.TxHash.String())
						seenMessages.visitCommitReport(src.Selector, mr.MinSeqNr, mr.MaxSeqNr)

						if mr.SourceChainSelector == src.Selector &&
							uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
							uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
							t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
							return event, nil
						}

						if !enforceSingleCommit && seenMessages.allCommited(src.Selector) {
							t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
							return event, nil
						}
					}
				}
			}
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		case <-timer.C:
			return nil, fmt.Errorf("timed out after waiting %s duration for commit report on chain selector %d from source selector %d expected seq nr range %s",
				duration.String(), dest.Selector, src.Selector, expectedSeqNumRange.String())
		case report := <-sink:
			if len(report.MerkleRoots) > 0 {
				// Check the interval of sequence numbers and make sure it matches
				// the expected range.
				for _, mr := range report.MerkleRoots {
					t.Logf("Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
						mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, src.Selector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates)

					seenMessages.visitCommitReport(src.Selector, mr.MinSeqNr, mr.MaxSeqNr)

					if mr.SourceChainSelector == src.Selector &&
						uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
						uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
						t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
						return report, nil
					}

					if !enforceSingleCommit && seenMessages.allCommited(src.Selector) {
						t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
						return report, nil
					}
				}
			}
		}
	}
}

// ConfirmExecWithSeqNrsForAll waits for all chains in the environment to execute the given expectedSeqNums.
// If successful, it returns a map that maps the SourceDestPair to the expected sequence number
// to its execution state.
// expectedSeqNums is a map of SourceDestPair to a slice of expected sequence numbers to be executed.
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmExecWithSeqNrsForAll(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	expectedSeqNums map[SourceDestPair][]uint64,
	startBlocks map[uint64]*uint64,
) (executionStates map[SourceDestPair]map[uint64]int) {
	var (
		wg errgroup.Group
		mx sync.Mutex
	)
	executionStates = make(map[SourceDestPair]map[uint64]int)
	for src, srcChain := range e.Chains {
		for dest, dstChain := range e.Chains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Go(func() error {
				var startBlock *uint64
				if startBlocks != nil {
					startBlock = startBlocks[dstChain.Selector]
				}

				expectedSeqNum, ok := expectedSeqNums[SourceDestPair{
					SourceChainSelector: srcChain.Selector,
					DestChainSelector:   dstChain.Selector,
				}]
				if !ok || len(expectedSeqNum) == 0 {
					return nil
				}

				innerExecutionStates, err := ConfirmExecWithSeqNrs(
					t,
					srcChain,
					dstChain,
					state.Chains[dstChain.Selector].OffRamp,
					startBlock,
					expectedSeqNum,
				)
				if err != nil {
					return err
				}

				mx.Lock()
				executionStates[SourceDestPair{
					SourceChainSelector: srcChain.Selector,
					DestChainSelector:   dstChain.Selector,
				}] = innerExecutionStates
				mx.Unlock()

				return nil
			})
		}
	}

	require.NoError(t, wg.Wait())
	return executionStates
}

// ConfirmExecWithSeqNrs waits for an execution state change on the destination chain with the expected sequence number.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
// Returns a map that maps the expected sequence number to its execution state.
func ConfirmExecWithSeqNrs(
	t *testing.T,
	source, dest deployment.Chain,
	offRamp *offramp.OffRamp,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	if len(expectedSeqNrs) == 0 {
		return nil, errors.New("no expected sequence numbers provided")
	}

	timer := time.NewTimer(tests.WaitTimeout(t))
	defer timer.Stop()
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error to subscribe ExecutionStateChanged : %w", err)
	}
	defer subscription.Unsubscribe()

	// some state to efficiently track the execution states
	// of all the expected sequence numbers.
	executionStates = make(map[uint64]int)
	seqNrsToWatch := make(map[uint64]struct{})
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = struct{}{}
	}
	for {
		select {
		case <-tick.C:
			for expectedSeqNr := range seqNrsToWatch {
				scc, executionState := GetExecutionState(t, source, dest, offRamp, expectedSeqNr)
				t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
					dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
				if executionState == EXECUTION_STATE_SUCCESS || executionState == EXECUTION_STATE_FAILURE {
					t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
						executionStateToString(executionState), dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
					executionStates[expectedSeqNr] = int(executionState)
					delete(seqNrsToWatch, expectedSeqNr)
					if len(seqNrsToWatch) == 0 {
						return executionStates, nil
					}
				}
			}
		case execEvent := <-sink:
			t.Logf("Received ExecutionStateChanged (state %s) for seqNum %d on chain %d (offramp %s) from chain %d",
				executionStateToString(execEvent.State), execEvent.SequenceNumber, dest.Selector, offRamp.Address().String(),
				source.Selector,
			)

			_, found := seqNrsToWatch[execEvent.SequenceNumber]
			if found && execEvent.SourceChainSelector == source.Selector {
				t.Logf("Received ExecutionStateChanged (state %s) on chain %d (offramp %s) from chain %d with expected sequence number %d",
					executionStateToString(execEvent.State), dest.Selector, offRamp.Address().String(), source.Selector, execEvent.SequenceNumber)
				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(seqNrsToWatch, execEvent.SequenceNumber)
				if len(seqNrsToWatch) == 0 {
					return executionStates, nil
				}
			}
		case <-timer.C:
			return nil, fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence numbers %+v",
				dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNrs)
		case subErr := <-subscription.Err():
			return nil, fmt.Errorf("subscription error: %w", subErr)
		}
	}
}

func ConfirmNoExecConsistentlyWithSeqNr(
	t *testing.T,
	source, dest deployment.Chain,
	offRamp *offramp.OffRamp,
	expectedSeqNr uint64,
	timeout time.Duration,
) {
	RequireConsistently(t, func() bool {
		scc, executionState := GetExecutionState(t, source, dest, offRamp, expectedSeqNr)
		t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
			dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
		if executionState == EXECUTION_STATE_UNTOUCHED {
			return true
		}
		t.Logf("Observed %s execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
			executionStateToString(executionState), dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
		return false
	}, timeout, 3*time.Second, "Expected no execution state change on chain %d (offramp %s) from chain %d with expected sequence number %d", dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
}

func GetExecutionState(t *testing.T, source, dest deployment.Chain, offRamp *offramp.OffRamp, expectedSeqNr uint64) (offramp.OffRampSourceChainConfig, uint8) {
	// if it's simulated backend, commit to ensure mining
	if backend, ok := source.Client.(*memory.Backend); ok {
		backend.Commit()
	}
	if backend, ok := dest.Client.(*memory.Backend); ok {
		backend.Commit()
	}
	scc, err := offRamp.GetSourceChainConfig(nil, source.Selector)
	require.NoError(t, err)
	executionState, err := offRamp.GetExecutionState(nil, source.Selector, expectedSeqNr)
	require.NoError(t, err)
	return scc, executionState
}

func RequireConsistently(t *testing.T, condition func() bool, duration time.Duration, tick time.Duration, msgAndArgs ...interface{}) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	tickTimer := time.NewTicker(tick)
	defer tickTimer.Stop()
	for {
		select {
		case <-tickTimer.C:
			if !condition() {
				require.FailNow(t, "Condition failed", msgAndArgs...)
			}
		case <-timer.C:
			return
		}
	}
}

func SeqNumberRangeToSlice(seqRanges map[SourceDestPair]ccipocr3.SeqNumRange) map[SourceDestPair][]uint64 {
	flatten := make(map[SourceDestPair][]uint64)

	for srcDst, seqRange := range seqRanges {
		if _, ok := flatten[srcDst]; !ok {
			flatten[srcDst] = make([]uint64, 0, seqRange.End()-seqRange.Start()+1)
		}

		for i := seqRange.Start(); i <= seqRange.End(); i++ {
			flatten[srcDst] = append(flatten[srcDst], uint64(i))
		}
	}

	return flatten
}

const (
	EXECUTION_STATE_UNTOUCHED  = 0
	EXECUTION_STATE_INPROGRESS = 1
	EXECUTION_STATE_SUCCESS    = 2
	EXECUTION_STATE_FAILURE    = 3
)

func executionStateToString(state uint8) string {
	switch state {
	case EXECUTION_STATE_UNTOUCHED:
		return "UNTOUCHED"
	case EXECUTION_STATE_INPROGRESS:
		return "IN_PROGRESS"
	case EXECUTION_STATE_SUCCESS:
		return "SUCCESS"
	case EXECUTION_STATE_FAILURE:
		return "FAILURE"
	default:
		return "UNKNOWN"
	}
}

func AssertEqualFeeConfig(t *testing.T, want, have fee_quoter.FeeQuoterDestChainConfig) {
	assert.Equal(t, want.DestGasOverhead, have.DestGasOverhead)
	assert.Equal(t, want.IsEnabled, have.IsEnabled)
	assert.Equal(t, want.ChainFamilySelector, have.ChainFamilySelector)
	assert.Equal(t, want.DefaultTokenDestGasOverhead, have.DefaultTokenDestGasOverhead)
	assert.Equal(t, want.DefaultTokenFeeUSDCents, have.DefaultTokenFeeUSDCents)
	assert.Equal(t, want.DefaultTxGasLimit, have.DefaultTxGasLimit)
	assert.Equal(t, want.DestGasPerPayloadByteBase, have.DestGasPerPayloadByteBase)
	assert.Equal(t, want.DestGasPerPayloadByteHigh, have.DestGasPerPayloadByteHigh)
	assert.Equal(t, want.DestGasPerPayloadByteThreshold, have.DestGasPerPayloadByteThreshold)
	assert.Equal(t, want.DestGasPerDataAvailabilityByte, have.DestGasPerDataAvailabilityByte)
	assert.Equal(t, want.DestDataAvailabilityMultiplierBps, have.DestDataAvailabilityMultiplierBps)
	assert.Equal(t, want.DestDataAvailabilityOverheadGas, have.DestDataAvailabilityOverheadGas)
	assert.Equal(t, want.MaxDataBytes, have.MaxDataBytes)
	assert.Equal(t, want.MaxNumberOfTokensPerMsg, have.MaxNumberOfTokensPerMsg)
	assert.Equal(t, want.MaxPerMsgGasLimit, have.MaxPerMsgGasLimit)
}

// AssertTimelockOwnership asserts that the ownership of the contracts has been transferred
// to the appropriate timelock contract on each chain.
func AssertTimelockOwnership(
	t *testing.T,
	e DeployedEnv,
	chains []uint64,
	state changeset.CCIPOnChainState,
) {
	// check that the ownership has been transferred correctly
	for _, chain := range chains {
		for _, contract := range []common.Address{
			state.Chains[chain].OnRamp.Address(),
			state.Chains[chain].OffRamp.Address(),
			state.Chains[chain].FeeQuoter.Address(),
			state.Chains[chain].NonceManager.Address(),
			state.Chains[chain].RMNRemote.Address(),
		} {
			owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.Chains[chain].Client)
			require.NoError(t, err)
			require.Equal(t, state.Chains[chain].Timelock.Address(), owner)
		}
	}

	// check home chain contracts ownership
	homeChainTimelockAddress := state.Chains[e.HomeChainSel].Timelock.Address()
	for _, contract := range []common.Address{
		state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
		state.Chains[e.HomeChainSel].CCIPHome.Address(),
		state.Chains[e.HomeChainSel].RMNHome.Address(),
	} {
		owner, _, err := commonchangeset.LoadOwnableContract(contract, e.Env.Chains[e.HomeChainSel].Client)
		require.NoError(t, err)
		require.Equal(t, homeChainTimelockAddress, owner)
	}
}
