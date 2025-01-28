package ccip

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

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
		testhelpers.WithExtraConfigTomls([]string{
			t.Name() + ".toml",
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
	reorgSourceSelector := allChains[0]
	reorgSourceChain, ok := chainsel.ChainBySelector(reorgSourceSelector)
	require.True(t, ok)
	reorgLogPollerService := fmt.Sprintf("EVM.%d.LogPoller", reorgSourceChain.EvmChainID)
	l.Info().
		Msgf("reorging log poller service name: %s", reorgLogPollerService)
	destChainSelector := allChains[1]

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chain.ChainID, 10), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}

	reorgSourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[reorgSourceSelector], nil)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, reorgSourceSelector, destChainSelector, false)

	// wait for log poller filters to get registered.
	l.Info().Msg("waiting for log poller filters to get registered")
	time.Sleep(15 * time.Second)
	reorgingMsgEvent := testhelpers.TestSendRequest(t, e.Env, state, reorgSourceSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.
		Info().
		Str("messageID", hexutil.Encode(reorgingMsgEvent.Message.Header.MessageId[:])).
		Uint64("messageBlockNumber", reorgingMsgEvent.Raw.BlockNumber).
		Str("messageBlockHashBeforeReorg", reorgingMsgEvent.Raw.BlockHash.String()).
		Msg("sent CCIP message that will get re-orged before getting finalized")

	const reorgDepth = 7
	var minChainBlockNumberBeforeReorg = reorgingMsgEvent.Raw.BlockNumber + reorgDepth - 1
	// let reorgDepth - 1 blocks pass by before re-orging the message.
	// This will effectively rewind the chain to a block where the message didn't exist.
	require.Eventually(t, func() bool {
		bn, err := reorgSourceClient.BlockNumber()
		require.NoError(t, err)
		l.Info().
			Int64("blockNumber", bn).
			Uint64("targetBlockNumber", minChainBlockNumberBeforeReorg).
			Msg("Waiting for chain to progress above target block number")
		return bn >= int64(minChainBlockNumberBeforeReorg) //nolint:gosec
	}, 1*time.Minute, 500*time.Millisecond, "timeout exceeded: chain did not progress above the target block number")

	// Run reorg below finality depth
	l.Info().
		Uint64("messageBlockNumber", reorgingMsgEvent.Raw.BlockNumber).
		Int("reorgDepth", reorgDepth).
		Uint64("sourceChainSelector", reorgSourceSelector).
		Msg("starting blockchain reorg on Simulated Geth chain")
	err = reorgSourceClient.GethSetHead(reorgDepth)
	require.NoError(t, err, "error starting blockchain reorg on Simulated Geth chain")

	bnAfterReorg, err := reorgSourceClient.BlockNumber()
	require.NoError(t, err, "error getting block number after reorg")

	l.Info().Int64("blockNumberAfterReorg", bnAfterReorg).Msg("block number after reorg")

	reorgingMsgEvent = testhelpers.TestSendRequest(t, e.Env, state, reorgSourceSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.
		Info().
		Str("blockHashAfterReorg", reorgingMsgEvent.Raw.BlockHash.String()).
		Str("messageID", hexutil.Encode(reorgingMsgEvent.Message.Header.MessageId[:])).
		Uint64("messageBlockNumber", reorgingMsgEvent.Raw.BlockNumber).
		Str("messageBlockHash", reorgingMsgEvent.Raw.BlockHash.String()).
		Msgf("re-sent CCIP message after the re-org")

	nodeAPIs := dockerEnv.GetCLClusterTestEnv().ClCluster.NodeAPIs()
	nonBootstrapCount := len(nodeAPIs) - 1
	l.Info().Msgf("waiting for %d non-bootstrap nodes to NOT report finality violation on the logpoller, since re-org is less than finality", nonBootstrapCount)
	gomega.NewWithT(t).Consistently(func() bool {
		violatedResponses := make(map[string]struct{})
		for _, node := range nodeAPIs {
			// skip bootstrap nodes, they won't have any logpoller filters
			p2pKeys, err := node.MustReadP2PKeys()
			require.NoError(t, err)

			require.GreaterOrEqual(t, len(p2pKeys.Data), 1)
			if !slices.Contains(nonBootstrapP2PIDs, p2pKeys.Data[0].Attributes.PeerID) {
				continue
			}

			resp, _, err := node.Health()
			require.NoError(t, err)
			for _, d := range resp.Data {
				if d.Attributes.Name == reorgLogPollerService &&
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
					getLogPollerHealth(reorgLogPollerService, resp.Data),
				)
			}
		}

		l.Info().Msgf("%d nodes reported finality violation", len(violatedResponses))
		return len(violatedResponses) == 0
	}, time.Minute, 10*time.Second).Should(gomega.BeTrue())

	// expect the commit to still go through on the non-reorged source chain.
	_, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		reorgSourceSelector,
		e.Env.Chains[destChainSelector],
		state.Chains[destChainSelector].OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)
}

