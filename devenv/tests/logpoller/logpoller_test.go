package logpoller

// start normal automation environment
// in the TOML config pass the pluginConfig and similar
// instead of using hardcoded values from integration-tests/actions/automation_ocr_helpers_local.go
// and BuildAutoOCR2ConfigVarsWithKeyIndexLocal() and it should work
// or maybe even try without passing these values and use default ones

// just make sure that finalityTag is correctly handled and it should be enough

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	common_logger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
)

type GeneratorType = string

const (
	GeneratorType_WASP   = "wasp"
	GeneratorType_Looped = "looped"
)

type Config struct {
	General      *General      `toml:"General"`
	ChaosConfig  *ChaosConfig  `toml:"Chaos"`
	Wasp         *WaspConfig   `toml:"Wasp"`
	LoopedConfig *LoopedConfig `toml:"Looped"`
}

func (c *Config) Validate() error {
	if c.General == nil {
		return errors.New("General config must be set")
	}

	err := c.General.Validate()
	if err != nil {
		return fmt.Errorf("General config validation failed: %w", err)
	}

	switch *c.General.Generator {
	case GeneratorType_WASP:
		if c.Wasp == nil {
			return errors.New("wasp config is nil")
		}
		err = c.Wasp.Validate()
		if err != nil {
			return fmt.Errorf("wasp config validation failed: %w", err)
		}
	case GeneratorType_Looped:
		if c.LoopedConfig == nil {
			return errors.New("looped config is nil")
		}
		err = c.LoopedConfig.Validate()
		if err != nil {
			return fmt.Errorf("looped config validation failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown generator type: %s", *c.General.Generator)
	}

	if c.ChaosConfig != nil {
		if err := c.ChaosConfig.Validate(); err != nil {
			return fmt.Errorf("chaos config validation failed: %w", err)
		}
	}

	return nil
}

type LoopedConfig struct {
	ExecutionCount    *int `toml:"execution_count"`
	MinEmitWaitTimeMs *int `toml:"min_emit_wait_time_ms"`
	MaxEmitWaitTimeMs *int `toml:"max_emit_wait_time_ms"`
}

func (l *LoopedConfig) Validate() error {
	if l.ExecutionCount == nil || *l.ExecutionCount == 0 {
		return errors.New("execution_count must be set and > 0")
	}

	if l.MinEmitWaitTimeMs == nil || *l.MinEmitWaitTimeMs == 0 {
		return errors.New("min_emit_wait_time_ms must be set and > 0")
	}

	if l.MaxEmitWaitTimeMs == nil || *l.MaxEmitWaitTimeMs == 0 {
		return errors.New("max_emit_wait_time_ms must be set and > 0")
	}

	return nil
}

type General struct {
	Generator    *string     `toml:"generator"`
	EventsToEmit []abi.Event `toml:"-"`
	Contracts    *int        `toml:"contracts"`
	EventsPerTx  *int        `toml:"events_per_tx"`
}

func (g *General) Validate() error {
	if g.Generator == nil || *g.Generator == "" {
		return errors.New("generator is empty")
	}

	if g.Contracts == nil || *g.Contracts == 0 {
		return errors.New("contracts is 0, but must be > 0")
	}

	if g.EventsPerTx == nil || *g.EventsPerTx == 0 {
		return errors.New("events_per_tx is 0, but must be > 0")
	}

	return nil
}

type ChaosConfig struct {
	ExperimentCount *int    `toml:"experiment_count"`
	TargetComponent *string `toml:"target_component"`
}

func (c *ChaosConfig) Validate() error {
	if c.ExperimentCount != nil && *c.ExperimentCount == 0 {
		return errors.New("experiment_count must be > 0")
	}

	return nil
}

type WaspConfig struct {
	RPS                   int64         `toml:"rps"`
	LPS                   int64         `toml:"lps"`
	RateLimitUnitDuration time.Duration `toml:"rate_limit_unit_duration"`
	Duration              time.Duration `toml:"duration"`
	CallTimeout           time.Duration `toml:"call_timeout"`
}

func (w *WaspConfig) Validate() error {
	if w.RPS == 0 && w.LPS == 0 {
		return errors.New("either RPS or LPS needs to be a positive integer")
	}
	if w.RPS != 0 && w.LPS != 0 {
		return errors.New("only one of RPS or LPS can be set")
	}
	if w.Duration == 0 {
		return errors.New("duration must be set and > 0")
	}
	if w.CallTimeout == 0 {
		return errors.New("call_timeout must be set and > 0")
	}
	if w.RateLimitUnitDuration == 0 {
		return errors.New("rate_limit_unit_duration  must be set and > 0")
	}

	return nil
}

// var logScannerSettings = test_env.GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(testreporters.NewAllowedLogMessage(
// 	"SLOW SQL QUERY",
// 	"It is expected, because we are pausing the Postgres container",
// 	zapcore.DPanicLevel,
// 	testreporters.WarnAboutAllowedMsgs_No,
// ))

// consistency test with no network disruptions with approximate emission of 1500-1600 logs per second for ~110-120 seconds
// 6 filters are registered
func TestLogPollerFewFiltersFixedDepth(t *testing.T) {
	executeBasicLogPollerTest(t, nil)
	// executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

func TestLogPollerFewFiltersFinalityTag(t *testing.T) {
	executeBasicLogPollerTest(t, nil)
	// executeBasicLogPollerTest(t, test_env.DefaultChainlinkNodeLogScannerSettings)
}

// consistency test with no network disruptions with approximate emission of 1000-1100 logs per second for ~110-120 seconds
// 900 filters are registered
func XTestLogPollerManyFiltersFixedDepth(t *testing.T) {
	t.Skip("Execute manually, when needed as it runs for a long time, remove the X from the test name to run it")
	executeBasicLogPollerTest(t, nil)
}

func XTestLogPollerManyFiltersFinalityTag(t *testing.T) {
	t.Skip("Execute manually, when needed as it runs for a long time, remove the X from the test name to run it")
	executeBasicLogPollerTest(t, nil)
}

// consistency test that introduces random disruptions by pausing either Chainlink or Postgres containers for random interval of 5-20 seconds
// with approximate emission of 520-550 logs per second for ~110 seconds
// 6 filters are registered
func TestLogPollerWithChaosFixedDepth(t *testing.T) {
	executeBasicLogPollerTest(t, nil)
	// executeBasicLogPollerTest(t, logScannerSettings)
}

func TestLogPollerWithChaosFinalityTag(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-562")
	// executeBasicLogPollerTest(t, logScannerSettings)
	executeBasicLogPollerTest(t, nil)
}

func TestLogPollerWithChaosPostgresFixedDepth(t *testing.T) {
	// executeBasicLogPollerTest(t, logScannerSettings)
	executeBasicLogPollerTest(t, nil)
}

func TestLogPollerWithChaosPostgresFinalityTag(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-563")
	// executeBasicLogPollerTest(t, logScannerSettings)
	executeBasicLogPollerTest(t, nil)
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
// func executeBasicLogPollerTest(t *testing.T, logScannerSettings test_env.ChainlinkNodeLogScannerSettings) {
func executeBasicLogPollerTest(t *testing.T, logScannerSettings any) {
	// testConfig, err := tc.GetConfig([]string{t.Name()}, tc.LogPoller)
	// require.NoError(t, err, "Error getting config")
	// overrideEphemeralAddressesCount(&testConfig)

	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
	require.NoError(t, err)

	// TODO pass this as argument
	cfg, err := de.LoadOutput[Config]("log_poller.toml")
	require.NoError(t, err)

	//TODO: // overrideEphemeralAddressesCount(&testConfig)

	var eventsToEmit []abi.Event
	for _, event := range EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit

	l := framework.L
	// TODO use chainlink-common?
	// coreLogger := core_logger.TestLogger(t) // needed by ORM ¯\_(ツ)_/¯

	var config *automation.Automation
	for _, candidate := range pdConfig.Config {
		if candidate.MustGetRegistryVersion() == ethereum.RegistryVersion_2_1 {
			config = candidate
			break
		}
	}
	require.NotNil(t, config, "failed to find matching config with registry version 2.1")

	upKeepsNeeded := *cfg.General.Contracts * len(cfg.General.EventsToEmit)

	chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
	require.NoError(t, err, "Failed to parse chain ID")

	var chainClient *seth.Client
	pks := []string{products.NetworkPrivateKey()}

	bcNode := in.Blockchains[0].Out.Nodes[0]
	c, _, _, err := products.ETHClient(
		t.Context(),
		bcNode.ExternalWSUrl,
		config.GasSettings.FeeCapMultiplier,
		config.GasSettings.TipCapMultiplier,
	)

	// on simulated network create new ephemeral addresses if insufficient private keys were provided
	// we require +1 private keys, because key at index 0 is the root key, which is not used during the test
	// for contract deployment and interaction
	// we create new addresses only on the simulated network to protect against fund loss
	keysRequired := upKeepsNeeded*2 + 1
	// TODO: how many addresses should we fund?
	if chainID == 1337 && len(pks) != keysRequired {
		for range keysRequired - 1 {
			address, pk, err := seth.NewAddress()
			require.NoError(t, err, "Failed to generate new address")

			cErr := products.FundNodeEIP1559(t.Context(), c, products.NetworkPrivateKey(), address, config.TestKeysMinFundingEth)
			require.NoError(t, cErr, "Failed to fund node")
			pks = append(pks, pk)
		}
	}

	require.GreaterOrEqual(t, len(pks), defaultAmountOfUpkeeps+1, "you must provide at least %d private keys", defaultAmountOfUpkeeps+1)

	if os.Getenv(seth.CONFIG_FILE_ENV_VAR) != "" {
		sethCfg, err := seth.ReadConfig()
		require.NoError(t, err, "Failed to read seth config")

		chainClient, err = seth.NewClientBuilderWithConfig(sethCfg).
			UseNetworkWithChainId(chainID).
			WithPrivateKeys(pks).
			WithRpcUrl(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl).
			Build()
	} else {
		chainClient, err = seth.NewClientBuilder().
			WithPrivateKeys(pks).
			WithRpcUrl(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl).
			Build()
	}
	require.NoError(t, err, "Failed to create chain client")

	lpTestEnv := logPollerEnvironment{
		chainClient: chainClient,
		config:      config,
		nodes:       in.NodeSets[0],
	}

	err = lpTestEnv.LoadContracts()
	require.NoError(t, err, "Failed to load contracts")

	// lpTestEnv := prepareEnvironment(l, t, &testConfig, logScannerSettings)

	err = lpTestEnv.linkToken.Transfer(config.DeployedContracts.Registry, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(upKeepsNeeded))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	// err = automation.CreateOCRKeeperJobsLocal(l, nodeClients, registry.Address(), network.ChainID, 0, registryVersion)
	// require.NoError(t, err, "Error creating OCR Keeper Jobs")
	// ocrConfig, err := automation.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, lpTestEnv.registry.RegistryOwnerAddress(), lpTestEnv.registry.ChainModuleAddress(), lpTestEnv.registry.ReorgProtectionEnabled())
	// require.NoError(t, err, "Error building OCR config vars")
	// err = lpTestEnv.registry.SetConfigTypeSafe(ocrConfig)
	// require.NoError(t, err, "Registry config should be set successfully")

	_, upkeepIDs := automation.DeployLegacyConsumers(t, lpTestEnv.chainClient, lpTestEnv.registry, lpTestEnv.registrar, lpTestEnv.linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

	err = AssertUpkeepIdsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := UploadLogEmitterContracts(l, t, chainClient, cfg)
	err = AssertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	lpTestEnv.logEmitters = logEmitters
	lpTestEnv.upkeepIDs = upkeepIDs
	lpTestEnv.upKeepsNeeded = upKeepsNeeded

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		//TODO fix
		// _ = automation.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(testEnv.ClCluster.NodeAPIs()))
	})

	ctx := t.Context()

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = RegisterFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	l.Info().Msg("No duplicate filters found. OK!")

	expectedFilters := GetExpectedFilters(lpTestEnv.logEmitters, cfg)
	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx, l, nil, t, lpTestEnv, expectedFilters)

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := lpTestEnv.chainClient.Client.BlockNumber(t.Context())
	require.NoError(t, err, "Error getting latest block number")
	if sb > math.MaxInt64 {
		t.Fatalf("start block overflows int64: %d", sb)
	}
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()

	// Start chaos experimnents by randomly pausing random containers (Chainlink nodes or their DBs)
	// chaosDoneCh := make(chan error, 1)
	// go func() {
	// TODO fix me
	// ExecuteChaosExperiment(l, testEnv, lpTestEnv.chainClient, &testConfig, chaosDoneCh)
	// }()

	totalLogsEmitted, err := ExecuteGenerator(t, cfg, lpTestEnv.chainClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")

	expectedLogsEmitted := GetExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	eb, err := lpTestEnv.chainClient.Client.BlockNumber(t.Context())
	require.NoError(t, err, "Error getting latest block number")

	l.Info().
		Int("Total logs emitted", totalLogsEmitted).
		Uint64("Probable last block with logs", eb).
		Int64("Expected total logs emitted", expectedLogsEmitted).
		Str("Duration", fmt.Sprintf("%d sec", duration)).
		Str("LPS", fmt.Sprintf("~%d/sec", totalLogsEmitted/duration)).
		Msg("FINISHED EVENT EMISSION")

	l.Info().Msg("Waiting before proceeding with test until all chaos experiments finish")
	// chaosError := <-chaosDoneCh
	// require.NoError(t, chaosError, "Error encountered during chaos experiment")

	if eb > math.MaxInt64 {
		t.Fatalf("end block overflows int64: %d", eb)
	}
	// use ridciuously high end block so that we don't have to find out the block number of the last block in which logs were emitted
	// as that's not trivial to do (i.e.  just because chain was at block X when log emission ended it doesn't mean all events made it to that block)
	endBlock := int64(eb) + 10000

	allNodesLogCountMatches, err := FluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l, nil, t, allNodesLogCountMatches, lpTestEnv, cfg, startBlock, endBlock, "5m")
}

func executeLogPollerReplay(t *testing.T, consistencyTimeout string) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[automation.Configurator](outputFile)
	require.NoError(t, err)

	// TODO pass this as argument
	cfg, err := de.LoadOutput[Config]("log_replay.toml")
	require.NoError(t, err)

	var eventsToEmit []abi.Event
	for _, event := range EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit
	upKeepsNeeded := *cfg.General.Contracts * len(cfg.General.EventsToEmit)

	l := framework.L
	// TODO use chainlink-common?
	// coreLogger := core_logger.TestLogger(t) // needed by ORM ¯\_(ツ)_/¯

	var config *automation.Automation
	for _, candidate := range pdConfig.Config {
		if candidate.MustGetRegistryVersion() == ethereum.RegistryVersion_2_1 {
			config = candidate
			break
		}
	}
	require.NotNil(t, config, "failed to find matching config with registry version 2.1")

	chainID, err := strconv.ParseUint(in.Blockchains[0].ChainID, 10, 64)
	require.NoError(t, err, "Failed to parse chain ID")

	var chainClient *seth.Client
	pks := []string{products.NetworkPrivateKey()}

	bcNode := in.Blockchains[0].Out.Nodes[0]
	c, _, _, err := products.ETHClient(
		t.Context(),
		bcNode.ExternalWSUrl,
		config.GasSettings.FeeCapMultiplier,
		config.GasSettings.TipCapMultiplier,
	)

	// on simulated network create new ephemeral addresses if insufficient private keys were provided
	// we require +1 private keys, because key at index 0 is the root key, which is not used during the test
	// for contract deployment and interaction
	// we create new addresses only on the simulated network to protect against fund loss
	keysRequired := upKeepsNeeded*2 + 1
	// TODO: how many addresses should we fund?
	if chainID == 1337 && len(pks) != keysRequired {
		for range keysRequired - 1 {
			address, pk, err := seth.NewAddress()
			require.NoError(t, err, "Failed to generate new address")

			cErr := products.FundNodeEIP1559(t.Context(), c, products.NetworkPrivateKey(), address, config.TestKeysMinFundingEth)
			require.NoError(t, cErr, "Failed to fund node")
			pks = append(pks, pk)
		}
	}

	require.GreaterOrEqual(t, len(pks), defaultAmountOfUpkeeps+1, "you must provide at least %d private keys", defaultAmountOfUpkeeps+1)

	if os.Getenv(seth.CONFIG_FILE_ENV_VAR) != "" {
		sethCfg, err := seth.ReadConfig()
		require.NoError(t, err, "Failed to read seth config")

		chainClient, err = seth.NewClientBuilderWithConfig(sethCfg).
			UseNetworkWithChainId(chainID).
			WithPrivateKeys(pks).
			WithRpcUrl(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl).
			Build()
	} else {
		chainClient, err = seth.NewClientBuilder().
			WithPrivateKeys(pks).
			WithRpcUrl(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl).
			Build()
	}
	require.NoError(t, err, "Failed to create chain client")

	lpTestEnv := logPollerEnvironment{
		chainClient: chainClient,
		config:      config,
		nodes:       in.NodeSets[0],
	}

	err = lpTestEnv.LoadContracts()
	require.NoError(t, err, "Failed to load contracts")

	// lpTestEnv := prepareEnvironment(l, t, &testConfig, logScannerSettings)

	err = lpTestEnv.linkToken.Transfer(config.DeployedContracts.Registry, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(upKeepsNeeded))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	// err = automation.CreateOCRKeeperJobsLocal(l, nodeClients, registry.Address(), network.ChainID, 0, registryVersion)
	// require.NoError(t, err, "Error creating OCR Keeper Jobs")
	// ocrConfig, err := automation.BuildAutoOCR2ConfigVarsLocal(l, workerNodes, registryConfig, registrar.Address(), 30*time.Second, lpTestEnv.registry.RegistryOwnerAddress(), lpTestEnv.registry.ChainModuleAddress(), lpTestEnv.registry.ReorgProtectionEnabled())
	// require.NoError(t, err, "Error building OCR config vars")
	// err = lpTestEnv.registry.SetConfigTypeSafe(ocrConfig)
	// require.NoError(t, err, "Registry config should be set successfully")

	_, upkeepIDs := automation.DeployLegacyConsumers(t, lpTestEnv.chainClient, lpTestEnv.registry, lpTestEnv.registrar, lpTestEnv.linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

	err = AssertUpkeepIdsUniqueness(upkeepIDs)
	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
	l.Info().Msg("No duplicate upkeep IDs found. OK!")

	// Deploy Log Emitter contracts
	logEmitters := UploadLogEmitterContracts(l, t, chainClient, cfg)
	err = AssertContractAddressUniquneness(logEmitters)
	require.NoError(t, err, "Error asserting contract addresses uniqueness")
	l.Info().Msg("No duplicate contract addresses found. OK!")

	lpTestEnv.logEmitters = logEmitters
	lpTestEnv.upkeepIDs = upkeepIDs
	lpTestEnv.upKeepsNeeded = upKeepsNeeded

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		// TODO: fix
		// _ = automation.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(testEnv.ClCluster.NodeAPIs()))
	})

	// Save block number before starting to emit events, so that we can later use it when querying logs
	sb, err := chainClient.Client.BlockNumber(t.Context())
	require.NoError(t, err, "Error getting latest block number")
	if sb > math.MaxInt64 {
		t.Fatalf("start block overflows int64: %d", sb)
	}
	startBlock := int64(sb)

	l.Info().Int64("Starting Block", startBlock).Msg("STARTING EVENT EMISSION")
	startTime := time.Now()
	totalLogsEmitted, err := ExecuteGenerator(t, cfg, chainClient, lpTestEnv.logEmitters)
	endTime := time.Now()
	require.NoError(t, err, "Error executing event generator")
	expectedLogsEmitted := GetExpectedLogCount(cfg)
	duration := int(endTime.Sub(startTime).Seconds())

	// Save block number after finishing to emit events, so that we can later use it when querying logs
	eb, err := chainClient.Client.BlockNumber(t.Context())
	require.NoError(t, err, "Error getting latest block number")

	if eb > math.MaxInt64 {
		t.Fatalf("end block overflows int64: %d", eb)
	}

	require.NotNil(t, pdConfig.Config[0].EVMNetworkSettings, "EVMNetworkSettings must not be nil in log poller tests")

	endBlock, err := GetEndBlockToWaitFor(int64(eb), int64(*pdConfig.Config[0].EVMNetworkSettings.FinalityDepth), *pdConfig.Config[0])
	require.NoError(t, err, "Error getting end block to wait for")

	require.NotZero(t, duration, "test duration cannot be zero")
	l.Info().Int64("Ending Block", endBlock).Int("Total logs emitted", totalLogsEmitted).Int64("Expected total logs emitted", expectedLogsEmitted).Str("Duration", fmt.Sprintf("%d sec", duration)).Str("LPS", fmt.Sprintf("%d/sec", totalLogsEmitted/duration)).Msg("FINISHED EVENT EMISSION")

	// Lets make sure no logs are in DB yet
	expectedFilters := GetExpectedFilters(lpTestEnv.logEmitters, cfg)
	logCountMatches, err := ClNodesHaveExpectedLogCount(startBlock, endBlock, big.NewInt(chainClient.ChainID), 0, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")
	require.True(t, logCountMatches, "Some CL nodes already had logs in DB")
	l.Info().Msg("No logs were saved by CL nodes yet, as expected. Proceeding.")

	// Register log triggered upkeep for each combination of log emitter contract and event signature (topic)
	// We need to register a separate upkeep for each event signature, because log trigger doesn't support multiple topics (even if log poller does)
	err = RegisterFiltersAndAssertUniquness(l, lpTestEnv.registry, lpTestEnv.upkeepIDs, lpTestEnv.logEmitters, cfg, lpTestEnv.upKeepsNeeded)
	require.NoError(t, err, "Error registering filters")

	waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(t.Context(), l, nil, t, lpTestEnv, expectedFilters)

	blockFinalisationWaitDuration := "5m"
	l.Warn().Str("Duration", blockFinalisationWaitDuration).Msg("Waiting for all CL nodes to have end block finalised")
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFinalised, err := LogPollerHasFinalisedEndBlock(endBlock, big.NewInt(chainClient.ChainID), l, nil, lpTestEnv)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if nodes have finalised end block. Retrying...")
		}
		g.Expect(hasFinalised).To(gomega.BeTrue(), "Some nodes have not finalised end block")
	}, blockFinalisationWaitDuration, "10s").Should(gomega.Succeed())

	// Trigger replay

	cl, err := clclient.New(lpTestEnv.nodes.Out.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	l.Info().Msg("Triggering log poller's replay")
	for i := 1; i < len(lpTestEnv.nodes.NodeSpecs); i++ {
		nodeName := lpTestEnv.nodes.Out.CLNodes[i].Node.ContainerName
		response, _, err := ReplayLogPollerFromBlock(l, cl[i], startBlock, chainClient.ChainID)
		require.NoError(t, err, "Error triggering log poller's replay on node %s", nodeName)
		require.Equal(t, "Replay started", response.Data.Attributes.Message, "Unexpected response message from log poller's replay")
	}

	// so that we don't have to look for block number of the last block in which logs were emitted as that's not trivial to do
	endBlock += 10000
	l.Warn().Str("Duration", consistencyTimeout).Msg("Waiting for replay logs to be processed by all nodes")

	// logCountWaitDuration, err := time.ParseDuration("5m")
	allNodesLogCountMatches, err := FluentlyCheckIfAllNodesHaveLogCount("5m", startBlock, endBlock, totalLogsEmitted, expectedFilters, l, nil, lpTestEnv)
	require.NoError(t, err, "Error checking if CL nodes have expected log count")

	conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l, nil, t, allNodesLogCountMatches, lpTestEnv, cfg, startBlock, endBlock, "5m")
}

