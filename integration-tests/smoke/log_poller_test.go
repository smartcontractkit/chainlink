package smoke

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	logpoller "github.com/smartcontractkit/chainlink/integration-tests/universal/log_poller"

	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

var logScannerSettings = test_env.GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(testreporters.NewAllowedLogMessage(
	"SLOW SQL QUERY",
	"It is expected, because we are pausing the Postgres container",
	zapcore.DPanicLevel,
	testreporters.WarnAboutAllowedMsgs_No,
))

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPollerFewFiltersFixedDepth(t *testing.T) {
	executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

func TestLogPollerFewFiltersFinalityTag(t *testing.T) {
	executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

// consistency test with no network disruptions with approximate emission of 1000-1100 logs per second for ~110-120 seconds
// 900 filters are registered
func XTestLogPollerManyFiltersFixedDepth(t *testing.T) {
	t.Skip("Execute manually, when needed as it runs for a long time, remove the X from the test name to run it")
	executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

func XTestLogPollerManyFiltersFinalityTag(t *testing.T) {
	t.Skip("Execute manually, when needed as it runs for a long time, remove the X from the test name to run it")
	executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

// consistency test that introduces random disruptions by pausing either Chainlink or Postgres containers for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerWithChaosFixedDepth(t *testing.T) {
	executeBasicLogPollerTest(t, logScannerSettings)
}

func TestLogPollerWithChaosFinalityTag(t *testing.T) {
	executeBasicLogPollerTest(t, logScannerSettings)
}

func TestLogPollerWithChaosPostgresFixedDepth(t *testing.T) {
	executeBasicLogPollerTest(t, logScannerSettings)
}

func TestLogPollerWithChaosPostgresFinalityTag(t *testing.T) {
	executeBasicLogPollerTest(t, logScannerSettings)
}

// consistency test that registers filters after events were emitted and then triggers replay via API
// unfortunately there is no way to make sure that logs that are indexed are only picked up by replay
// and not by backup poller
// with approximate emission of 24 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerReplayFixedDepth(t *testing.T) {
	executeLogPollerReplay(t, "5m")
}

func TestLogPollerReplayFinalityTag(t *testing.T) {
	executeLogPollerReplay(t, "5m")
}

// HELPER FUNCTIONS
func executeBasicLogPollerTest(t *testing.T, logScannerSettings test_env.ChainlinkNodeLogScannerSettings) {
	testConfig, err := tc.GetConfig([]string{t.Name()}, tc.LogPoller)
	require.NoError(t, err, "Error getting config")
	overrideEphemeralAddressesCount(&testConfig)

	var eventsToEmit []abi.Event
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg := testConfig.LogPoller
	cfg.General.EventsToEmit = eventsToEmit

	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	lpTestEnv := prepareEnvironment(l, t, &testConfig, logScannerSettings)
	testEnv := lpTestEnv.testEnv
	sethClient := lpTestEnv.sethClient

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(testEnv.ClCluster.NodeAPIs()))
	})

	ctx := testcontext.Get(t)

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = logpoller.RegisterFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	l.Info().Msg("No duplicate filters found. OK!")

	expectedFilters := logpoller.GetExpectedFilters(lpTestEnv.logEmitters, cfg)
	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx, l, coreLogger, t, lpTestEnv, expectedFilters)

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := sethClient.Client.BlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")
	if sb > math.MaxInt64 {
		t.Fatalf("start block overflows int64: %d", sb)
	}
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()

	// Start chaos experimnents by randomly pausing random containers (Chainlink nodes or their DBs)
	chaosDoneCh := make(chan error, 1)
	go func() {
		logpoller.ExecuteChaosExperiment(l, testEnv, sethClient, &testConfig, chaosDoneCh)
	}()

	totalLogsEmitted, err := logpoller.ExecuteGenerator(t, cfg, sethClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")

	expectedLogsEmitted := logpoller.GetExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	eb, err := sethClient.Client.BlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")

	l.Info().
		Int("Total logs emitted", totalLogsEmitted).
		Uint64("Probable last block with logs", eb).
		Int64("Expected total logs emitted", expectedLogsEmitted).
		Str("Duration", fmt.Sprintf("%d sec", duration)).
		Str("LPS", fmt.Sprintf("~%d/sec", totalLogsEmitted/duration)).
		Msg("FINISHED EVENT EMISSION")

	l.Info().Msg("Waiting before proceeding with test until all chaos experiments finish")
	chaosError := <-chaosDoneCh
	require.NoError(t, chaosError, "Error encountered during chaos experiment")

	if eb > math.MaxInt64 {
		t.Fatalf("end block overflows int64: %d", eb)
	}
	// use ridciuously high end block so that we don't have to find out the block number of the last block in which logs were emitted
	// as that's not trivial to do (i.e.  just because chain was at block X when log emission ended it doesn't mean all events made it to that block)
	endBlock := int64(eb) + 10000

	allNodesLogCountMatches, err := logpoller.FluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, coreLogger, testEnv, sethClient.ChainID)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l, coreLogger, t, allNodesLogCountMatches, lpTestEnv, &testConfig, startBlock, endBlock, "5m")
}