func Test_CCIPReorg_BelowFinality_OnDest(t *testing.T) {
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
		testhelpers.WithExtraConfigTomls([]string{
			t.Name() + ".toml",
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
	sourceSelector := allChains[0]
	destSelector := allChains[1]
	destChain, ok := chainsel.ChainBySelector(destSelector)
	require.True(t, ok)
	reorgLogPollerService := fmt.Sprintf("EVM.%d.LogPoller", destChain.EvmChainID)
	l.Info().
		Msgf("reorging log poller service name: %s", reorgLogPollerService)

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chain.ChainID, 10), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}

	reorgDestClient := ctf_client.NewRPCClient(chainSelToRPCURL[destSelector], nil)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destSelector, false)

	// wait for log poller filters to get registered.
	l.Info().Msg("waiting for log poller filters to get registered")
	time.Sleep(15 * time.Second)
	reorgingMsgEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceSelector, destSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.
		Info().
		Str("messageID", hexutil.Encode(reorgingMsgEvent.Message.Header.MessageId[:])).
		Uint64("messageBlockNumber", reorgingMsgEvent.Raw.BlockNumber).
		Str("messageBlockHash", reorgingMsgEvent.Raw.BlockHash.String()).
		Msgf("sent CCIP message that whose commit report will re-org on dest before getting finalized")

	// expect the commit to still go through for the message.
	reportEvent, err := testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.Chains[destSelector],
		state.Chains[destSelector].OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)

	l.
		Info().
		Uint64("reportBlockNumber", reportEvent.Raw.BlockNumber).
		Str("reportBlockHash", reportEvent.Raw.BlockHash.String()).
		Msg("got commit report on dest, preparing to re-org it")

	// re-org the dest chain less than finality blocks.
	const reorgDepth = 7
	var minChainBlockNumberBeforeReorg = reportEvent.Raw.BlockNumber + reorgDepth - 1
	// Wait for chain to progress
	require.Eventually(t, func() bool {
		bn, err := reorgDestClient.BlockNumber()
		require.NoError(t, err)
		l.Info().
			Int64("blockNumber", bn).
			Uint64("targetBlockNumber", minChainBlockNumberBeforeReorg).
			Msg("Waiting for chain to progress above target block number")
		return bn >= int64(minChainBlockNumberBeforeReorg) //nolint:gosec
	}, 1*time.Minute, 500*time.Millisecond, "timeout exceeded: chain did not progress above the target block number")

	// Run reorg below finality depth
	l.Info().
		Uint64("messageBlockNumber", reorgingMsgEvent.Raw.BlockNumber).
		Int("reorgDepth", reorgDepth).
		Uint64("sourceChainSelector", sourceSelector).
		Msg("starting blockchain reorg on Simulated Geth chain")
	err = reorgDestClient.GethSetHead(reorgDepth)
	require.NoError(t, err, "error starting blockchain reorg on Simulated Geth chain")

	bnAfterReorg, err := reorgDestClient.BlockNumber()
	require.NoError(t, err, "error getting block number after reorg")

	l.Info().Int64("blockNumberAfterReorg", bnAfterReorg).Msg("block number after reorg")

	// commit should be re-submitted after the re-org
	reportEvent, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		sourceSelector,
		e.Env.Chains[destSelector],
		state.Chains[destSelector].OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)
}

