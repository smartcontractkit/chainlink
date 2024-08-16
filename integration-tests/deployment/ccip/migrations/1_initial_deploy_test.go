package migrations

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

func Test0001_InitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := Context(t)
	chains := memory.NewMemoryChains(t, 3)
	homeChainSel := uint64(0)
	homeChainEVM := uint64(0)
	// First chain is home chain.
	for chainSel := range chains {
		homeChainEVM, _ = chainsel.ChainIdFromSelector(chainSel)
		homeChainSel = chainSel
		break
	}
	ab, err := ccipdeployment.DeployCapReg(lggr, chains, homeChainSel)
	require.NoError(t, err)

	addrs, err := ab.AddressesForChain(homeChainSel)
	require.NoError(t, err)
	require.Len(t, addrs, 2)
	capReg := common.Address{}
	for addr := range addrs {
		capReg = common.HexToAddress(addr)
		break
	}
	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, memory.RegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	})
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	state, err := ccipdeployment.GenerateOnchainState(e, ab)
	require.NoError(t, err)

	capabilities, err := state.CapabilityRegistry[homeChainSel].GetCapabilities(nil)
	require.NoError(t, err)
	require.Len(t, capabilities, 1)
	ccipCap, err := state.CapabilityRegistry[homeChainSel].GetHashedCapabilityId(nil,
		ccipdeployment.CapabilityLabelledName, ccipdeployment.CapabilityVersion)
	require.NoError(t, err)
	_, err = state.CapabilityRegistry[homeChainSel].GetCapability(nil, ccipCap)
	require.NoError(t, err)

	// Apply migration
	output, err := Apply0001(e, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel: homeChainSel,
		// Capreg/config already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccipdeployment.GenerateOnchainState(e, output.AddressBook)
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
				Receiver:     common.LeftPadBytes(state.Receivers[dest].Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil, // TODO: no tokens for now
				FeeToken:     state.Weth9s[src].Address(),
				ExtraArgs:    nil, // TODO: no extra args for now, falls back to default
			}
			fee, err := state.Routers[src].GetFee(
				&bind.CallOpts{Context: context.Background()}, dest, msg)
			require.NoError(t, err, deployment.MaybeDataErr(err))
			tx, err := state.Weth9s[src].Deposit(&bind.TransactOpts{
				From:   e.Chains[src].DeployerKey.From,
				Signer: e.Chains[src].DeployerKey.Signer,
				Value:  fee,
			})
			require.NoError(t, err)
			require.NoError(t, srcChain.Confirm(tx.Hash()))

			// TODO: should be able to avoid this by using native?
			tx, err = state.Weth9s[src].Approve(e.Chains[src].DeployerKey,
				state.Routers[src].Address(), fee)
			require.NoError(t, err)
			require.NoError(t, srcChain.Confirm(tx.Hash()))

			t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
				src, dest)
			tx, err = state.Routers[src].CcipSend(e.Chains[src].DeployerKey, dest, msg)
			require.NoError(t, err)
			require.NoError(t, srcChain.Confirm(tx.Hash()))
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
				waitForCommitWithInterval(t, srcChain, dstChain, state.EvmOffRampsV160[dest], ccipocr3.SeqNumRange{1, 1})
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
			go func(src, dest uint64) {
				defer wg.Done()
				waitForExecWithSeqNr(t, srcChain, dstChain, state.EvmOffRampsV160[dest], 1)
			}(src, dest)
		}
	}
	wg.Wait()

	// TODO: Apply the proposal.
}

func ReplayAllLogs(nodes map[string]memory.Node, chains map[uint64]deployment.Chain) error {
	for _, node := range nodes {
		for sel := range chains {
			chainID, _ := chainsel.ChainIdFromSelector(sel)
			if err := node.App.ReplayFromBlock(big.NewInt(int64(chainID)), 1, false); err != nil {
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
	offRamp *evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp,
	expectedSeqNumRange ccipocr3.SeqNumRange,
) {
	sink := make(chan *evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReportAccepted)
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
						uint64(expectedSeqNumRange.Start()) == mr.Interval.Min &&
						uint64(expectedSeqNumRange.End()) == mr.Interval.Max {
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
	offramp *evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp,
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