func executeLogPollerReplay(t *testing.T, consistencyTimeout string) {
	testConfig, err := tc.GetConfig([]string{t.Name()}, tc.LogPoller)
	require.NoError(t, err, "Error getting config")
	overrideEphemeralAddressesCount(&testConfig)

	var eventsToEmit []abi.Event
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg := testConfig.LogPoller
	cfg.General.EventsToEmit = eventsToEmit

	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	lpTestEnv := prepareEnvironment(l, t, &testConfig, test_env.DefaultChainlinkNodeLogScannerSettings)
	testEnv := lpTestEnv.testEnv
	sethClient := lpTestEnv.sethClient

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(testEnv.ClCluster.NodeAPIs()))
	})

	ctx := testcontext.Get(t)
	evmNetwork, err := testEnv.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := sethClient.Client.BlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")
	if sb > math.MaxInt64 {
		t.Fatalf("start block overflows int64: %d", sb)
	}
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()
	totalLogsEmitted, err := logpoller.ExecuteGenerator(t, cfg, sethClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := logpoller.GetExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := sethClient.Client.BlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")

	if eb > math.MaxInt64 {
		t.Fatalf("end block overflows int64: %d", eb)
	}
	endBlock, err := logpoller.GetEndBlockToWaitFor(int64(eb), *evmNetwork, cfg)
	require.NoError(t, err, "Error getting end block to wait for")

	l.Info().Int64("Ending Block", endBlock).Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Lets make sure no logs are in DB yet
	expectedFilters := logpoller.GetExpectedFilters(lpTestEnv.logEmitters, cfg)
	logCountMatches, err := logpoller.ClNodesHaveExpectedLogCount(startBlock, endBlock, big.NewInt(sethClient.ChainID), 0, expectedFilters, l, coreLogger, testEnv.ClCluster)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")
	require.True(t, logCountMatches, "Some CL nodes already had logs in DB")
	l.Info().Msg("No logs were saved by CL nodes yet, as expected. Proceeding.")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = logpoller.RegisterFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx, l, coreLogger, t, lpTestEnv, expectedFilters)

	blockFinalisationWaitDuration := "5m"
	l.Warn().Str("Duration", blockFinalisationWaitDuration).Msg("Waiting for all CL nodes to have end block finalised")
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFinalised, err := logpoller.LogPollerHasFinalisedEndBlock(endBlock, big.NewInt(sethClient.ChainID), l, coreLogger, testEnv.ClCluster)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if nodes have finalised end block. Retrying...")
		}
		g.Expect(hasFinalised).To(gomega.BeTrue(), "Some nodes have not finalised end block")
	}, blockFinalisationWaitDuration, "10s").Should(gomega.Succeed())

	// Trigger replay
	l.Info().Msg("Triggering log poller's replay")
	for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
		nodeName := testEnv.ClCluster.Nodes[i].ContainerName
		response, _, err := testEnv.ClCluster.Nodes[i].API.ReplayLogPollerFromBlock(startBlock, sethClient.ChainID)
		require.NoError(t, err, "Error triggering log poller's replay on node %s", nodeName)
		require.Equal(t, "Replay started", response.Data.Attributes.Message, "Unexpected response message from log poller's replay")
	}

	// so that we don't have to look for block number of the last block in which logs were emitted as that's not trivial to do
	endBlock = endBlock + 10000
	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for replay logs to be processed by all nodes")

	// logCountWaitDuration, err := time.ParseDuration("5m")
	allNodesLogCountMatches, err := logpoller.FluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, coreLogger, testEnv, sethClient.ChainID)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l, coreLogger, t, allNodesLogCountMatches, lpTestEnv, &testConfig, startBlock, endBlock, "5m")
}

type logPollerEnvironment struct {
	sethClient    *seth.Client
	logEmitters   []*contracts.LogEmitter
	testEnv       *test_env.CLClusterTestEnv
	registry      contracts.KeeperRegistry
	upkeepIDs     []*big.Int
	upKeepsNeeded int
}

