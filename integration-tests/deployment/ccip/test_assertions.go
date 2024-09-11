package ccipdeployment

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
)

func WaitForCommitForAllWithInterval(
	t *testing.T,
	e deployment.Environment,
	state CCIPOnChainState,
	expectedSeqNumRange ccipocr3.SeqNumRange,
) {
	var wg sync.WaitGroup
	for src, srcChain := range e.Chains {
		for dest, dstChain := range e.Chains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Add(1)
			go func(src, dest uint64) {
				defer wg.Done()
				WaitForCommitWithInterval(t, srcChain, dstChain, state.Chains[dest].EvmOffRampV160, expectedSeqNumRange)
			}(src, dest)
		}
	}
	wg.Wait()
}

func WaitForCommitWithInterval(
	t *testing.T,
	src deployment.Chain,
	dest deployment.Chain,
	offRamp *offramp.OffRamp,
	expectedSeqNumRange ccipocr3.SeqNumRange,
) {
	sink := make(chan *offramp.OffRampCommitReportAccepted)
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: context.Background(),
	}, sink)
	require.NoError(t, err)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	//revive:disable
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
			t.Fatalf("Subscription error: %+v", subErr)
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
						return
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
) {
	var wg sync.WaitGroup
	for src, srcChain := range e.Chains {
		for dest, dstChain := range e.Chains {
			if src == dest {
				continue
			}
			srcChain := srcChain
			dstChain := dstChain
			wg.Add(1)
			go func(src, dest deployment.Chain) {
				defer wg.Done()
				WaitForExecWithSeqNr(t, src, dest, state.Chains[dest.Selector].EvmOffRampV160, expectedSeqNr)
			}(srcChain, dstChain)
		}
	}
	wg.Wait()
}

func WaitForExecWithSeqNr(t *testing.T,
	source, dest deployment.Chain,
	offramp *offramp.OffRamp,
	expectedSeqNr uint64) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	for range tick.C {
		// TODO: Clean this up
		// if it's simulated backend, commit to ensure mining
		if backend, ok := source.Client.(*backends.SimulatedBackend); ok {
			backend.Commit()
		}
		if backend, ok := dest.Client.(*backends.SimulatedBackend); ok {
			backend.Commit()
		}
		scc, err := offramp.GetSourceChainConfig(nil, source.Selector)
		require.NoError(t, err)
		t.Logf("Waiting for ExecutionStateChanged on chain  %d from chain %d with expected sequence number %d, current onchain minSeqNr: %d",
			dest.Selector, source.Selector, expectedSeqNr, scc.MinSeqNr)
		iter, err := offramp.FilterExecutionStateChanged(nil,
			[]uint64{source.Selector}, []uint64{expectedSeqNr}, nil)
		require.NoError(t, err)
		var count int
		for iter.Next() {
			if iter.Event.SequenceNumber == expectedSeqNr && iter.Event.SourceChainSelector == source.Selector {
				count++
			}
		}
		if count == 1 {
			t.Logf("Received ExecutionStateChanged on chain %d from chain %d with expected sequence number %d",
				dest.Selector, source.Selector, expectedSeqNr)
			return
		}
	}
}
