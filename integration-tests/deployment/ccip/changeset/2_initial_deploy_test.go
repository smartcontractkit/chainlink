package changeset

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccipdeployment.Context(t)
	tenv := ccipdeployment.NewDeployedTestEnvironment(t, lggr)
	e := tenv.Env
	nodes := tenv.Nodes
	chains := e.Chains

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	// Apply migration
	output, err := Apply0002(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel: tenv.HomeChainSel,
		// Capreg/config already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(e, output.AddressBook)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	require.NoError(t, ReplayAllLogs(nodes, chains))

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
	require.NoError(t, ReplayAllLogs(nodes, chains))

	// Send a request from every router
	// Add all lanes
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				require.NoError(t, ccipdeployment.AddLane(e, state, source, dest))
			}
		}
	}

	// Send a message from each chain to every other chain.
	for src, srcChain := range e.Chains {
		for dest := range e.Chains {
			if src == dest {
				continue
			}
			msg := router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil, // TODO: no tokens for now
				FeeToken:     state.Chains[src].Weth9.Address(),
				ExtraArgs:    nil, // TODO: no extra args for now, falls back to default
			}
			fee, err := state.Chains[src].Router.GetFee(
				&bind.CallOpts{Context: context.Background()}, dest, msg)
			require.NoError(t, err, deployment.MaybeDataErr(err))
			tx, err := state.Chains[src].Weth9.Deposit(&bind.TransactOpts{
				From:   e.Chains[src].DefaultKey().From,
				Signer: e.Chains[src].DefaultKey().Signer,
				Value:  fee,
			})
			require.NoError(t, err)
			_, err = srcChain.Confirm(tx.Hash())
			require.NoError(t, err)

			// TODO: should be able to avoid this by using native?
			tx, err = state.Chains[src].Weth9.Approve(e.Chains[src].DefaultKey(),
				state.Chains[src].Router.Address(), fee)
			require.NoError(t, err)
			_, err = srcChain.Confirm(tx.Hash())
			require.NoError(t, err)

			t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
				src, dest)
			tx, err = state.Chains[src].Router.CcipSend(e.Chains[src].DefaultKey(), dest, msg)
			require.NoError(t, err)
			_, err = srcChain.Confirm(tx.Hash())
			require.NoError(t, err)
		}
	}

	// Wait for all commit reports to land.
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
				waitForCommitWithInterval(t, srcChain, dstChain, state.Chains[dest].EvmOffRampV160, ccipocr3.SeqNumRange{1, 1})
			}(src, dest)
		}
	}
	wg.Wait()

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
				waitForExecWithSeqNr(t, src, dest, state.Chains[dest.Selector].EvmOffRampV160, 1)
			}(srcChain, dstChain)
		}
	}
	wg.Wait()

	// TODO: Apply the proposal.
}

func ReplayAllLogs(nodes map[string]memory.Node, chains map[uint64]deployment.Chain) error {
	for _, node := range nodes {
		for sel := range chains {
			if err := node.ReplayLogs(map[uint64]uint64{sel: 1}); err != nil {
				return err
			}
		}
	}
	return nil
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

func waitForExecWithSeqNr(t *testing.T,
	source, dest deployment.Chain,
	offramp *offramp.OffRamp,
	expectedSeqNr uint64) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	for range tick.C {
		// TODO: Clean this up
		source.Client.(*backends.SimulatedBackend).Commit()
		dest.Client.(*backends.SimulatedBackend).Commit()
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
