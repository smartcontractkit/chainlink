package ccipdeployment

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

type DeployedTestEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	Nodes        map[string]memory.Node
}

// NewDeployedEnvironment creates a new CCIP environment
// with capreg and nodes set up.
func NewEnvironmentWithCR(t *testing.T, lggr logger.Logger, numChains int) DeployedTestEnvironment {
	ctx := Context(t)
	chains := memory.NewMemoryChains(t, numChains)
	// Lower chainSel is home chain.
	var chainSels []uint64
	// Say first chain is home chain.
	for chainSel := range chains {
		chainSels = append(chainSels, chainSel)
	}
	sort.Slice(chainSels, func(i, j int) bool {
		return chainSels[i] < chainSels[j]
	})
	// Take lowest for determinism.
	homeChainSel := chainSels[0]
	homeChainEVM, _ := chainsel.ChainIdFromSelector(homeChainSel)
	ab, capReg, err := DeployCapReg(lggr, chains, homeChainSel)
	require.NoError(t, err)

	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, memory.RegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	})
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedTestEnvironment{
		Ab:           ab,
		Env:          e,
		HomeChainSel: homeChainSel,
		Nodes:        nodes,
	}
}

func NewEnvironmentWithCRAndJobs(t *testing.T, lggr logger.Logger, numChains int) DeployedTestEnvironment {
	ctx := Context(t)
	e := NewEnvironmentWithCR(t, lggr, numChains)
	jbs, err := NewCCIPJobSpecs(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for nodeID, jobs := range jbs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
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
	require.NoError(t, ReplayAllLogs(e.Nodes, e.Env.Chains))
	return e
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

func SendRequest(t *testing.T, e deployment.Environment, state CCIPOnChainState, src, dest uint64, testRouter bool) uint64 {
	msg := router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello"),
		TokenAmounts: nil, // TODO: no tokens for now
		// Pay native.
		FeeToken:  common.HexToAddress("0x0"),
		ExtraArgs: nil, // TODO: no extra args for now, falls back to default
	}
	router := state.Chains[src].Router
	if testRouter {
		router = state.Chains[src].TestRouter
	}
	fee, err := router.GetFee(
		&bind.CallOpts{Context: context.Background()}, dest, msg)
	require.NoError(t, err, deployment.MaybeDataErr(err))

	t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
		src, dest)
	e.Chains[src].DeployerKey.Value = fee
	tx, err := router.CcipSend(
		e.Chains[src].DeployerKey,
		dest,
		msg)
	require.NoError(t, err)
	blockNum, err := e.Chains[src].Confirm(tx)
	require.NoError(t, err)
	it, err := state.Chains[src].OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{dest})
	require.NoError(t, err)
	require.True(t, it.Next())
	return it.Event.Message.Header.SequenceNumber
}

func ConfirmExecution(t *testing.T,
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
