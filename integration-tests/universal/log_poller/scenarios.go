package logpoller

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	core_logger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
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
		err      error
		testName = "basic-log-poller"
	)

	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, testName, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, false, time.Duration(500*time.Millisecond), cfg.General.UseFinalityTag,
	)

	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)
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
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < cfg.General.Contracts; i++ {
		logEmitter, err := testEnv.ContractDeployer.DeployLogEmitterContract()
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
	}

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	for i := 0; i < len(upkeepIDs); i++ {
		emitterAddress := (*logEmitters[i%cfg.General.Contracts]).Address()
		upkeepID := upkeepIDs[i]
		topicId := cfg.General.EventsToEmit[i%len(cfg.General.EventsToEmit)].ID

		l.Info().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicId.Hex()).Msg("Registering log trigger for log emitter")
		err = registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicId)
		randomWait(50, 200)
		require.NoError(t, err, "Error registering log trigger for log emitter")
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Msg("Fetching filters from log poller's DB")

			hasFilters, err := nodeHasExpectedFilters(expectedFilters, coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			if err != nil {
				l.Warn().Err(err).Msg("Error checking if node has expected filters. Retrying...")
				return
			}

			g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
		}
	}, "30s", "1s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")
	l.Info().Int("Count", len(expectedFilters)).Msg("Expected filters count")

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Msg("STARTING EVENT EMISSION")
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
	l.Info().Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	// +10 to be safe, but this should be done fluently, but in such a way that can handle reorgs, so making sure that each trx is included in finalised block might not suffice
	// as some trx might be included in a block that was pruned
	endBlock := int64(eb) + 10

	l.Info().Msg("Waiting before proceeding with test until all chaos experiments finish")
	chaosError := <-chaosDoneCh
	require.NoError(t, chaosError, "Error encountered during chaos experiment")

	// Wait until last block in which events were emitted has been finalised
	// how long should we wait here until all logs are processed? wait for block X to be processed by all nodes?
	waitDuration := "15m"
	l.Warn().Str("Duration", waitDuration).Msg("Waiting for logs to be processed by all nodes and for chain to advance beyond finality")

	gom.Eventually(func(g gomega.Gomega) {
		hasAdvanced, err := chainHasAdvancedBeyondFinality(l, testEnv.EVMClient, endBlock)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if chain has advanced beyond finality. Retrying...")
		}
		g.Expect(hasAdvanced).To(gomega.BeTrue(), "Chain has not advanced beyond finality")
	}, waitDuration, "30s").Should(gomega.Succeed())

	var fibonacciBackoff = func(attempt int) int64 {
		if attempt <= 0 {
			return 0
		}
		a, b := 0, 1
		for i := 2; i <= attempt; i++ {
			a, b = b, a+b
		}

		return int64(b)
	}

	attempt := 0
	currentBackoff := int64(0)
	maxBackoff := int64(100)

	gom.Eventually(func(g gomega.Gomega) {
		logCountMatches, err := clNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), totalLogsEmitted, expectedFilters, l, coreLogger, testEnv.ClCluster)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if CL nodes have expected log count. Retrying...")
		}
		if !logCountMatches && currentBackoff < maxBackoff {
			currentBackoff = fibonacciBackoff(attempt)
			endBlock += currentBackoff
			attempt += 1
			l.Info().Int64("New end block", endBlock).Msg("Increasing end block")
		}
		g.Expect(logCountMatches).To(gomega.BeTrue(), "Not all CL nodes have expected log count")
	}, waitDuration, "5s").Should(gomega.Succeed())

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	logConsistencyWaitDuration := "1m"
	l.Warn().Str("Duration", logConsistencyWaitDuration).Msg("Waiting for CL nodes to have all the logs that EVM node has")

	gom.Eventually(func(g gomega.Gomega) {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err != nil {
			l.Warn().Err(err).Msg("Error getting missing logs. Retrying...")
		}

		if !missingLogs.IsEmpty() {
			printMissingLogsByType(missingLogs, l, cfg)
		}
		g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
	}, logConsistencyWaitDuration, "5s").Should(gomega.Succeed())
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
		err      error
		testName = "replay-log-poller"
	)

	// we set blockBackfillDepth to 0, to make sure nothing will be backfilled and won't interfere with our test
	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, testName, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, false, time.Duration(1000*time.Millisecond), cfg.General.UseFinalityTag,
	)

	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)
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
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < cfg.General.Contracts; i++ {
		logEmitter, err := testEnv.ContractDeployer.DeployLogEmitterContract()
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
	}

	time.Sleep(5 * time.Second) //wait for contracts to be uploaded to chain, TODO: could make this wait fluent

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Msg("STARTING EVENT EMISSION")
	startTime := time.Now()
	totalLogsEmitted, err := executeGenerator(t, cfg, logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := getExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())
	l.Info().Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	// +10 to be safe, but this should be done fluently, but in such a way that can handle reorgs, so making sure that each trx is included in finalised block might not suffice
	// as some trx might be included in a block that was pruned
	endBlock := int64(eb) + 10

	// Lets make sure no logs are in DB yet
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	logCountMatches, err := clNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), 0, expectedFilters, l, coreLogger, testEnv.ClCluster)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")
	require.True(t, logCountMatches, "Some CL nodes already had logs in DB")
	l.Info().Msg("No logs were saved by CL nodes yet, as expected. Proceeding.")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	for i := 0; i < len(upkeepIDs); i++ {
		emitterAddress := (*logEmitters[i%cfg.General.Contracts]).Address()
		upkeepID := upkeepIDs[i]
		topicId := cfg.General.EventsToEmit[i%len(cfg.General.EventsToEmit)].ID

		l.Info().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicId.Hex()).Msg("Registering log trigger for log emitter")
		err = registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicId)
		require.NoError(t, err, "Error registering log trigger for log emitter")
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Msg("Fetching filters from log poller's DB")

			hasFilters, err := nodeHasExpectedFilters(expectedFilters, coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			if err != nil {
				l.Warn().Err(err).Msg("Error checking if node has expected filters. Retrying...")
				return
			}

			g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
		}
	}, "30s", "1s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")
	l.Info().Int("Count", len(expectedFilters)).Msg("Expected filters count")

	// Trigger replay
	l.Info().Msg("Triggering log poller's replay")
	for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
		nodeName := testEnv.ClCluster.Nodes[i].ContainerName
		response, _, err := testEnv.ClCluster.Nodes[i].API.ReplayLogPollerFromBlock(startBlock, testEnv.EVMClient.GetChainID().Int64())
		require.NoError(t, err, "Error triggering log poller's replay on node %s", nodeName)
		require.Equal(t, "Replay started", response.Data.Attributes.Message, "Unexpected response message from log poller's replay")
	}

	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for logs to be processed by all nodes and for chain to advance beyond finality")

	gom.Eventually(func(g gomega.Gomega) {
		logCountMatches, err := clNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), totalLogsEmitted, expectedFilters, l, coreLogger, testEnv.ClCluster)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if CL nodes have expected log count. Retrying...")
		}
		g.Expect(logCountMatches).To(gomega.BeTrue(), "Not all CL nodes have expected log count")
	}, consistencyTimeout, "30s").Should(gomega.Succeed())

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for CL nodes to have all the logs that EVM node has")

	gom.Eventually(func(g gomega.Gomega) {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err != nil {
			l.Warn().Err(err).Msg("Error getting missing logs. Retrying...")
		}

		if !missingLogs.IsEmpty() {
			printMissingLogsByType(missingLogs, l, cfg)
		}
		g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
	}, consistencyTimeout, "10s").Should(gomega.Succeed())
}

