package ccip

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	ctf_client "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

var (
	logsToIgnoreOpt = testhelpers.WithLogMessagesToIgnore([]testhelpers.LogMessageToIgnore{
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
	})
)

const (
	greaterThanFinalityReorgDepth = 60
	lessThanFinalityReorgDepth    = 7
)

func Test_CCIPReorg_BelowFinality_OnSource(t *testing.T) {
	e, l, dockerEnv, nonBootstrapP2PIDs, state, _ := setupReorgTest(t,
		testhelpers.WithExtraConfigTomls([]string{t.Name() + ".toml"}),
	)

	// Chain setup
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceSelector := allChains[0]
	destSelector := allChains[1]

	// Build RPC map and get clients
	chainSelToRPCURL := buildChainSelectorToRPCURLMap(t, dockerEnv)
	sourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[sourceSelector], nil)

	// Setup CCIP lane
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destSelector, false)
	waitForLogPollerFilters(l)

	// Send initial message
	msgBeforeReorg := sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)

	// Wait and perform reorg
	minBlock := msgBeforeReorg.Raw.BlockNumber + lessThanFinalityReorgDepth - 1
	waitForBlockNumber(t, sourceClient, minBlock, 1*time.Minute, 500*time.Millisecond, l)
	performReorg(t, sourceClient, lessThanFinalityReorgDepth, l)

	// Verify message consistency
	msgAfterReorg := sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)
	require.Equal(t, msgBeforeReorg.Message.Header.SequenceNumber, msgAfterReorg.Message.Header.SequenceNumber)
	require.Equal(t, msgBeforeReorg.Message.Header.MessageId, msgAfterReorg.Message.Header.MessageId)

	// Check node health
	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()
	checkFinalityViolations(
		t,
		nodeAPIs,
		nonBootstrapP2PIDs,
		getHeadTrackerService(t, sourceSelector),
		getLogPollerService(t, sourceSelector),
		l,
		0,              // no nodes reporting finality violation
		1*time.Minute,  // timeout
		10*time.Second, // interval
	)

	// Verify commit
	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.BlockChains.EVMChains()[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err, "Commit verification failed")
}

func Test_CCIPReorg_BelowFinality_OnDest(t *testing.T) {
	e, l, dockerEnv, _, state, _ := setupReorgTest(t,
		testhelpers.WithExtraConfigTomls([]string{t.Name() + ".toml"}),
	)
	evmChains := e.Env.BlockChains.EVMChains()

	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceSelector := allChains[0]
	destSelector := allChains[1]

	// Chain setup
	chainSelToRPCURL := buildChainSelectorToRPCURLMap(t, dockerEnv)
	destClient := ctf_client.NewRPCClient(chainSelToRPCURL[destSelector], nil)

	// Test setup
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destSelector, false)
	waitForLogPollerFilters(l)

	// Initial operation
	msg := sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)
	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		evmChains[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err, "Commit verification failed before reorg")

	// Reorg execution
	minBlock := msg.Raw.BlockNumber + lessThanFinalityReorgDepth - 1
	waitForBlockNumber(t, destClient, minBlock, 1*time.Minute, 500*time.Millisecond, l)
	performReorg(t, destClient, lessThanFinalityReorgDepth, l)

	// Post-reorg validation
	_, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		evmChains[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err, "Commit verification failed after reorg")
}

