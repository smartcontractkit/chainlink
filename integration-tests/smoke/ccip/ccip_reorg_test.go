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
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipcs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func Test_CCIPReorg_BelowFinality_OnSource(t *testing.T) {
}

func Test_CCIPReorg_BelowFinality_OnDest(t *testing.T) {

}

func Test_CCIPReorg_GreaterThanFinality_OnDest(t *testing.T) {
}

// This test sends a ccip message and re-orgs the chain
// after the message block has been finalized.
// The result should be that the plugin does not process
// messages from the re-orged chain anymore.
// However, it should gracefully process messages from non-reorged chains.
func Test_CCIPReorg_GreaterThanFinality_OnSource(t *testing.T) {
	require.Equal(
		t,
		os.Getenv(testhelpers.ENVTESTTYPE),
		string(testhelpers.Docker),
		"Reorg tests are only supported in docker environments",
	)

	l := logging.GetTestLogger(t)

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
			{
				Msg:    "Failed to poll and save logs due to finality violation, retrying later",
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

	l.Info().Msgf("nonBootstrapP2PIDs: %s", nonBootstrapP2PIDs)

	state, err := ccipcs.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChains := e.Env.AllChainSelectors()
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceChainSelector := allChains[0]
	sourceChain, ok := chainsel.ChainBySelector(sourceChainSelector)
	require.True(t, ok)
	logPollerService := fmt.Sprintf("EVM.%d.LogPoller", sourceChain.EvmChainID)
	l.Info().Msgf("log poller service name: %s", logPollerService)
	destChainSelector := allChains[1]

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(fmt.Sprintf("%d", chain.ChainID), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}

	sourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[sourceChainSelector], nil)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChainSelector, destChainSelector, false)

	// wait for log poller filters to get registered.
	l.Info().Msg("waiting for log poller filters to get registered")
	time.Sleep(30 * time.Second)
	msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChainSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.Info().Msgf("sent CCIP message, msgSentEvent: %+v", msgSentEvent)

	// Run reorg above finality depth
	const reorgDepth = 50
	l.Info().Int("reorgDepth", reorgDepth).Msgf("starting blockchain reorg on Simulated Geth chain")
	err = sourceClient.GethSetHead(reorgDepth)
	require.NoError(t, err, "error starting blockchain reorg on Simulated Geth chain")

	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()
	nonBootstrapCount := len(nodeAPIs) - 1
	l.Info().Msgf("waiting for %d non-bootstrap nodes to report finality violation on the logpoller", nonBootstrapCount)
	require.Eventually(t, func() bool {
		violatedResponses := make(map[string]struct{})
		for _, node := range nodeAPIs {
			// skip bootstrap nodes, they won't have any logpoller filters
			p2pKeys, err := node.MustReadP2PKeys()
			require.NoError(t, err)

			l.Debug().Msgf("got p2pKeys from node API: %+v", p2pKeys)

			require.GreaterOrEqual(t, len(p2pKeys.Data), 1)
			if !slices.Contains(nonBootstrapP2PIDs, p2pKeys.Data[0].Attributes.PeerID) {
				l.Info().Msgf("skipping bootstrap node w/ p2p id %s", p2pKeys.Data[0].Attributes.PeerID)
				continue
			}

			resp, _, err := node.Health()
			require.NoError(t, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == logPollerService &&
					d.Attributes.Output == "finality violated" &&
					d.Attributes.Status == "failing" {
					violatedResponses[p2pKeys.Data[0].Attributes.PeerID] = struct{}{}
					break
				}
			}

			if _, ok := violatedResponses[p2pKeys.Data[0].Attributes.PeerID]; ok {
				l.Info().Msgf("node %s reported finality violation", p2pKeys.Data[0].Attributes.PeerID)
			} else {
				l.Info().Msgf("node %s did not report finality violation, log poller response: %+v",
					p2pKeys.Data[0].Attributes.PeerID,
					getLogPollerHealth(logPollerService, resp.Data),
				)
			}
		}

		l.Info().Msgf("%d nodes reported finality violation", len(violatedResponses))
		return len(violatedResponses) == nonBootstrapCount
	}, 3*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	l.Info().Msg("All nodes reported finality violation")

	// TODO: assert that the plugin does not commit the message?
}

func getLogPollerHealth(logPollerService string, healthResponses []nodeclient.HealthResponseDetail) nodeclient.HealthCheck {
	for _, d := range healthResponses {
		if d.Attributes.Name == logPollerService {
			return d.Attributes
		}
	}

	return nodeclient.HealthCheck{}
}
