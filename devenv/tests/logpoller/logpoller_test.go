package logpoller

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
)

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered

// Execute both on environment with finalityTagEnabled and with finalityDepth
func TestLogPoller(t *testing.T) {
	cfg := &Config{
		General: General{
			Generator:        "looped",
			Contracts:        2,
			EventsPerTx:      4,
			FundingAmountEth: 100,
		},
		LoopedConfig: &LoopedConfig{
			ExecutionCount:    30,
			MinEmitWaitTimeMs: 400,
			MaxEmitWaitTimeMs: 600,
		},
	}
	executePollerTest(t, cfg)
}

// consistency test that registers filters after events were emitted and then triggers replay via API
// unfortunately there is no way to make sure that logs that are indexed are only picked up by replay
// and not by backup poller
// with approximate emission of 24 logs per second for ~110 seconds
// 6 filters are registered

// Execute both on environment with finalityTagEnabled and with finalityDepth
func TestLogPollerReplay(t *testing.T) {
	cfg := &Config{
		General: General{
			Generator:        "looped",
			Contracts:        2,
			EventsPerTx:      4,
			FundingAmountEth: 100,
		},
		LoopedConfig: &LoopedConfig{
			ExecutionCount:    100,
			MinEmitWaitTimeMs: 400,
			MaxEmitWaitTimeMs: 600,
		},
	}
	executeLogPollerReplay(t, cfg, "5m")
}

// TODO: adjust these values to what we are/can really run, they seem to be outdated
// consistency test with no network disruptions with approximate emission of 1000-1100 logs per second for ~110-120 seconds
// 900 filters are registered

// Execute both on environment with finalityTagEnabled and with finalityDepth
func XTestLogPollerHeavyLoad(t *testing.T) {
	t.Skip("Execute manually, when needed as it runs for a long time, remove the X from the test name to run it")
	cfg := &Config{
		General: General{
			Generator:        "looped",
			Contracts:        20,
			EventsPerTx:      30,
			FundingAmountEth: 100,
		},
		LoopedConfig: &LoopedConfig{
			ExecutionCount:    30,
			MinEmitWaitTimeMs: 400,
			MaxEmitWaitTimeMs: 600,
		},
	}

	executePollerTest(t, cfg)
}

// consistency test that introduces random disruptions by pausing Chainlink containers for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered

// Execute both on environment with finalityTagEnabled and with finalityDepth
func TestLogPollerChaosChainlinkNodes(t *testing.T) {
	cfg := &Config{
		General: General{
			Generator:        "looped",
			Contracts:        2,
			EventsPerTx:      4,
			FundingAmountEth: 100,
		},
		LoopedConfig: &LoopedConfig{
			ExecutionCount:    50,
			MinEmitWaitTimeMs: 400,
			MaxEmitWaitTimeMs: 600,
		},
		ChaosConfig: &ChaosConfig{
			ExperimentCount: 3,
			TargetComponent: "chainlink",
		},
	}

	executePollerTest(t, cfg)
}

// consistency test that introduces random disruptions by pausing Postgres container for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered

// Execute both on environment with finalityTagEnabled and with finalityDepth
func TestLogPollerChaosPostgres(t *testing.T) {
	cfg := &Config{
		General: General{
			Generator:        "looped",
			Contracts:        2,
			EventsPerTx:      4,
			FundingAmountEth: 25,
		},
		LoopedConfig: &LoopedConfig{
			ExecutionCount:    50,
			MinEmitWaitTimeMs: 400,
			MaxEmitWaitTimeMs: 600,
		},
		ChaosConfig: &ChaosConfig{
			ExperimentCount: 3,
			TargetComponent: "postgres",
		},
	}

	executePollerTest(t, cfg, products.NewAllowedLogMessage(
		"SLOW SQL QUERY",
		"It is expected, because we are pausing the Postgres container",
		zapcore.DPanicLevel,
		products.WarnAboutAllowedMsgs_No,
	))
}