func Test_CCIPReorg_GreaterThanFinality_OnDest(t *testing.T) {
	e, l, dockerEnv, nonBootstrapP2PIDs, state, _ := setupReorgTest(t, logsToIgnoreOpt)

	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(allChains), 2)
	sourceSelector := allChains[0]
	destSelector := allChains[1]
	l.
		Info().
		Uint64("sourceSelector", sourceSelector).
		Uint64("destSelector", destSelector).
		Msg("Chain selectors")

	// Chain setup
	chainSelToRPCURL := buildChainSelectorToRPCURLMap(t, dockerEnv)
	destClient := ctf_client.NewRPCClient(chainSelToRPCURL[destSelector], nil)

	// Test setup
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destSelector, false)
	waitForLogPollerFilters(l)

	// Initial operation
	sendCCIPMessage(t, e.Env, state, sourceSelector, destSelector, l)
	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.BlockChains.EVMChains()[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err, "Commit verification failed before reorg")

	// Deep reorg execution
	performReorg(t, destClient, greaterThanFinalityReorgDepth, l)
	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()

	// Finality violation check
	checkFinalityViolations(
		t,
		nodeAPIs,
		nonBootstrapP2PIDs,
		getHeadTrackerService(t, destSelector),
		getLogPollerService(t, destSelector),
		l,
		len(nodeAPIs)-1,
		tests.WaitTimeout(t),
		10*time.Second,
	)

	// Commit absence check
	// TODO: timing out right away, needs to get fixed.
	// gomega.NewWithT(t).Consistently(func() bool {
	// 	it, err := state.Chains[destSelector].OffRamp.FilterCommitReportAccepted(&bind.FilterOpts{Start: 0})
	// 	require.NoError(t, err)
	// 	return !it.Next()
	// }, 1*time.Minute, 10*time.Second, tests.WaitTimeout(t)).Should(gomega.BeTrue())
}

// This test sends a ccip message and re-orgs the chain
// after the message block has been finalized.
// The result should be that the plugin does not process
// messages from the re-orged chain anymore.
// However, it should gracefully process messages from non-reorged chains.
func Test_CCIPReorg_GreaterThanFinality_OnSource(t *testing.T) {
	e, l, dockerEnv, nonBootstrapP2PIDs, state, _ := setupReorgTest(t, logsToIgnoreOpt)

	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(allChains), 3)
	reorgSource := allChains[0]
	nonReorgSource := allChains[1]
	destSelector := allChains[2]

	// Chain setup
	chainSelToRPCURL := buildChainSelectorToRPCURLMap(t, dockerEnv)
	reorgSourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[reorgSource], nil)

	// Multi-lane setup
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, reorgSource, destSelector, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, nonReorgSource, destSelector, false)
	waitForLogPollerFilters(l)

	// Send messages from both sources
	sendCCIPMessage(t, e.Env, state, reorgSource, destSelector, l)
	sendCCIPMessage(t, e.Env, state, nonReorgSource, destSelector, l)

	// Deep reorg execution
	performReorg(t, reorgSourceClient, greaterThanFinalityReorgDepth, l)
	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()

	// Finality check
	checkFinalityViolations(
		t,
		nodeAPIs,
		nonBootstrapP2PIDs,
		getHeadTrackerService(t, reorgSource),
		getLogPollerService(t, reorgSource),
		l,
		len(nodeAPIs)-1,
		tests.WaitTimeout(t),
		10*time.Second,
	)

	// Validate non-reorged source
	_, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		nonReorgSource,
		e.Env.BlockChains.EVMChains()[destSelector],
		state.MustGetEVMChainState(destSelector).OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err, "Commit verification failed for non-reorged source")

	// Commit absence check on the reorged source
	gomega.NewWithT(t).Consistently(func() bool {
		it, err := state.MustGetEVMChainState(destSelector).OffRamp.FilterCommitReportAccepted(&bind.FilterOpts{Start: 0})
		require.NoError(t, err)
		var found bool
	outer:
		for it.Next() {
			for _, mr := range it.Event.UnblessedMerkleRoots {
				if mr.SourceChainSelector == reorgSource {
					found = true
					break outer
				}
			}
		}
		return !found
	}, 1*time.Minute, 10*time.Second).Should(gomega.BeTrue())
}

func getServiceHealth(serviceName string, healthResponses []nodeclient.HealthResponseDetail) nodeclient.HealthCheck {
	for _, d := range healthResponses {
		if d.Attributes.Name == serviceName {
			return d.Attributes
		}
	}

	return nodeclient.HealthCheck{}
}

