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

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

func WaitForCommitForAllWithInterval(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	expectedSeqNumRange ccipocr3.SeqNumRange,
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
				return func(src, dest uint64) error {
					var startBlock *uint64
					if startBlocks != nil {
						startBlock = startBlocks[dest]
					}
					return WaitForCommitWithInterval(t, srcChain, dstChain, state.Chains[dest].EvmOffRampV160, startBlock, expectedSeqNumRange)
				}(src, dest)
			})
		}
	}
	require.NoError(t, wg.Wait())
}

func WaitForCommitWithInterval(
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
	timer := time.NewTimer(5 * time.Minute)
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
			return fmt.Errorf("timed out waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				dest.Selector, src.Selector, expectedSeqNumRange.String())
		case report := <-sink:
			if len(report.Report.MerkleRoots) > 0 {
				// Check the interval of sequence numbers and make sure it matches
				// the expected range.
				for _, mr := range report.Report.MerkleRoots {
					if mr.SourceChainSelector == src.Selector &&
						uint64(expectedSeqNumRange.Start()) == mr.MinSeqNr &&
						uint64(expectedSeqNumRange.End()) == mr.MaxSeqNr {
						t.Logf("Received commit report on selector %d from source selector %d expected seq nr range %s",
							dest.Selector, src.Selector, expectedSeqNumRange.String())
						return nil
					}
				}
			}
		}
	}
}

func WaitForExecWithSeqNrForAll(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	expectedSeqNr uint64,
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
				return func(src, dest deployment.Chain) error {
					var startBlock *uint64
					if startBlocks != nil {
						startBlock = startBlocks[dest.Selector]
					}
					return WaitForExecWithSeqNr(t, src, dest, state.Chains[dest.Selector].EvmOffRampV160, startBlock, expectedSeqNr)
				}(srcChain, dstChain)
			})
		}
	}
	require.NoError(t, wg.Wait())
}

func WaitForExecWithSeqNr(
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
			t.Logf("Waiting for ExecutionStateChanged on chain  %d from chain %d with expected sequence number %d, current onchain minSeqNr: %d",
				dest.Selector, source.Selector, expectedSeqNr, scc.MinSeqNr)
		case execEvent := <-sink:
			if execEvent.SequenceNumber == expectedSeqNr && execEvent.SourceChainSelector == source.Selector {
				t.Logf("Received ExecutionStateChanged on chain %d from chain %d with expected sequence number %d",
					dest.Selector, source.Selector, expectedSeqNr)
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("timed out waiting for ExecutionStateChanged on chain %d from chain %d with expected sequence number %d",
				dest.Selector, source.Selector, expectedSeqNr)
		case subErr := <-subscription.Err():
			return fmt.Errorf("Subscription error: %w", subErr)
		}
	}
}
