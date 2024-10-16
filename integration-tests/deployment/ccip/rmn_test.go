package ccipdeployment

import (
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestRMN(t *testing.T) {
	t.Skip("Local only")

	envWithRMN, rmnCluster := NewLocalDevEnvironmentWithRMN(t, logger.TestLogger(t))
	for rmnNode, rmn := range rmnCluster.Nodes {
		t.Log(rmnNode, rmn.Proxy.PeerID, rmn.RMN.OffchainPublicKey, rmn.RMN.EVMOnchainPublicKey)
	}
	pprint(t, "envWithRmn: ", envWithRMN)

	onChainState, err := LoadOnchainState(envWithRMN.Env, envWithRMN.Ab)
	require.NoError(t, err)
	pprint(t, "onChainState", onChainState)

	// Use peerIDs to set RMN config.
	// Add a lane, send a message.

	jobSpecs, err := NewCCIPJobSpecs(envWithRMN.Env.NodeIDs, envWithRMN.Env.Offchain)
	require.NoError(t, err)

	ctx := Context(t)

	for nodeID, jobs := range jobSpecs {
		for _, job := range jobs {
			_, err := envWithRMN.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, AddLanesForAll(envWithRMN.Env, onChainState))

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)

	// Send one message from one chain to another.
	chains := maps.Values(envWithRMN.Env.Chains)
	sort.Slice(chains, func(i int, j int) bool { return chains[i].Selector < chains[j].Selector })
	srcChain := chains[0]
	dstChain := chains[1]
	require.True(t, srcChain.Selector != dstChain.Selector)
	t.Logf("source chain is %d dest chain is %d", srcChain.Selector, dstChain.Selector)

	latesthdr, err := dstChain.Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[dstChain.Selector] = &block
	seqNum := SendRequest(t, envWithRMN.Env, onChainState, srcChain.Selector, dstChain.Selector, false)
	t.Logf("expected seqNum: %d", seqNum)

	expectedSeqNum := make(map[uint64]uint64)
	expectedSeqNum[dstChain.Selector] = seqNum

	t.Logf("waiting for commit report...")
	ConfirmCommitForAllWithExpectedSeqNums(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
	t.Logf("got commit report")

	t.Logf("waiting for execute report...")
	ConfirmExecWithSeqNrForAll(t, envWithRMN.Env, onChainState, expectedSeqNum, startBlocks)
	t.Logf("got execute report")
}

func pprint(t *testing.T, msg string, v interface{}) {
	t.Logf("%s %#v", msg, v)
}