// prepareEnvironment prepares environment for log poller tests by starting DON, private Ethereum network,
// deploying registry and log emitter contracts and registering log triggered upkeeps
func prepareEnvironment(l zerolog.Logger, t *testing.T, testConfig *tc.TestConfig, logScannerSettings test_env.ChainlinkNodeLogScannerSettings) logPollerEnvironment {
	cfg := testConfig.LogPoller
	if len(cfg.General.EventsToEmit) == 0 {
		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
		for _, event := range logpoller.EmitterABI.Events {
			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
		}
	}

	l.Info().Msg("Starting basic log poller test")

	var (
		err           error
		upKeepsNeeded = *cfg.General.Contracts * len(cfg.General.EventsToEmit)
	)

	chainClient, _, linkToken, registry, registrar, testEnv, _ := logpoller.SetupLogPollerTestDocker(
		t,
		ethereum.RegistryVersion_2_1,
		logpoller.DefaultOCRRegistryConfig,
		upKeepsNeeded,
		*cfg.General.UseFinalityTag,
		testConfig,
		logScannerSettings,
	)

	_, upkeepIDs := actions.DeployLegacyConsumers(t, chainClient, registry, registrar, linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

	err = logpoller.AssertUpkeepIdsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := logpoller.UploadLogEmitterContracts(l, t, chainClient, testConfig)
	err = logpoller.AssertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	return logPollerEnvironment{
		logEmitters:   logEmitters,
		registry:      registry,
		upkeepIDs:     upkeepIDs,
		upKeepsNeeded: upKeepsNeeded,
		testEnv:       testEnv,
		sethClient:    chainClient,
	}
}

// waitForAllNodesToHaveExpectedFiltersRegisteredOrFail waits until all nodes have expected filters registered until timeout
func waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx context.Context, l zerolog.Logger, coreLogger core_logger.SugaredLogger, t *testing.T, lpTestEnv logPollerEnvironment, expectedFilters []logpoller.ExpectedFilter) {
	// Make sure that all nodes have expected filters registered before starting to emit events

	sethClient := lpTestEnv.sethClient
	testEnv := lpTestEnv.testEnv

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFilters := false
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().
				Str("Node name", nodeName).
				Msg("Fetching filters from log poller's DB")
			var message string
			var err error

			hasFilters, message, err = logpoller.NodeHasExpectedFilters(ctx, expectedFilters, coreLogger, big.NewInt(sethClient.ChainID), testEnv.ClCluster.Nodes[i].PostgresDb)
			if !hasFilters || err != nil {
				l.Warn().
					Str("Details", message).
					Msg("Some filters were missing, but we will retry")
				break
			}
		}
		g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
	}, "5m", "10s").Should(gomega.Succeed())

	l.Info().
		Msg("All nodes have expected filters registered")
	l.Info().
		Int("Count", len(expectedFilters)).
		Msg("Expected filters count")
}

// conditionallyWaitUntilNodesHaveTheSameLogsAsEvm checks whether all CL nodes have the same number of logs as EVM node
// if not, then it prints missing logs and wait for some time and checks again
func conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l zerolog.Logger, coreLogger core_logger.SugaredLogger, t *testing.T, allNodesLogCountMatches bool, lpTestEnv logPollerEnvironment, testConfig *tc.TestConfig, startBlock, endBlock int64, waitDuration string) {
	logCountWaitDuration, err := time.ParseDuration(waitDuration)
	require.NoError(t, err, "Error parsing log count wait duration")

	allNodesHaveAllExpectedLogs := false
	if !allNodesLogCountMatches {
		missingLogs, err := logpoller.GetMissingLogs(startBlock, endBlock, lpTestEnv.logEmitters, lpTestEnv.sethClient, lpTestEnv.testEnv.ClCluster, l, coreLogger, testConfig.LogPoller)
		if err == nil {
			if !missingLogs.IsEmpty() {
				logpoller.PrintMissingLogsInfo(missingLogs, l, testConfig.LogPoller)
			} else {
				allNodesHaveAllExpectedLogs = true
				l.Info().Msg("All CL nodes have all the logs that EVM node has")
			}
		}
	}

	require.True(t, allNodesLogCountMatches, "Not all CL nodes had expected log count afer %s", logCountWaitDuration)

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	// but only in the rare case that first attempt to do it failed (basically here want to know not only
	// if log count matches, but whether details of every single log match)
	if !allNodesHaveAllExpectedLogs {
		logConsistencyWaitDuration := "5m"
		l.Info().
			Str("Duration", logConsistencyWaitDuration).
			Msg("Waiting for CL nodes to have all the logs that EVM node has")

		gom := gomega.NewGomegaWithT(t)
		gom.Eventually(func(g gomega.Gomega) {
			missingLogs, err := logpoller.GetMissingLogs(startBlock, endBlock, lpTestEnv.logEmitters, lpTestEnv.sethClient, lpTestEnv.testEnv.ClCluster, l, coreLogger, testConfig.LogPoller)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("Error getting missing logs. Retrying...")
			}

			if !missingLogs.IsEmpty() {
				logpoller.PrintMissingLogsInfo(missingLogs, l, testConfig.LogPoller)
			}
			g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
		}, logConsistencyWaitDuration, "10s").Should(gomega.Succeed())
	}
}

func overrideEphemeralAddressesCount(testConfig *tc.TestConfig) {
	// override whatever is in the config file to avoid a situatiation where we don't have enough ephemeral addresses
	// to emit events from all contracts
	minContracts := int64(*testConfig.LogPoller.General.Contracts * 20)
	if testConfig.Seth.EphemeralAddrs != nil && *testConfig.Seth.EphemeralAddrs > minContracts {
		return
	}

	testConfig.Seth.EphemeralAddrs = &minContracts
}