// TODO: move to CTF
func ReplayLogPollerFromBlock(l zerolog.Logger, c *clclient.ChainlinkClient, fromBlock, evmChainID int64) (*ReplayResponse, *http.Response, error) {
	specObj := &ReplayResponse{}
	l.Info().Str("NodeURL", c.Config.URL).Int64("From block", fromBlock).Int64("EVM chain ID", evmChainID).Msg("Replaying Log Poller from block")
	resp, err := c.APIClient.R().
		SetResult(&specObj).
		SetQueryParams(map[string]string{
			"family":  "evm",
			"ChainID": strconv.FormatInt(evmChainID, 10),
		}).
		SetPathParams(map[string]string{
			"fromBlock": strconv.FormatInt(fromBlock, 10),
		}).
		Post("/v2/replay_from_block/{fromBlock}")
	if err != nil {
		return nil, nil, err
	}

	return specObj, resp.RawResponse, err
}

type ReplayResponse struct {
	Data ReplayResponseData `json:"data"`
}

type ReplayResponseData struct {
	Attributes ReplayResponseAttributes `json:"attributes"`
}

type ReplayResponseAttributes struct {
	Message    string   `json:"message"`
	EVMChainID *big.Int `json:"evmChainID"`
}

// OCR2ExportKey is the model that represents the exported VRF key
type OCR2ExportKey struct {
	KeyType           string `json:"keyType"`
	ChainType         string `json:"chainType"`
	ID                string `json:"id"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffchainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
	Crypto            struct {
		Cipher       string `json:"cipher"`
		Ciphertext   string `json:"ciphertext"`
		Cipherparams struct {
			Iv string `json:"iv"`
		} `json:"cipherparams"`
		Kdf       string `json:"kdf"`
		Kdfparams struct {
			Dklen int    `json:"dklen"`
			N     int    `json:"n"`
			P     int    `json:"p"`
			R     int    `json:"r"`
			Salt  string `json:"salt"`
		} `json:"kdfparams"`
		Mac string `json:"mac"`
	} `json:"crypto"`
}

type logPollerEnvironment struct {
	logger zerolog.Logger
	nodes  *nodeset.Input

	chainClient *seth.Client

	config *automation.Automation

	logEmitters   []*contracts.LogEmitter
	upkeepIDs     []*big.Int
	upKeepsNeeded int

	registry  contracts.KeeperRegistry
	registrar contracts.KeeperRegistrar
	linkToken contracts.LinkToken
}

func (l *logPollerEnvironment) dbPort() int {
	if l.nodes.DbInput.Port != 0 {
		return l.nodes.DbInput.Port
	}

	return postgres.ExposedStaticPort
}

func (l *logPollerEnvironment) LoadContracts() error {
	if err := l.LoadLINK(l.config.DeployedContracts.LinkToken); err != nil {
		return fmt.Errorf("error loading link token contract: %w", err)
	}

	if err := l.LoadRegistry(l.config.DeployedContracts.Registry, l.config.DeployedContracts.ChainModule); err != nil {
		return fmt.Errorf("error loading registry contract: %w", err)
	}

	if l.registry.RegistryOwnerAddress().String() != l.chainClient.MustGetRootKeyAddress().String() {
		return fmt.Errorf("registry owner address is not the root key address")
	}

	if err := l.LoadRegistrar(l.config.DeployedContracts.Registrar); err != nil {
		return fmt.Errorf("error loading registrar contract: %w", err)
	}

	return nil
}

func (l *logPollerEnvironment) LoadLINK(address string) error {
	linkToken, err := contracts.LoadLinkTokenContract(l.logger, l.chainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	l.linkToken = linkToken
	l.logger.Info().Str("LINK Token Address", l.linkToken.Address()).Msg("Successfully loaded LINK Token")
	return nil
}

func (l *logPollerEnvironment) LoadRegistry(registryAddress, chainModuleAddress string) error {
	registry, err := contracts.LoadKeeperRegistry(l.logger, l.chainClient, common.HexToAddress(registryAddress), ethereum.RegistryVersion_2_1, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return err
	}
	l.registry = registry
	l.logger.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", l.registry.Address()).Msg("Successfully loaded Registry")
	return nil
}

func (l *logPollerEnvironment) LoadRegistrar(address string) error {
	if l.registry == nil {
		return errors.New("registry must be deployed or loaded before registrar")
	}
	// l.RegistrarSettings.RegistryAddr = l.registry.Address()
	registrar, err := contracts.LoadKeeperRegistrar(l.chainClient, common.HexToAddress(address), ethereum.RegistryVersion_2_1)
	if err != nil {
		return err
	}
	l.logger.Info().Str("Registrar Address", registrar.Address()).Msg("Successfully loaded Registrar")
	l.registrar = registrar
	return nil
}

// prepareEnvironment prepares environment for log poller tests by starting DON, private Ethereum network,
// deploying registry and log emitter contracts and registering log triggered upkeeps
// func prepareEnvironment(l zerolog.Logger, t *testing.T, testConfig *tc.TestConfig, logScannerSettings test_env.ChainlinkNodeLogScannerSettings) logPollerEnvironment {
// 	cfg := testConfig.LogPoller
// 	if len(cfg.General.EventsToEmit) == 0 {
// 		l.Warn().Msg("No events to emit specified, using all events from log emitter contract")
// 		for _, event := range logpoller.EmitterABI.Events {
// 			cfg.General.EventsToEmit = append(cfg.General.EventsToEmit, event)
// 		}
// 	}

// 	l.Info().Msg("Starting basic log poller test")

// 	var (
// 		err           error
// 		upKeepsNeeded = *cfg.General.Contracts * len(cfg.General.EventsToEmit)
// 	)

// 	chainClient, _, linkToken, registry, registrar, testEnv, _ := logpoller.SetupLogPollerTestDocker(
// 		t,
// 		ethereum.RegistryVersion_2_1,
// 		logpoller.DefaultOCRRegistryConfig,
// 		upKeepsNeeded,
// 		*cfg.General.UseFinalityTag,
// 		testConfig,
// 		logScannerSettings,
// 	)

// 	_, upkeepIDs := actions.DeployLegacyConsumers(t, chainClient, registry, registrar, linkToken, upKeepsNeeded, big.NewInt(int64(9e18)), uint32(2500000), true, false, false, nil)

// 	err = logpoller.AssertUpkeepIdsUniqueness(upkeepIDs)
// 	require.NoError(t, err, "Error asserting upkeep ids uniqueness")
// 	l.Info().Msg("No duplicate upkeep IDs found. OK!")

// 	// Deploy Log Emitter contracts
// 	logEmitters := logpoller.UploadLogEmitterContracts(l, t, chainClient, testConfig)
// 	err = logpoller.AssertContractAddressUniquneness(logEmitters)
// 	require.NoError(t, err, "Error asserting contract addresses uniqueness")
// 	l.Info().Msg("No duplicate contract addresses found. OK!")

// 	return logPollerEnvironment{
// 		logEmitters:   logEmitters,
// 		registry:      registry,
// 		upkeepIDs:     upkeepIDs,
// 		upKeepsNeeded: upKeepsNeeded,
// 		testEnv:       testEnv,
// 		sethClient:    chainClient,
// 	}
// }

// waitForAllNodesToHaveExpectedFiltersRegisteredOrFail waits until all nodes have expected filters registered until timeout
func waitForAllNodesToHaveExpectedFiltersRegisteredOrFail(ctx context.Context, l zerolog.Logger, coreLogger common_logger.SugaredLogger, t *testing.T, testEnv logPollerEnvironment, expectedFilters []ExpectedFilter) {
	// Make sure that all nodes have expected filters registered before starting to emit events

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		hasFilters := false
		for i := 1; i < len(testEnv.nodes.NodeSpecs); i++ {
			nodeName := testEnv.nodes.Out.CLNodes[i].Node.ContainerName
			l.Info().
				Str("Node name", nodeName).
				Msg("Fetching filters from log poller's DB")
			var message string
			var err error

			hasFilters, message, err = NodeHasExpectedFilters(ctx, expectedFilters, coreLogger, big.NewInt(testEnv.chainClient.ChainID), i, testEnv.dbPort())
			if !hasFilters || err != nil {
				if message == "" {
					message = err.Error()
				}
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
func conditionallyWaitUntilNodesHaveTheSameLogsAsEvm(l zerolog.Logger, coreLogger common_logger.SugaredLogger, t *testing.T, allNodesLogCountMatches bool, lpTestEnv logPollerEnvironment, config *Config, startBlock, endBlock int64, waitDuration string) {
	logCountWaitDuration, err := time.ParseDuration(waitDuration)
	require.NoError(t, err, "Error parsing log count wait duration")

	allNodesHaveAllExpectedLogs := false
	if !allNodesLogCountMatches {
		missingLogs, err := GetMissingLogs(startBlock, endBlock, lpTestEnv, l, coreLogger, config)
		if err == nil {
			if !missingLogs.IsEmpty() {
				PrintMissingLogsInfo(missingLogs, l, config)
			} else {
				allNodesHaveAllExpectedLogs = true
				l.Info().Msg("All CL nodes have all the logs that EVM node has")
			}
		}
	}

	require.True(t, allNodesLogCountMatches, "Not all CL nodes had expected log count after %s", logCountWaitDuration)

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
			missingLogs, err := GetMissingLogs(startBlock, endBlock, lpTestEnv, l, coreLogger, config)
			if err != nil {
				l.Warn().
					Err(err).
					Msg("Error getting missing logs. Retrying...")
			}

			if !missingLogs.IsEmpty() {
				PrintMissingLogsInfo(missingLogs, l, config)
			}
			g.Expect(missingLogs.IsEmpty()).To(gomega.BeTrue(), "Some CL nodes were missing logs")
		}, logConsistencyWaitDuration, "10s").Should(gomega.Succeed())
	}
}

// func overrideEphemeralAddressesCount(config *Config) {
// 	// override whatever is in the config file to avoid a situatiation where we don't have enough ephemeral addresses
// 	// to emit events from all contracts
// 	minContracts := int64(*config.General.Contracts * 20)
// 	if config.Seth.EphemeralAddrs != nil && *config.Seth.EphemeralAddrs > minContracts {
// 		return
// 	}

// 	config.Seth.EphemeralAddrs = &minContracts
// }