func executePollerTest(t *testing.T, cfg *Config, allowedLogMessages ...products.AllowedLogMessage) {
	l := framework.L
	ctx := t.Context()

	t.Cleanup(func() {
		err := products.ScanLogs(l, products.DefaultSettings(allowedLogMessages...))
		require.NoError(t, err, "Found concerning logs in Chainlink Node logs")
	})

	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	require.Equal(t, "1337", in.Blockchains[0].ChainID, "log poller tests can only run on simulated network with chain ID 1337")
	pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
	require.NoError(t, err)

	eventsToEmit := []abi.Event{}
	for _, event := range EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}
	cfg.General.EventsToEmit = eventsToEmit
	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)

	var config *automation.Automation
	for _, candidate := range pdConfig.Config {
		if candidate.MustGetRegistryVersion() == contracts.RegistryVersion_2_1 {
			config = candidate
			break
		}
	}
	require.NotNil(t, config, "failed to find matching config with registry version 2.1")

	pks := []string{products.NetworkPrivateKey()}

	// assuming we need to have at least 2x more keys than upkeeps in order to generate the amount of events we need
	keysRequired := upKeepsNeeded * 2

	// on simulated network create new ephemeral addresses if insufficient private keys were provided
	// we ignore key at index 0, because it is the root key, which is not used during the test
	// for contract deployment and interaction
	if len(pks)-1 != keysRequired {
		bcNode := in.Blockchains[0].Out.Nodes[0]
		c, _, _, err := products.ETHClient(
			ctx,
			bcNode.ExternalWSUrl,
			config.GasSettings.FeeCapMultiplier,
			config.GasSettings.TipCapMultiplier,
		)
		require.NoError(t, err, "Failed to create ETH client")

		checkRequiredBalance(t, keysRequired, c, cfg.General.FundingAmountEth)
		newPks, err := products.FundNewAddresses(ctx, keysRequired, c, cfg.General.FundingAmountEth)
		require.NoError(t, err, "Failed to fund new addresses")
		pks = append(pks, newPks...)
	}
	require.GreaterOrEqual(t, len(pks), keysRequired+1, "you must provide at least %d private keys", keysRequired+1)

	chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
	require.NoError(t, err, "Failed to parse chain ID")

	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
	require.NoError(t, err, "Failed to create chain client")

	lpTestEnv, err := newLpTestEnvironment(chainClient, config, in)
	require.NoError(t, err, "Failed to create log poller test environment")

	err = lpTestEnv.linkToken.Transfer(config.DeployedContracts.Registry, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(upKeepsNeeded))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	_, upkeepIDs := automation.DeployLegacyConsumers(t, lpTestEnv.chainClient, lpTestEnv.registry, lpTestEnv.registrar, lpTestEnv.linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

	err = assertUpkeepIDsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := uploadLogEmitterContracts(l, t, chainClient, cfg)
	err = assertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	lpTestEnv.logEmitters = logEmitters
	lpTestEnv.upkeepIDs = upkeepIDs
	lpTestEnv.upKeepsNeeded = upKeepsNeeded

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		cl, err := clclient.New(lpTestEnv.nodes.Out.CLNodes)
		require.NoError(t, err, "failed to create chainlink clients")
		if in.Blockchains[0].ChainID != "1337" {
			_ = products.ReturnFundsFromNodes(l, chainClient, cl)
		}
	})

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = registerFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	l.Info().Msg("No duplicate filters found. OK!")

	expectedFilters := getExpectedFilters(lpTestEnv.logEmitters, cfg)
	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx, l, nil, t, lpTestEnv, expectedFilters)

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := lpTestEnv.chainClient.Client.BlockNumber(ctx)
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
		executeChaosExperiment(ctx, l, in.NodeSets[0], lpTestEnv.chainClient, cfg, chaosDoneCh)
	}()

	totalLogsEmitted, err := executeGenerator(t, cfg, lpTestEnv.chainClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")

	expectedLogsEmitted := getExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	eb, err := lpTestEnv.chainClient.Client.BlockNumber(ctx)
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

	allNodesLogCountMatches, err := checkIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	waitUntilNodesHaveTheSameLogsAsEvm(l, nil, t, allNodesLogCountMatches, lpTestEnv, cfg, startBlock, endBlock, "5m")
}

