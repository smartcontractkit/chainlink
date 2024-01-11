package logpoller

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

func ExecuteBasicLogPollerTest(t *testing.T, cfg *Config) {
	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	if cfg.General.EventsToEmit == nil || len(cfg.General.EventsToEmit) == 0 {
		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
		for _, event := range EmitterABI.Events {
			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
		}
	}

	l.Info().Msg("Starting basic log poller test")

	var (
		err           error
		upKeepsNeeded = cfg.General.Contracts * len(cfg.General.EventsToEmit)
	)

	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, upKeepsNeeded, time.Duration(500*time.Millisecond), cfg.General.UseFinalityTag,
	)

	_, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		upKeepsNeeded,
		big.NewInt(automationDefaultLinkFunds),
		automationDefaultUpkeepGasLimit,
		true,
		false,
	)

	err = assertUpkeepIdsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := uploadLogEmitterContractsAndWaitForFinalisation(l, t, testEnv, cfg)
	err = assertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = registerFiltersAndAssertUniquness(l, registry, upkeepIDs, logEmitters, cfg, upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	l.Info().Msg("No duplicate filters found. OK!")

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFilters := false
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().
				Str("Node name", nodeName).
				Msg("Fetching filters from log poller's DB")
			var message string

			hasFilters, message, err = nodeHasExpectedFilters(expectedFilters, coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			if !hasFilters {
				l.Warn().
					Str("Details", message).
					Msg("Some filters were missing, but we will retry")
				break
			}
		}
		g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
	}, "5m", "10s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")
	l.Info().
		Int("Count", len(expectedFilters)).
		Msg("Expected filters count")

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()

	// Start chaos experimnents by randomly pausing random containers (Chainlink nodes or their DBs)
	chaosDoneCh := make(chan error, 1)
	go func() {
		executeChaosExperiment(l, testEnv, cfg, chaosDoneCh)
	}()

	totalLogsEmitted, err := executeGenerator(t, cfg, logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")

	expectedLogsEmitted := getExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	eb, err := testEnv.EVMClient.LatestBlockNumber(testcontext.Get(t))
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

	// use ridciuously high end block so that we don't have to find out the block number of the last block in which logs were emitted
	// as that's not trivial to do (i.e.  just because chain was at block X when log emission ended it doesn't mean all events made it to that block)
	endBlock := int64(eb) + 10000

	logCountWaitDuration, err := time.ParseDuration("5m")
	allNodesLogCountMatches, err := fluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, coreLogger, testEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	allNodesHaveAllExpectedLogs := false
	if !allNodesLogCountMatches {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err == nil {
			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, cfg)
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

		gom.Eventually(func(g gomega.Gomega) {
			missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("Error getting missing logs. Retrying...")
			}

			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, cfg)
			}
			g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
		}, logConsistencyWaitDuration, "10s").Should(gomega.Succeed())
	}
}

func ExecuteLogPollerReplay(t *testing.T, cfg *Config, consistencyTimeout string) {
	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	if cfg.General.EventsToEmit == nil || len(cfg.General.EventsToEmit) == 0 {
		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
		for _, event := range EmitterABI.Events {
			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
		}
	}

	l.Info().Msg("Starting replay log poller test")

	var (
		err           error
		upKeepsNeeded = cfg.General.Contracts * len(cfg.General.EventsToEmit)
	)

	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, upKeepsNeeded, time.Duration(1000*time.Millisecond), cfg.General.UseFinalityTag)

	_, upkeepIDs := actions.DeployConsumers(
		t,
		registry,
		registrar,
		linkToken,
		contractDeployer,
		chainClient,
		upKeepsNeeded,
		big.NewInt(automationDefaultLinkFunds),
		automationDefaultUpkeepGasLimit,
		true,
		false,
	)

	// Deploy Log Emitter contracts
	logEmitters := uploadLogEmitterContractsAndWaitForFinalisation(l, t, testEnv, cfg)

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()
	totalLogsEmitted, err := executeGenerator(t, cfg, logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := getExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := testEnv.EVMClient.LatestBlockNumber(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest block number")

	endBlock, err := GetEndBlockToWaitFor(int64(eb), testEnv.EVMClient.GetChainID().Int64(), cfg)
	require.NoError(t, err, "Error getting end block to wait for")

	l.Info().Int64("Ending Block", endBlock).Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Lets make sure no logs are in DB yet
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	logCountMatches, err := clNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), 0, expectedFilters, l, coreLogger, testEnv.ClCluster)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")
	require.True(t, logCountMatches, "Some CL nodes already had logs in DB")
	l.Info().Msg("No logs were saved by CL nodes yet, as expected. Proceeding.")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = registerFiltersAndAssertUniquness(l, registry, upkeepIDs, logEmitters, cfg, upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFilters := false
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Msg("Fetching filters from log poller's DB")
			var message string

			hasFilters, message, err = nodeHasExpectedFilters(expectedFilters, coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			if !hasFilters {
				l.Warn().Str("Details", message).Msg("Some filters were missing, but we will retry")
				break
			}
		}
		g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
	}, "5m", "10s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")
	l.Info().Int("Count", len(expectedFilters)).Msg("Expected filters count")

	blockFinalisationWaitDuration := "5m"
	l.Warn().Str("Duration", blockFinalisationWaitDuration).Msg("Waiting for all CL nodes to have end block finalised")
	gom.Eventually(func(g gomega.Gomega) {
		hasFinalised, err := logPollerHasFinalisedEndBlock(endBlock, testEnv.EVMClient.GetChainID(), l, coreLogger, testEnv.ClCluster)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if nodes have finalised end block. Retrying...")
		}
		g.Expect(hasFinalised).To(gomega.BeTrue(), "Some nodes have not finalised end block")
	}, blockFinalisationWaitDuration, "10s").Should(gomega.Succeed())

	// Trigger replay
	l.Info().Msg("Triggering log poller's replay")
	for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
		nodeName := testEnv.ClCluster.Nodes[i].ContainerName
		response, _, err := testEnv.ClCluster.Nodes[i].API.ReplayLogPollerFromBlock(startBlock, testEnv.EVMClient.GetChainID().Int64())
		require.NoError(t, err, "Error triggering log poller's replay on node %s", nodeName)
		require.Equal(t, "Replay started", response.Data.Attributes.Message, "Unexpected response message from log poller's replay")
	}

	// so that we don't have to look for block number of the last block in which logs were emitted as that's not trivial to do
	endBlock = endBlock + 10000
	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for replay logs to be processed by all nodes")

	logCountWaitDuration, err := time.ParseDuration("5m")
	allNodesLogCountMatches, err := fluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, coreLogger, testEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	allNodesHaveAllExpectedLogs := false
	if !allNodesLogCountMatches {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err == nil {
			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, cfg)
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

		gom.Eventually(func(g gomega.Gomega) {
			missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("Error getting missing logs. Retrying...")
			}

			if !missingLogs.IsEmpty() {
				printMissingLogsInfo(missingLogs, l, cfg)
			}
			g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
		}, logConsistencyWaitDuration, "10s").Should(gomega.Succeed())
	}
}