func Test_CCIPReorg_GreaterThanFinality_OnDest(t *testing.T) {
	require.Equal(
		t,
		os.Getenv(testhelpers.ENVTESTTYPE),
		string(testhelpers.Docker),
		"Reorg tests are only supported in docker environments",
	)

	l := logging.GetTestLogger(t)

	// This test sends a ccip message and re-orgs the chain
	// after the message block is finalized.
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
	sourceSelector := allChains[0]
	destChainSelector := allChains[1]
	destChain, ok := chainsel.ChainBySelector(destChainSelector)
	require.True(t, ok)
	destLogPollerService := fmt.Sprintf("EVM.%d.LogPoller", destChain.EvmChainID)
	l.
		Info().
		Msgf("reorging dest log poller service name: %s", destLogPollerService)

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chain.ChainID, 10), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}

	reorgDestClient := ctf_client.NewRPCClient(chainSelToRPCURL[destChainSelector], nil)

	// setup lanes from s1 and s2 to destChainSelector
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceSelector, destChainSelector, false)

	// wait for log poller filters to get registered.
	l.Info().Msg("waiting for log poller filters to get registered")
	time.Sleep(15 * time.Second)
	reorgingMsgEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.Info().Msgf("sent CCIP message that will get re-orged, msg id: %x", reorgingMsgEvent.Message.Header.MessageId)

	// Run reorg above finality depth
	const reorgDepth = 60
	l.Info().
		Int("reorgDepth", reorgDepth).
		Uint64("destChainSelector", destChainSelector).
		Msg("starting blockchain reorg on Simulated Geth chain")
	err = reorgDestClient.GethSetHead(reorgDepth)
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
				if d.Attributes.Name == destLogPollerService &&
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
					getLogPollerHealth(destLogPollerService, resp.Data),
				)
			}
		}

		l.Info().Msgf("%d nodes reported finality violation", len(violatedResponses))
		return len(violatedResponses) == nonBootstrapCount
	}, 2*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	l.Info().Msg("All nodes reported finality violation")

	// the commit should NOT go through on the re-orged dest chain.
	// TODO: this is not a great way to assert not-happenings.
	gomega.NewWithT(t).Consistently(func() bool {
		it, err := state.Chains[destChainSelector].OffRamp.FilterCommitReportAccepted(&bind.FilterOpts{
			Start: 0,
		})
		require.NoError(t, err)
		return !it.Next()
	}, 1*time.Minute, 10*time.Second).Should(gomega.BeTrue())
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
	// after the message block is finalized.
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
	require.GreaterOrEqual(t, len(allChains), 3)
	reorgSourceSelector := allChains[0]
	reorgSourceChain, ok := chainsel.ChainBySelector(reorgSourceSelector)
	require.True(t, ok)
	noreorgSourceSelector := allChains[1]
	noreorgSourceChain, ok := chainsel.ChainBySelector(noreorgSourceSelector)
	require.True(t, ok)
	reorgLogPollerService := fmt.Sprintf("EVM.%d.LogPoller", reorgSourceChain.EvmChainID)
	noreorgLogPollerService := fmt.Sprintf("EVM.%d.LogPoller", noreorgSourceChain.EvmChainID)
	l.Info().
		Msgf("reorging log poller service name: %s, no reorg log poller service name: %s",
			reorgLogPollerService, noreorgLogPollerService)
	destChainSelector := allChains[2]

	dockerEnv, ok := tEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok)

	chainSelToRPCURL := make(map[uint64]string)
	for _, chain := range dockerEnv.GetDevEnvConfig().Chains {
		require.GreaterOrEqual(t, len(chain.HTTPRPCs), 1)
		details, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chain.ChainID, 10), chainsel.FamilyEVM)
		require.NoError(t, err)

		chainSelToRPCURL[details.ChainSelector] = chain.HTTPRPCs[0].Internal
	}

	reorgSourceClient := ctf_client.NewRPCClient(chainSelToRPCURL[reorgSourceSelector], nil)

	// setup lanes from s1 and s2 to destChainSelector
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, reorgSourceSelector, destChainSelector, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, noreorgSourceSelector, destChainSelector, false)

	// wait for log poller filters to get registered.
	l.Info().Msg("waiting for log poller filters to get registered")
	time.Sleep(15 * time.Second)
	reorgingMsgEvent := testhelpers.TestSendRequest(t, e.Env, state, reorgSourceSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.Info().Msgf("sent CCIP message that will get re-orged, msg id: %x", reorgingMsgEvent.Message.Header.MessageId)
	msgEvent := testhelpers.TestSendRequest(t, e.Env, state, noreorgSourceSelector, destChainSelector, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[destChainSelector].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	l.Info().Msgf("sent CCIP message that will not get re-orged, msgSentEvent: %x", msgEvent.Message.Header.MessageId)

	// Run reorg above finality depth
	const reorgDepth = 50
	l.Info().
		Int("reorgDepth", reorgDepth).
		Uint64("sourceChainSelector", reorgSourceSelector).
		Msg("starting blockchain reorg on Simulated Geth chain")
	err = reorgSourceClient.GethSetHead(reorgDepth)
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
				if d.Attributes.Name == reorgLogPollerService &&
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
					getLogPollerHealth(reorgLogPollerService, resp.Data),
				)
			}
		}

		l.Info().Msgf("%d nodes reported finality violation", len(violatedResponses))
		return len(violatedResponses) == nonBootstrapCount
	}, 3*time.Minute, 5*time.Second, "not all the nodes report finality violation")
	l.Info().Msg("All nodes reported finality violation")

	// expect the commit to still go through on the non-reorged source chain.
	_, err = testhelpers.ConfirmCommitWithExpectedSeqNumRange(
		t,
		noreorgSourceSelector,
		e.Env.Chains[destChainSelector],
		state.Chains[destChainSelector].OffRamp,
		nil, // startBlock
		ccipocr3.NewSeqNumRange(1, 1),
		false, // enforceSingleCommit
	)
	require.NoError(t, err)

	// Works but super slow.
	// testhelpers.ConfirmExecWithSeqNrs(
	// 	t,
	// 	e.Env.Chains[noreorgSourceSelector],
	// 	e.Env.Chains[destChainSelector],
	// 	state.Chains[destChainSelector].OffRamp,
	// 	nil,         // startBlock
	// 	[]uint64{1}, // expectedSeqNrs
	// )
}

func getLogPollerHealth(logPollerService string, healthResponses []nodeclient.HealthResponseDetail) nodeclient.HealthCheck {
	for _, d := range healthResponses {
		if d.Attributes.Name == logPollerService {
			return d.Attributes
		}
	}

	return nodeclient.HealthCheck{}
}