// Helper to initialize common test components
func setupReorgTest(t *testing.T, testOpts ...testhelpers.TestOps) (
	testhelpers.DeployedEnv,
	logging.Logger,
	*testsetups.DeployedLocalDevEnvironment,
	[]string,
	stateview.CCIPOnChainState,
	devenv.RMNCluster,
) {
	require.Equal(t, os.Getenv(testhelpers.ENVTESTTYPE), string(testhelpers.Docker),
		"Reorg tests are only supported in docker environments")

	l := logging.GetTestLogger(t)
	e, rmnCluster, tEnv := testsetups.NewIntegrationEnvironment(t, testOpts...)
	dockerEnv := tEnv.(*testsetups.DeployedLocalDevEnvironment)

	nodeInfos, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	nonBootstrapP2PIDs := getNonBootstrapP2PIDs(nodeInfos)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	return e, l, dockerEnv, nonBootstrapP2PIDs, state, rmnCluster
}

// Extract non-bootstrap P2P IDs
func getNonBootstrapP2PIDs(nodeInfos deployment.Nodes) []string {
	ids := make([]string, 0, len(nodeInfos.NonBootstraps()))
	for _, n := range nodeInfos.NonBootstraps() {
		ids = append(ids, strings.TrimPrefix(n.PeerID.String(), "p2p_"))
	}
	return ids
}

// Build chain selector to RPC URL map
func buildChainSelectorToRPCURLMap(t *testing.T, dockerEnv *testsetups.DeployedLocalDevEnvironment) map[uint64]string {
	devEnvConfig := dockerEnv.GetDevEnvConfig()
	require.NotNil(t, devEnvConfig)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range devEnvConfig.Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(chain.ChainID, chainsel.FamilyEVM)
		require.NoError(t, err)
		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}
	return chainSelToRPCURL
}

// Get log poller service name
func getLogPollerService(t *testing.T, chainSelector uint64) string {
	chain, ok := chainsel.ChainBySelector(chainSelector)
	require.True(t, ok)
	return fmt.Sprintf("EVM.%d.LogPoller", chain.EvmChainID)
}

func getHeadTrackerService(t *testing.T, chainSelector uint64) string {
	chain, ok := chainsel.ChainBySelector(chainSelector)
	require.True(t, ok)
	return fmt.Sprintf("EVM.%d.HeadTracker", chain.EvmChainID)
}

// Send CCIP message helper
func sendCCIPMessage(
	t *testing.T,
	env cldf.Environment,
	state stateview.CCIPOnChainState,
	sourceSelector, destSelector uint64,
	l logging.Logger,
) *onramp.OnRampCCIPMessageSent {
	msgEvent := testhelpers.TestSendRequest(t, env, state, sourceSelector, destSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.MustGetEVMChainState(destSelector).Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})

	l.Info().
		Str("messageID", hexutil.Encode(msgEvent.Message.Header.MessageId[:])).
		Uint64("blockNumber", msgEvent.Raw.BlockNumber).
		Str("blockHash", msgEvent.Raw.BlockHash.String()).
		Msg("Sent CCIP message")

	return msgEvent
}

