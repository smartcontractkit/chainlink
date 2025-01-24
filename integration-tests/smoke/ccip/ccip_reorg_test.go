package ccip

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipcs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func Test_CCIPReorgBelowFinality(t *testing.T) {
}

func Test_CCIPReorgGreaterThanFinalityOnSource(t *testing.T) {
	t.Skip("Flakey")
	// This test sends a ccip message and re-orgs the chain
	// after the message block has been finalized.
	// The result should be that the plugin does not process
	// messages from the re-orged chain anymore.
	// However, it should gracefully process messages from non-reorged chains.
	require.Equal(
		t,
		os.Getenv(testhelpers.ENVTESTTYPE),
		string(testhelpers.Docker),
		"Reorg tests are only supported in docker environments",
	)

	// This test sends a ccip message and re-orgs the chain
	// prior to the message block being finalized.
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithLogMessagesToIgnore([]testhelpers.LogMessageToIgnore{
			{
				Msg:    "Got very old block.",
				Reason: "We are expecting a re-org beyond finality",
				Level:  zapcore.DPanicLevel,
			},
			{
				Msg:    "Reorg greater than finality depth detected",
				Reason: "We are expecting a re-org beyond finality",
				Level:  zapcore.DPanicLevel,
			},
		}),
	)

	nodeInfos, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	var nonBootstrapP2PIDs = make([]string, 0, len(nodeInfos.NonBootstraps()))
	for _, n := range nodeInfos.NonBootstraps() {
		nonBootstrapP2PIDs = append(nonBootstrapP2PIDs, strings.TrimPrefix(n.PeerID.String(), "p2p_"))
	}

	t.Log("nonBootstrapP2PIDs:", nonBootstrapP2PIDs)

	state, err := ccipcs.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := e.Env.AllChainSelectors()
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceChainSelector := allChains[0]
	sourceChain, ok := chainsel.ChainBySelector(sourceChainSelector)
	require.True(t, ok)
	headTrackerService := fmt.Sprintf("EVM.%d.HeadTracker", sourceChain.EvmChainID)
	logPollerService := fmt.Sprintf("EVM.%d.LogPoller", sourceChain.EvmChainID)
	t.Log("head tracker service name:", headTrackerService, "log poller service name:", logPollerService)
	destChainSelector := allChains[1]

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(fmt.Sprintf("%d", chain.ChainID), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0]
	}

	sourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[sourceChainSelector], nil)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChainSelector, destChainSelector, false)

	// wait for log poller filters to get registered.
	time.Sleep(30 * time.Second)
	msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChainSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	t.Log("msgSentEvent:", msgSentEvent)
	time.Sleep(10 * time.Second)
	// expectedSeqNum := map[testhelpers.SourceDestPair]uint64{
	// 	{
	// 		SourceChainSelector: sourceChainSelector,
	// 		DestChainSelector:   destChainSelector,
	// 	}: msgSentEvent.SequenceNumber,
	// }
	// expectedSeqNumExec := map[testhelpers.SourceDestPair][]uint64{
	// 	{
	// 		SourceChainSelector: sourceChainSelector,
	// 		DestChainSelector:   destChainSelector,
	// 	}: {msgSentEvent.SequenceNumber},
	// }

	// Run reorg above finality depth
	t.Log("starting blockchain reorg on Simulated Geth chain")
	const reorgDepth = 50
	err = sourceClient.GethSetHead(reorgDepth)
	require.NoError(t, err, "error starting blockchain reorg on Simulated Geth chain")

	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()
	t.Logf("waiting for %d non-bootstrap nodes to report finality violation on the logpoller", len(nodeAPIs)-1)
	require.Eventually(t, func() bool {
		violatedResponses := make(map[string]struct{})
		for _, node := range nodeAPIs {
			// skip bootstrap nodes, they won't have any logpoller filters
			p2pKeys, err := node.MustReadP2PKeys()
			require.NoError(t, err)

			t.Log("p2pKeys:", p2pKeys)

			require.GreaterOrEqual(t, len(p2pKeys.Data), 1)
			if !slices.Contains(nonBootstrapP2PIDs, p2pKeys.Data[0].Attributes.PeerID) {
				t.Log("skipping bootstrap node w/ p2p id", p2pKeys.Data[0].Attributes.PeerID)
				continue
			}

			resp, _, err := node.Health()
			require.NoError(t, err)
			var violated bool
			for _, d := range resp.Data {
				if d.Attributes.Name == logPollerService &&
					d.Attributes.Output == "finality violated" &&
					d.Attributes.Status == "failing" {
					violatedResponses[p2pKeys.Data[0].Attributes.PeerID] = struct{}{}
					violated = true
					break
				}
			}

			if violated {
				t.Log("node", p2pKeys.Data[0].Attributes.PeerID, "reported finality violation")
			} else {
				t.Log("node", p2pKeys.Data[0].Attributes.PeerID, "did not report finality violation, full response: %+v", resp)
			}
		}

		t.Logf("%d nodes reported finality violation", len(violatedResponses))
		return len(violatedResponses) == (len(nodeAPIs) - 1)
	}, 3*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	t.Log("All nodes reported finality violation")
}
