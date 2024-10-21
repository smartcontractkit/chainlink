package ccipdeployment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

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
	require.NoError(t, wg.Wait())
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
			// TODO: Clean this up
			// if it's simulated backend, commit to ensure mining
			if backend, ok := source.Client.(*backends.SimulatedBackend); ok {
				backend.Commit()
			}
			if backend, ok := dest.Client.(*backends.SimulatedBackend); ok {
				backend.Commit()
			}
			scc, err := offRamp.GetSourceChainConfig(nil, source.Selector)
			if err != nil {
				return fmt.Errorf("error to get source chain config : %w", err)
			}
			executionState, err := offRamp.GetExecutionState(nil, source.Selector, expectedSeqNr)
			if err != nil {
				return fmt.Errorf("error to get execution state : %w", err)
			}
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
