package ccipdeployment

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

func ConfirmGasPriceUpdatedForAll(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
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
					startBlock = startBlocks[srcChain.Selector]
				}
				return ConfirmGasPriceUpdated(
					t,
					dstChain,
					state.Chains[srcChain.Selector].FeeQuoter,
					*startBlock,
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
) error {
	it, err := srcFeeQuoter.FilterUsdPerUnitGasUpdated(&bind.FilterOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, []uint64{dest.Selector})

	require.NoError(t, err)
	require.Truef(t, it.Next(), "No gas price update event found on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	require.NotEqualf(t, InitialGasPrice, it.Event.Value, "Gas price not updated on chain %d, fee quoter %s",
		dest.Selector, srcFeeQuoter.Address().String())
	return nil
}

func ConfirmTokenPriceUpdatedForAll(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	startBlocks map[uint64]*uint64,
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
			tokenToPrice[linkAddress] = InitialLinkPrice
			tokenToPrice[wethAddress] = InitialWethPrice
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
		return fmt.Errorf("Not all tokens updated on chain  %d", chain.Selector)
	}

	return nil
}

// ConfirmCommitForAllWithExpectedSeqNums waits for all chains in the environment to commit the given expectedSeqNums.
// expectedSeqNums is a map of destinationchain selector to expected sequence number
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmCommitForAllWithExpectedSeqNums(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	expectedSeqNums map[uint64]uint64,
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

				if expectedSeqNums[dstChain.Selector] == 0 {
					return nil
				}

				return ConfirmCommitWithExpectedSeqNumRange(
					t,
					srcChain,
					dstChain,
					state.Chains[dstChain.Selector].OffRamp,
					startBlock,
					ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(expectedSeqNums[dstChain.Selector]),
						ccipocr3.SeqNum(expectedSeqNums[dstChain.Selector]),
					})
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
		3*time.Minute,
		1*time.Second,
		"all commitments did not confirm",
	)
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
) error {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink)
	if err != nil {
		return fmt.Errorf("error to subscribe CommitReportAccepted : %w", err)
	}

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
			if backend, ok := src.Client.(*backends.SimulatedBackend); ok {
				backend.Commit()
			}
			if backend, ok := dest.Client.(*backends.SimulatedBackend); ok {
				backend.Commit()
			}
			t.Logf("Waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, src.Selector, expectedSeqNumRange.String())
		case subErr := <-subscription.Err():
			return fmt.Errorf("subscription error: %w", subErr)
		case <-timer.C:
			return fmt.Errorf("timed out after waiting %s duration for commit report on chain selector %d from source selector %d expected seq nr range %s",
				duration.String(), dest.Selector, src.Selector, expectedSeqNumRange.String())
		case report := <-sink:
			if len(report.MerkleRoots) > 0 {
				// Check the interval of sequence numbers and make sure it matches
				// the expected range.
				for _, mr := range report.MerkleRoots {
					if mr.SourceChainSelector == src.Selector &&
						uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
						uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
						t.Logf("Received commit report for [%d, %d] on selector %d from source selector %d expected seq nr range %s, token prices: %v",
							mr.MinSeqNr, mr.MaxSeqNr, dest.Selector, src.Selector, expectedSeqNumRange.String(), report.PriceUpdates.TokenPriceUpdates)
						return nil
					}
				}
			}
		}
	}
}

// ConfirmExecWithSeqNrForAll waits for all chains in the environment to execute the given expectedSeqNums.
// expectedSeqNums is a map of destinationchain selector to expected sequence number
// startBlocks is a map of destination chain selector to start block number to start watching from.
// If startBlocks is nil, it will start watching from the latest block.
func ConfirmExecWithSeqNrForAll(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	expectedSeqNums map[uint64]uint64,
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

				if expectedSeqNums[dstChain.Selector] == 0 {
					return nil
				}

				return ConfirmExecWithSeqNr(
					t,
					srcChain,
					dstChain,
					state.Chains[dstChain.Selector].OffRamp,
					startBlock,
					expectedSeqNums[dstChain.Selector],
				)
			})
		}
	}
	require.NoError(t, wg.Wait())
}

// ConfirmExecWithSeqNr waits for an execution state change on the destination chain with the expected sequence number.
// startBlock is the block number to start watching from.
// If startBlock is nil, it will start watching from the latest block.
func ConfirmExecWithSeqNr(
	t *testing.T,
	source, dest deployment.Chain,
	offRamp *offramp.OffRamp,
	startBlock *uint64,
	expectedSeqNr uint64,
) error {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	sink := make(chan *offramp.OffRampExecutionStateChanged)
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: context.Background(),
		Start:   startBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error to subscribe ExecutionStateChanged : %w", err)
	}
	defer subscription.Unsubscribe()
	for {
		select {
		case <-tick.C:
			scc, executionState := GetExecutionState(t, source, dest, offRamp, expectedSeqNr)
			t.Logf("Waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d, current onchain minSeqNr: %d, execution state: %s",
				dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr, scc.MinSeqNr, executionStateToString(executionState))
			if executionState == EXECUTION_STATE_SUCCESS {
				t.Logf("Observed SUCCESS execution state on chain %d (offramp %s) from chain %d with expected sequence number %d",
					dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
				return nil
			}
		case execEvent := <-sink:
			t.Logf("Received ExecutionStateChanged for seqNum %d on chain %d (offramp %s) from chain %d",
				execEvent.SequenceNumber, dest.Selector, offRamp.Address().String(), source.Selector)
			if execEvent.SequenceNumber == expectedSeqNr && execEvent.SourceChainSelector == source.Selector {
				t.Logf("Received ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d",
					dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d (offramp %s) from chain %d with expected sequence number %d",
				dest.Selector, offRamp.Address().String(), source.Selector, expectedSeqNr)
		case subErr := <-subscription.Err():
			return fmt.Errorf("subscription error: %w", subErr)
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
	if backend, ok := source.Client.(*backends.SimulatedBackend); ok {
		backend.Commit()
	}
	if backend, ok := dest.Client.(*backends.SimulatedBackend); ok {
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
