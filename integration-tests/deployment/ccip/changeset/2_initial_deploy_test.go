package changeset

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccipdeployment.Context(t)
	tenv := ccipdeployment.NewEnvironmentWithCR(t, lggr, 3)
	e := tenv.Env
	nodes := tenv.Nodes
	chains := e.Chains

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	// Apply migration
	output, err := Apply0002(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		// Capreg/config already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(e, output.AddressBook)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	require.NoError(t, ccipdeployment.ReplayAllLogs(nodes, chains))

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}
	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)

	// Ensure job related logs are up to date.
	require.NoError(t, ccipdeployment.ReplayAllLogs(nodes, chains))

	// Add all lanes
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				require.NoError(t, ccipdeployment.AddLane(e, state, source, dest))
			}
		}
	}

	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)
	for src := range e.Chains {
		for dest := range e.Chains {
			if src == dest {
				continue
			}
			seqNum := ccipdeployment.SendRequest(t, e, state, src, dest, false)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	cStart := time.Now()
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
				waitForCommitWithInterval(t, srcChain, dstChain, state.Chains[dest].OffRamp,
					ccipocr3.SeqNumRange{ccipocr3.SeqNum(expectedSeqNum[dest]), ccipocr3.SeqNum(expectedSeqNum[dest])})
			}(src, dest)
		}
	}
	wg.Wait()
	cEnd := time.Now()

	// Wait for all exec reports to land
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
				ccipdeployment.ConfirmExecution(t,
					src, dest, state.Chains[dest.Selector].OffRamp,
					expectedSeqNum[dest.Selector])
			}(srcChain, dstChain)
		}
	}
	wg.Wait()
	eEnd := time.Now()
	t.Log("Commit time:", cEnd.Sub(cStart))
	t.Log("Exec time:", eEnd.Sub(cEnd))
	// TODO: Apply the proposal.
}

func waitForCommitWithInterval(
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
			src.Client.(*backends.SimulatedBackend).Commit()
			dest.Client.(*backends.SimulatedBackend).Commit()
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