// Check finality violations helper
// Uses require.Eventually or gomega.Consistently based on expected count
func checkFinalityViolations(
	t *testing.T,
	nodeAPIs []*nodeclient.ChainlinkClient,
	nonBootstrapP2PIDs []string,
	headTrackerServiceName,
	logPollerServiceName string,
	l logging.Logger,
	expectedCount int,
	timeout time.Duration,
	interval time.Duration,
) {
	checkFunc := func() bool {
		violated := 0
		headTrackerViolated := 0
		logPollerViolated := 0
		for _, node := range nodeAPIs {
			p2pKeys, err := node.MustReadP2PKeys()
			require.NoError(t, err)
			if len(p2pKeys.Data) == 0 || !slices.Contains(nonBootstrapP2PIDs, p2pKeys.Data[0].Attributes.PeerID) {
				l.Info().Str("p2pKey", p2pKeys.Data[0].Attributes.PeerID).Msg("Skipping bootstrap node")
				continue
			}

			resp, _, err := node.Health()
			require.NoError(t, err)
			lpHealth := getServiceHealth(logPollerServiceName, resp.Data)
			htHealth := getServiceHealth(headTrackerServiceName, resp.Data)
			var lpViolated, htViolated bool
			if lpHealth.Status == "failing" && lpHealth.Output == commontypes.ErrFinalityViolated.Error() {
				logPollerViolated++
				lpViolated = true
			}

			if htHealth.Status == "failing" && (htHealth.Output == commontypes.ErrFinalityViolated.Error() || strings.Contains(htHealth.Output, "finality violated")) {
				headTrackerViolated++
				htViolated = true
				l.Info().
					Str("p2pKey", p2pKeys.Data[0].Attributes.PeerID).
					Str("output", htHealth.Output).
					Str("headTrackerService", headTrackerServiceName).
					Msg("Head tracker finality violation")
			}

			l.Info().
				Str("p2pKey", p2pKeys.Data[0].Attributes.PeerID).
				Str("htStatus", htHealth.Status).
				Str("htOutput", htHealth.Output).
				Str("lpStatus", lpHealth.Status).
				Str("lpOutput", lpHealth.Output).
				Str("headTrackerService", headTrackerServiceName).
				Bool("headTrackerFinalityViolated", htViolated).
				Bool("logPollerFinalityViolated", lpViolated).
				Str("logPollerService", logPollerServiceName).
				Msg("finality violation checks")

			if lpViolated || htViolated {
				violated++
			}
		}

		l.
			Info().
			Int("violated", violated).
			Int("expectedCount", expectedCount).
			Int("headTrackerViolated", headTrackerViolated).
			Int("logPollerViolated", logPollerViolated).
			Msg("Checking finality violations")
		return violated == expectedCount
	}
	if expectedCount > 0 {
		require.Eventually(t, checkFunc, timeout, interval)
	} else {
		gomega.NewWithT(t).Consistently(checkFunc, timeout, interval).Should(gomega.BeTrue())
	}
}

// waitForLogPollerFilters waits for log poller filters to be registered
func waitForLogPollerFilters(l logging.Logger) {
	l.Info().Msg("Waiting for log poller filters to get registered")
	time.Sleep(30 * time.Second) // Consider making duration configurable if needed
}

// waitForBlockNumber waits until chain reaches target block number
func waitForBlockNumber(
	t *testing.T,
	client *ctf_client.RPCClient,
	targetBlock uint64,
	timeout time.Duration,
	checkInterval time.Duration,
	l logging.Logger,
) {
	require.Eventually(t, func() bool {
		bn, err := client.BlockNumber()
		require.NoError(t, err)
		l.Info().
			Int64("currentBlock", bn).
			Uint64("targetBlock", targetBlock).
			Msg("Waiting for chain progression")
		return bn >= int64(targetBlock) //nolint:gosec // no risk of overflow
	}, timeout, checkInterval, "Timeout waiting for block number")
}

// performReorg executes a chain reorg and verifies its success
func performReorg(
	t *testing.T,
	client *ctf_client.RPCClient,
	reorgDepth int,
	l logging.Logger,
) {
	l.Info().
		Int("reorgDepth", reorgDepth).
		Msg("Starting blockchain reorg")

	// Get current block before reorg for verification
	preReorgBlock, err := client.BlockNumber()
	require.NoError(t, err)

	err = client.GethSetHead(reorgDepth)
	require.NoError(t, err, "Failed to execute reorg")

	// Verify post-reorg state
	postReorgBlock, err := client.BlockNumber()
	require.NoError(t, err)

	l.Info().
		Int64("preReorgBlock", preReorgBlock).
		Int64("postReorgBlock", postReorgBlock).
		Msg("Reorg completed")

	require.Less(t, postReorgBlock, preReorgBlock,
		"Block number should decrease after reorg")
}