func executeLogPollerReplay(t *testing.T, cfg *Config, consistencyTimeout string) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	require.Equal(t, "1337", in.Blockchains[0].ChainID, "log poller tests can only run on simulated network with chain ID 1337")

	pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
	require.NoError(t, err)
	require.NotNil(t, pdConfig.Config[0].EVMNetworkSettings, "EVMNetworkSettings must not be nil in log poller tests")
	finalitySet := (pdConfig.Config[0].EVMNetworkSettings.FinalityTagEnabled != nil && *pdConfig.Config[0].EVMNetworkSettings.FinalityTagEnabled) || (pdConfig.Config[0].EVMNetworkSettings.FinalityDepth != nil && *pdConfig.Config[0].EVMNetworkSettings.FinalityDepth > 0)
	require.True(t, finalitySet, "Either FinalityTagEnabled must be true or FinalityDepth must be greater than 0")

	l := framework.L
	ctx := t.Context()
	t.Cleanup(func() {
		err := products.ScanLogs(l, products.DefaultSettings())
		require.NoError(t, err, "Found concerning logs in Chainlink Node logs")
	})

	eventsToEmit := []abi.Event{}
	for _, event := range EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}
	cfg.General.EventsToEmit = eventsToEmit
	upKeepsNeeded := cfg.General.Contracts * len(cfg.General.EventsToEmit)

	var config *automation.Automation
	for _, candidate := range pdConfig.Config {
		if candidate.MustGetRegistryVersion() == contracts.RegistryVersion_2_1 {
			config = candidate
			break
		}
	}
	require.NotNil(t, config, "failed to find matching config with registry version 2.1")

	pks := []string{products.NetworkPrivateKey()}

	// assuming we need to have at least 2x more keys than upkeeps in order to generate the amount of events we need
	// +1 because we won't be using the root key at index 0
	keysRequired := upKeepsNeeded*2 + 1

	// on simulated network create new ephemeral addresses if insufficient private keys were provided
	// we require +1 private keys, because key at index 0 is the root key, which is not used during the test
	// for contract deployment and interaction
	if len(pks) != keysRequired {
		bcNode := in.Blockchains[0].Out.Nodes[0]
		c, _, _, err := products.ETHClient(
			ctx,
			bcNode.ExternalWSUrl,
			config.GasSettings.FeeCapMultiplier,
			config.GasSettings.TipCapMultiplier,
		)
		require.NoError(t, err, "Failed to create ETH client")

		checkRequiredBalance(t, keysRequired, c, cfg.General.FundingAmountEth)

		newPks, err := products.FundNewAddresses(ctx, keysRequired, c, cfg.General.FundingAmountEth)
		require.NoError(t, err, "Failed to fund new addresses")
		pks = append(pks, newPks...)
	}
	require.GreaterOrEqual(t, len(pks), defaultAmountOfUpkeeps+1, "you must provide at least %d private keys", defaultAmountOfUpkeeps+1)

	chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
	require.NoError(t, err, "Failed to parse chain ID")

	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pks, &chainID)
	require.NoError(t, err, "Failed to create chain client")

	lpTestEnv, err := newLpTestEnvironment(chainClient, config, in)
	require.NoError(t, err, "Failed to create log poller test environment")

	err = lpTestEnv.linkToken.Transfer(config.DeployedContracts.Registry, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(upKeepsNeeded))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	_, upkeepIDs := automation.DeployLegacyConsumers(t, lpTestEnv.chainClient, lpTestEnv.registry, lpTestEnv.registrar, lpTestEnv.linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

	err = assertUpkeepIDsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := uploadLogEmitterContracts(l, t, chainClient, cfg)
	err = assertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	lpTestEnv.logEmitters = logEmitters
	lpTestEnv.upkeepIDs = upkeepIDs
	lpTestEnv.upKeepsNeeded = upKeepsNeeded

	cl, err := clclient.New(lpTestEnv.nodes.Out.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		if in.Blockchains[0].ChainID != "1337" {
			_ = products.ReturnFundsFromNodes(l, chainClient, cl)
		}
	})

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := chainClient.Client.BlockNumber(ctx)
	require.NoError(t, err, "Error getting latest block number")
	if sb > math.MaxInt64 {
		t.Fatalf("start block overflows int64: %d", sb)
	}
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()
	totalLogsEmitted, err := executeGenerator(t, cfg, chainClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := getExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := chainClient.Client.BlockNumber(ctx)
	require.NoError(t, err, "Error getting latest block number")

	if eb > math.MaxInt64 {
		t.Fatalf("end block overflows int64: %d", eb)
	}

	endBlock, err := getEndBlockToWaitFor(int64(eb), pdConfig.Config[0])
	require.NoError(t, err, "Error getting end block to wait for")

	require.NotZero(t, duration, "test duration cannot be zero")
	l.Info().Int64("Ending Block", endBlock).Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Lets make sure no logs are in DB yet
	expectedFilters := getExpectedFilters(lpTestEnv.logEmitters, cfg)
	logCountMatches, err := nodesHaveExpectedLogCount(startBlock, endBlock, big.NewInt(chainClient.ChainID), 0, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")
	require.True(t, logCountMatches, "Some CL nodes already had logs in DB")
	l.Info().Msg("No logs were saved by CL nodes yet, as expected. Proceeding.")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = registerFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx, l, nil, t, lpTestEnv, expectedFilters)

	blockFinalisationWaitDuration := "5m"
	l.Warn().Str("Duration", blockFinalisationWaitDuration).Msg("Waiting for all CL nodes to have end block finalised")
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFinalised, err := logPollerHasFinalisedEndBlock(endBlock, big.NewInt(chainClient.ChainID), l, nil, lpTestEnv)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if nodes have finalised end block. Retrying...")
		}
		g.Expect(hasFinalised).To(gomega.BeTrue(), "Some nodes have not finalised end block")
	}, blockFinalisationWaitDuration, "10s").Should(gomega.Succeed())

	// Trigger replay
	l.Info().Msg("Triggering log poller's replay")
	for i := 1; i < len(lpTestEnv.nodes.NodeSpecs); i++ {
		nodeName := lpTestEnv.nodes.Out.CLNodes[i].Node.ContainerName
		response, _, err := cl[i].ReplayLogPollerFromBlock(startBlock, chainClient.ChainID)
		require.NoError(t, err, "Error triggering log poller's replay on node %s", nodeName)
		require.Equal(t, "Replay started", response.Data.Attributes.Message, "Unexpected response message from log poller's replay")
	}

	// so that we don't have to look for block number of the last block in which logs were emitted as that's not trivial to do
	endBlock += 10000
	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for replay logs to be processed by all nodes")

	// logCountWaitDuration, err := time.ParseDuration("5m")
	allNodesLogCountMatches, err := checkIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	waitUntilNodesHaveTheSameLogsAsEvm(l, nil, t, allNodesLogCountMatches, lpTestEnv, cfg, startBlock, endBlock, "5m")
}