type FinalityBlockFn = func(chainId int64, endBlock int64) (int64, error)

func ExecuteCILogPollerTest(t *testing.T, cfg *Config, getFinalBlockFn FinalityBlockFn) {
	l := logging.GetTestLogger(t)
	coreLogger := core_logger.TestLogger(t) //needed by ORM ¯\_(ツ)_/¯

	if cfg.General.EventsToEmit == nil || len(cfg.General.EventsToEmit) == 0 {
		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
		for _, event := range EmitterABI.Events {
			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
		}
	}

	l.Info().Msg("Starting CI log poller test")

	var (
		err      error
		testName = "ci-log-poller"
	)

	chainClient, _, contractDeployer, linkToken, registry, registrar, testEnv := setupLogPollerTestDocker(
		t, testName, ethereum.RegistryVersion_2_1, defaultOCRRegistryConfig, false, time.Duration(1000*time.Millisecond), cfg.General.UseFinalityTag,
	)

	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)
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
	logEmitters := make([]*contracts.LogEmitter, 0)
	for i := 0; i < cfg.General.Contracts; i++ {
		logEmitter, err := testEnv.ContractDeployer.DeployLogEmitterContract()
		logEmitters = append(logEmitters, &logEmitter)
		require.NoError(t, err, "Error deploying log emitter contract")
		l.Info().Str("Contract address", logEmitter.Address().Hex()).Msg("Log emitter contract deployed")
	}

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	for i := 0; i < len(upkeepIDs); i++ {
		emitterAddress := (*logEmitters[i%cfg.General.Contracts]).Address()
		upkeepID := upkeepIDs[i]
		topicId := cfg.General.EventsToEmit[i%len(cfg.General.EventsToEmit)].ID

		l.Info().Int("Upkeep id", int(upkeepID.Int64())).Str("Emitter address", emitterAddress.String()).Str("Topic", topicId.Hex()).Msg("Registering log trigger for log emitter")
		err = registerSingleTopicFilter(registry, upkeepID, emitterAddress, topicId)
		randomWait(50, 200)
		require.NoError(t, err, "Error registering log trigger for log emitter")
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error encountered when waiting for setting trigger config for upkeeps")

	// Make sure that all nodes have expected filters registered before starting to emit events
	expectedFilters := getExpectedFilters(logEmitters, cfg)
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		for i := 1; i < len(testEnv.ClCluster.Nodes); i++ {
			nodeName := testEnv.ClCluster.Nodes[i].ContainerName
			l.Info().Str("Node name", nodeName).Msg("Fetching filters from log poller's DB")

			hasFilters, err := nodeHasExpectedFilters(expectedFilters, coreLogger, testEnv.EVMClient.GetChainID(), testEnv.ClCluster.Nodes[i].PostgresDb)
			if err != nil {
				l.Warn().Err(err).Msg("Error checking if node has expected filters. Retrying...")
				return
			}

			g.Expect(hasFilters).To(gomega.BeTrue(), "Not all expected filters were found in the DB")
		}
	}, "30s", "1s").Should(gomega.Succeed())
	l.Info().Msg("All nodes have expected filters registered")
	l.Info().Int("Count", len(expectedFilters)).Msg("Expected filters count")

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	startBlock := int64(sb)

	l.Info().Msg("STARTING EVENT EMISSION")
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
	l.Info().Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := testEnv.EVMClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")

	endBlock, err := getFinalBlockFn(testEnv.EVMClient.GetChainID().Int64(), int64(eb))
	require.NoError(t, err, "Error getting final block")

	l.Info().Msg("Waiting before proceeding with test until all chaos experiments finish")
	chaosError := <-chaosDoneCh
	require.NoError(t, chaosError, "Error encountered during chaos experiment")

	// Wait until last block in which events were emitted has been finalised (with buffer)
	waitDuration := "45m"
	l.Warn().Str("Duration", waitDuration).Msg("Waiting for chain to advance beyond finality")

	gom.Eventually(func(g gomega.Gomega) {
		hasAdvanced, err := chainHasAdvancedBeyondFinality(l, testEnv.EVMClient, endBlock)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if chain has advanced beyond finality. Retrying...")
		}
		g.Expect(hasAdvanced).To(gomega.BeTrue(), "Chain has not advanced beyond finality")
	}, waitDuration, "30s").Should(gomega.Succeed())

	// Wait until all CL nodes have exactly the same logs emitted by test contracts as the EVM node has
	logConsistencyWaitDuration := "5m"
	l.Warn().Str("Duration", logConsistencyWaitDuration).Msg("Waiting for CL nodes to have all the logs that EVM node has")

	gom.Eventually(func(g gomega.Gomega) {
		missingLogs, err := getMissingLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, testEnv.ClCluster, l, coreLogger, cfg)
		if err != nil {
			l.Warn().Err(err).Msg("Error getting missing logs. Retrying...")
		}

		if !missingLogs.IsEmpty() {
			printMissingLogsByType(missingLogs, l, cfg)
		}
		g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
	}, logConsistencyWaitDuration, "20s").Should(gomega.Succeed())

	evmLogs, _ := getEVMLogs(startBlock, endBlock, logEmitters, testEnv.EVMClient, l, cfg)

	if totalLogsEmitted != len(evmLogs) {
		l.Warn().Int("Total logs emitted", totalLogsEmitted).Int("Total logs in EVM", len(evmLogs)).Msg("Test passed, but total logs emitted does not match total logs in EVM")
	}
}