func checkRequiredBalance(t *testing.T, keysRequired int, c *ethclient.Client, fundingAmountEth float64) {
	fundingTxCount := keysRequired - 1
	requiredBalanceEth := new(big.Float).Mul(big.NewFloat(float64(fundingTxCount)), big.NewFloat(fundingAmountEth))
	requiredBalanceWeiFloat := new(big.Float).Mul(requiredBalanceEth, big.NewFloat(1e18))
	requiredBalanceWei, _ := requiredBalanceWeiFloat.Int(nil)

	gasCtx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	feeCapWei, err := c.SuggestGasPrice(gasCtx)
	cancel()
	require.NoError(t, err, "Failed to fetch gas price for funding txs")
	gasPerFundingTx := new(big.Int).Mul(big.NewInt(int64(products.DefaultNativeTransferGasPrice)), feeCapWei)
	totalFundingGasWei := new(big.Int).Mul(gasPerFundingTx, big.NewInt(int64(fundingTxCount)))
	requiredBalanceWei.Add(requiredBalanceWei, totalFundingGasWei)

	privateKeyStr := strings.TrimPrefix(products.NetworkPrivateKey(), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	require.NoError(t, err, "Failed to parse private key")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok, "error casting public key to ECDSA")
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	balanceAt, bErr := c.BalanceAt(t.Context(), fromAddress, nil)
	require.NoError(t, bErr, "Failed to get balance")
	require.GreaterOrEqual(t, balanceAt.Cmp(requiredBalanceWei), 0, "Insufficient balance. Need %s wei but have %s", requiredBalanceWei.String(), balanceAt.String())
}
