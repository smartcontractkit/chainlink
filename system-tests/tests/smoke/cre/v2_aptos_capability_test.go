package cre

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	aptoslib "github.com/aptos-labs/aptos-go-sdk"
	aptoscrypto "github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	aptosbind "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptosdatafeeds "github.com/smartcontractkit/chainlink-aptos/bindings/data_feeds"
	aptosplatformsecondary "github.com/smartcontractkit/chainlink-aptos/bindings/platform_secondary"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	blockchains_aptos "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	blockchains_evm "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	aptoswrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite/config"
	aptoswriteroundtrip_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswriteroundtrip/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	aptosLocalMaxGasAmount        uint64 = 200_000
	aptosLocalGasUnitPrice        uint64 = 100
	aptosWorkerFundingAmountOctas uint64 = 1_000_000_000_000

	aptosWorkflowTimeout              = 4 * time.Minute
	aptosOnchainBenchmarkTimeout      = 2 * time.Minute
	aptosOnchainBenchmarkPollInterval = 3 * time.Second

	aptosTestFeedIDSuffix        = byte(1)
	aptosWriteBenchmarkValue     = uint64(123456789)
	aptosRoundtripBenchmarkValue = uint64(987654321)

	aptosScenarioOverrideEnv = "CRE_APTOS_SCENARIOS"
)

var aptosForwarderVersion = semver.MustParse("1.0.0")
var aptosWorkflowNameSeq uint64

// ExecuteAptosTest runs the Aptos CRE suite with the current CI scenario set by
// default. Individual scenarios still remain available for local/manual
// execution via CRE_APTOS_SCENARIOS.
func ExecuteAptosTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, resolveAptosScenarios(t))
}

type aptosScenario struct {
	name          string
	requiresWrite bool
	run           func(
		t *testing.T,
		tenv *configuration.TestEnvironment,
		aptosChain blockchains.Blockchain,
		userLogsCh <-chan *workflowevents.UserLogs,
		baseMessageCh <-chan *commonevents.BaseMessage,
	)
}

func aptosDefaultScenarios() []aptosScenario {
	return []aptosScenario{
		{name: "Aptos Write Read Roundtrip", requiresWrite: true, run: ExecuteAptosWriteReadRoundtripTest},
		{name: "Aptos Write Expected Failure", requiresWrite: true, run: ExecuteAptosWriteExpectedFailureTest},
	}
}

func resolveAptosScenarios(t *testing.T) []aptosScenario {
	t.Helper()

	raw := strings.TrimSpace(os.Getenv(aptosScenarioOverrideEnv))
	if raw == "" {
		return aptosDefaultScenarios()
	}
	if strings.EqualFold(raw, "ci") {
		t.Logf("running Aptos scenarios from %s=%q", aptosScenarioOverrideEnv, raw)
		return aptosDefaultScenarios()
	}

	available := map[string]aptosScenario{
		"read": {
			name:          "Aptos Read",
			requiresWrite: false,
			run:           ExecuteAptosReadTest,
		},
		"write": {
			name:          "Aptos Write",
			requiresWrite: true,
			run:           ExecuteAptosWriteTest,
		},
		"roundtrip": {
			name:          "Aptos Write Read Roundtrip",
			requiresWrite: true,
			run:           ExecuteAptosWriteReadRoundtripTest,
		},
		"write-expected-failure": {
			name:          "Aptos Write Expected Failure",
			requiresWrite: true,
			run:           ExecuteAptosWriteExpectedFailureTest,
		},
	}

	parts := strings.Split(raw, ",")
	scenarios := make([]aptosScenario, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		key := strings.ToLower(strings.TrimSpace(part))
		if key == "" {
			continue
		}
		scenario, ok := available[key]
		require.Truef(t, ok, "unknown Aptos scenario %q in %s", key, aptosScenarioOverrideEnv)
		if _, duplicate := seen[scenario.name]; duplicate {
			continue
		}
		seen[scenario.name] = struct{}{}
		scenarios = append(scenarios, scenario)
	}

	require.NotEmptyf(t, scenarios, "%s was set but did not resolve to any Aptos scenarios", aptosScenarioOverrideEnv)
	t.Logf("running Aptos scenarios from %s=%q", aptosScenarioOverrideEnv, raw)

	return scenarios
}

func executeAptosScenarios(t *testing.T, tenv *configuration.TestEnvironment, scenarios []aptosScenario) {
	aptosChain := mustAptosChainInEnv(t, tenv)
	lggr := framework.L

	writeDon := findAptosDonForChain(t, tenv, aptosChain.ChainID())
	assertAptosWorkerRuntimeKeysMatchMetadata(t, writeDon)
	if aptosScenariosRequireWriteSetup(scenarios) {
		ensureAptosWriteWorkersFunded(t, aptosChain, writeDon)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			if parallelEnabled && fanoutEnabled {
				t.Parallel()
			}

			scenarioEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, tenv.TestConfig)
			scenarioAptosChain := mustAptosChainInEnv(t, scenarioEnv)

			userLogsCh := make(chan *workflowevents.UserLogs, 1000)
			baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
			server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(lggr, userLogsCh, baseMessageCh))
			t.Cleanup(func() {
				server.Shutdown(t.Context())
				close(userLogsCh)
				close(baseMessageCh)
			})

			scenario.run(t, scenarioEnv, scenarioAptosChain, userLogsCh, baseMessageCh)
		})
	}
}

func aptosScenariosRequireWriteSetup(scenarios []aptosScenario) bool {
	for _, scenario := range scenarios {
		if scenario.requiresWrite {
			return true
		}
	}
	return false
}

func mustAptosChainInEnv(t *testing.T, tenv *configuration.TestEnvironment) blockchains.Blockchain {
	t.Helper()

	require.NotNil(t, tenv, "Aptos suite requires a test environment")
	require.NotNil(t, tenv.CreEnvironment, "Aptos suite requires a CRE environment")
	require.NotEmpty(t, tenv.CreEnvironment.Blockchains, "Aptos suite expects at least one blockchain in the environment")

	for _, bc := range tenv.CreEnvironment.Blockchains {
		if bc.IsFamily(blockchain.FamilyAptos) {
			return bc
		}
	}

	require.FailNow(t, "Aptos suite expects an Aptos chain in the environment (use config workflow-gateway-don-aptos.toml)")
	return nil
}

func assertAptosWorkerRuntimeKeysMatchMetadata(t *testing.T, writeDon *crelib.Don) {
	t.Helper()

	workers, err := writeDon.Workers()
	require.NoError(t, err, "failed to list Aptos write DON workers")
	require.NotEmpty(t, workers, "Aptos write DON workers list is empty")

	for _, worker := range workers {
		require.NotNil(t, worker.Keys, "worker %q is missing metadata keys", worker.Name)
		require.NotNil(t, worker.Keys.Aptos, "worker %q is missing metadata Aptos key", worker.Name)

		expectedAccount, err := crecrypto.NormalizeAptosAccount(worker.Keys.Aptos.Account)
		require.NoError(t, err, "worker %q has invalid metadata Aptos account", worker.Name)
		expectedPublicKey := normalizeHexValue(worker.Keys.Aptos.PublicKey)
		require.NotEmpty(t, expectedPublicKey, "worker %q is missing metadata Aptos public key", worker.Name)

		var runtimeKeys struct {
			Data []struct {
				Attributes struct {
					Account   string `json:"account"`
					PublicKey string `json:"publicKey"`
				} `json:"attributes"`
			} `json:"data"`
		}
		resp, err := worker.Clients.RestClient.APIClient.R().
			SetResult(&runtimeKeys).
			Get("/v2/keys/aptos")
		require.NoError(t, err, "failed to read runtime Aptos keys for worker %q", worker.Name)
		require.Equal(t, http.StatusOK, resp.StatusCode(), "worker %q Aptos keys endpoint returned unexpected status", worker.Name)
		require.Len(t, runtimeKeys.Data, 1, "worker %q must expose exactly one Aptos runtime key", worker.Name)

		runtimeKey := runtimeKeys.Data[0].Attributes
		actualAccount, err := crecrypto.NormalizeAptosAccount(runtimeKey.Account)
		require.NoError(t, err, "worker %q exposed invalid runtime Aptos account", worker.Name)
		require.Equal(t, expectedAccount, actualAccount, "worker %q runtime Aptos account does not match metadata-generated account", worker.Name)
		require.Equal(t, expectedPublicKey, normalizeHexValue(runtimeKey.PublicKey), "worker %q runtime Aptos public key does not match metadata-generated key", worker.Name)
	}
}

// ExecuteAptosReadTest deploys a workflow that reads 0x1::coin::name() on Aptos local devnet
// in a consensus read step and asserts the expected value.
func ExecuteAptosReadTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	lggr := framework.L

	ensureAptosLedgerVersionPositive(t, aptosChain)

	// Fixed name so re-runs against the same DON overwrite the same workflow instead of accumulating multiple (e.g. aptos-read-workflow-4838 and aptos-read-workflow-5736).
	const workflowName = "aptos-read-workflow"
	workflowConfig := t_helpers.AptosReadWorkflowConfig{
		ChainSelector:    aptosChain.ChainSelector(),
		WorkflowName:     workflowName,
		ExpectedCoinName: "Aptos Coin", // 0x1::coin::name<0x1::aptos_coin::AptosCoin>() on local devnet
	}

	const workflowFileLocation = "./aptos/aptosread/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	expectedLog := "Aptos read consensus succeeded"
	t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, aptosWorkflowTimeout, t_helpers.WithUserLogWorkflowID(workflowID))
	lggr.Info().Str("expected_log", expectedLog).Msg("Aptos read capability test passed")
}

// The local Aptos devnet can legitimately remain at ledger version 0 until the
// first Aptos transaction lands, which makes the read capability reject the
// consensus height before the workflow has a chance to execute.
func ensureAptosLedgerVersionPositive(t *testing.T, aptosChain blockchains.Blockchain) {
	t.Helper()

	aptosBC, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")

	client, err := aptosBC.NodeClient()
	require.NoError(t, err, "failed to create Aptos node client")

	info, err := client.Info()
	require.NoError(t, err, "failed to fetch Aptos node info")
	if info.LedgerVersion() > 0 {
		return
	}

	deployer, err := aptosBC.LocalDeployerAccount()
	require.NoError(t, err, "failed to create Aptos local deployer account")
	deployerAddress := deployer.AccountAddress()

	require.NoError(t,
		aptosBC.Fund(t.Context(), deployerAddress.StringLong(), 1),
		"failed to bump Aptos ledger version via faucet funding",
	)

	require.Eventuallyf(t, func() bool {
		nodeInfo, nodeErr := client.Info()
		return nodeErr == nil && nodeInfo.LedgerVersion() > 0
	}, 30*time.Second, time.Second, "expected Aptos ledger version to become positive after funding")
}

func ExecuteAptosWriteTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	lggr := framework.L
	scenario := prepareAptosWriteScenario(t, tenv, aptosChain)

	workflowName := uniqueAptosWorkflowName("aptos-write-workflow")
	workflowConfig := aptoswrite_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        scenario.receiverHex,
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		// Keep within the current local Aptos transaction max-gas bound.
		MaxGasAmount: aptosLocalMaxGasAmount,
		GasUnitPrice: aptosLocalGasUnitPrice,
	}

	const workflowFileLocation = "./aptos/aptoswrite/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	txHash := waitForAptosWriteSuccessLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, workflowID, aptosWorkflowTimeout)
	assertAptosReceiverUpdatedOnChain(t, aptosChain, scenario.receiverHex, scenario.expectedBenchmarkValue)
	assertAptosWriteTxOnChain(t, aptosChain, txHash, scenario.receiverHex)
	lggr.Info().
		Str("tx_hash", txHash).
		Str("receiver", scenario.receiverHex).
		Msg("Aptos write capability test passed with onchain verification")
}

func ExecuteAptosWriteReadRoundtripTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	lggr := framework.L
	scenario := prepareAptosRoundtripScenario(t, tenv, aptosChain)

	workflowName := uniqueAptosWorkflowName("aptos-write-read-roundtrip-workflow")
	roundtripCfg := aptoswriteroundtrip_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        scenario.receiverHex,
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		MaxGasAmount:       aptosLocalMaxGasAmount,
		GasUnitPrice:       aptosLocalGasUnitPrice,
		FeedIDHex:          scenario.feedIDHex,
		ExpectedBenchmark:  scenario.expectedBenchmarkValue,
	}

	workflowID := t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &roundtripCfg, "./aptos/aptoswriteroundtrip/main.go")
	t_helpers.WatchWorkflowLogs(
		t,
		lggr,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		"Aptos write/read consensus succeeded",
		aptosWorkflowTimeout,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)
	lggr.Info().
		Str("receiver", scenario.receiverHex).
		Uint64("expected_benchmark", scenario.expectedBenchmarkValue).
		Str("feed_id", scenario.feedIDHex).
		Msg("Aptos write/read roundtrip capability test passed")
}

func ExecuteAptosWriteExpectedFailureTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	lggr := framework.L
	scenario := prepareAptosWriteFailureScenario(t, tenv, aptosChain)

	workflowName := uniqueAptosWorkflowName("aptos-write-expected-failure-workflow")
	workflowConfig := aptoswrite_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        "0x0", // Intentionally invalid write receiver to force onchain failure path.
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		MaxGasAmount:       aptosLocalMaxGasAmount,
		GasUnitPrice:       aptosLocalGasUnitPrice,
		ExpectFailure:      true,
	}

	const workflowFileLocation = "./aptos/aptoswrite/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	txHash := waitForAptosWriteExpectedFailureLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, workflowID, aptosWorkflowTimeout)
	assertAptosWriteFailureTxOnChain(t, aptosChain, txHash)

	lggr.Info().
		Str("tx_hash", txHash).
		Msg("Aptos expected write-failure workflow test passed")
}

type aptosWriteScenario struct {
	chainSelector          uint64
	receiverHex            string
	reportPayloadHex       string
	feedIDHex              string
	expectedBenchmarkValue uint64
	requiredSignatures     int
	writeDon               *crelib.Don
}

func prepareAptosWriteScenario(t *testing.T, tenv *configuration.TestEnvironment, aptosChain blockchains.Blockchain) aptosWriteScenario {
	return prepareAptosWriteScenarioWithBenchmark(
		t,
		tenv,
		aptosChain,
		aptosTestFeedID(),
		aptosWriteBenchmarkValue,
	)
}

func prepareAptosRoundtripScenario(t *testing.T, tenv *configuration.TestEnvironment, aptosChain blockchains.Blockchain) aptosWriteScenario {
	return prepareAptosWriteScenarioWithBenchmark(
		t,
		tenv,
		aptosChain,
		aptosTestFeedID(),
		aptosRoundtripBenchmarkValue,
	)
}

func prepareAptosWriteFailureScenario(t *testing.T, tenv *configuration.TestEnvironment, aptosChain blockchains.Blockchain) aptosWriteScenario {
	t.Helper()

	writeDon := findAptosDonForChain(t, tenv, aptosChain.ChainID())
	workers, workerErr := writeDon.Workers()
	require.NoError(t, workerErr, "failed to list Aptos write DON workers")
	f := (len(workers) - 1) / 3
	require.GreaterOrEqual(t, f, 1, "Aptos write DON requires f>=1")

	return aptosWriteScenario{
		chainSelector:      aptosChain.ChainSelector(),
		reportPayloadHex:   hex.EncodeToString(buildAptosDataFeedsBenchmarkPayloadFor(aptosTestFeedID(), aptosWriteBenchmarkValue)),
		requiredSignatures: f + 1,
		writeDon:           writeDon,
	}
}

func prepareAptosWriteScenarioWithBenchmark(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	feedID []byte,
	expectedBenchmark uint64,
) aptosWriteScenario {
	t.Helper()

	forwarderHex := aptosForwarderAddress(tenv, aptosChain.ChainSelector())
	require.NotEmpty(t, forwarderHex, "Aptos write test requires forwarder address for chainSelector=%d", aptosChain.ChainSelector())
	require.False(t, isZeroAptosAddress(forwarderHex), "Aptos write test requires non-zero forwarder address for chainSelector=%d", aptosChain.ChainSelector())

	writeDon := findAptosDonForChain(t, tenv, aptosChain.ChainID())
	workers, workerErr := writeDon.Workers()
	require.NoError(t, workerErr, "failed to list Aptos write DON workers")
	f := (len(workers) - 1) / 3
	require.GreaterOrEqual(t, f, 1, "Aptos write DON requires f>=1")

	return aptosWriteScenario{
		chainSelector:          aptosChain.ChainSelector(),
		receiverHex:            deployAptosDataFeedsReceiverForWrite(t, tenv, aptosChain, forwarderHex, feedID),
		reportPayloadHex:       hex.EncodeToString(buildAptosDataFeedsBenchmarkPayloadFor(feedID, expectedBenchmark)),
		feedIDHex:              hex.EncodeToString(feedID),
		expectedBenchmarkValue: expectedBenchmark,
		requiredSignatures:     f + 1,
		writeDon:               writeDon,
	}
}

func uniqueAptosWorkflowName(base string) string {
	return fmt.Sprintf("%s-%d-%d", base, time.Now().UnixNano(), atomic.AddUint64(&aptosWorkflowNameSeq, 1))
}

func findAptosDonForChain(t *testing.T, tenv *configuration.TestEnvironment, chainID uint64) *crelib.Don {
	t.Helper()
	require.NotNil(t, tenv.Dons, "test environment DON metadata is required")

	for _, don := range tenv.Dons.List() {
		if !don.HasFlag("aptos") {
			continue
		}
		chainIDs, err := don.GetEnabledChainIDsForCapability("aptos")
		require.NoError(t, err, "failed to read enabled chain ids for DON %q", don.Name)
		for _, id := range chainIDs {
			if id == chainID {
				return don
			}
		}
	}

	require.FailNowf(t, "missing Aptos DON", "could not find aptos DON for chainID=%d", chainID)
	return nil
}

func isZeroAptosAddress(addr string) bool {
	trimmed := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(addr)), "0x")
	if trimmed == "" {
		return true
	}
	for _, ch := range trimmed {
		if ch != '0' {
			return false
		}
	}
	return true
}

func aptosForwarderAddress(tenv *configuration.TestEnvironment, chainSelector uint64) string {
	return crecontracts.MustGetAddressFromDataStore(
		tenv.CreEnvironment.CldfEnvironment.DataStore,
		chainSelector,
		"AptosForwarder",
		aptosForwarderVersion,
		"",
	)
}

var aptosTxHashInLogRe = regexp.MustCompile(`txHash=([^\s"]+)`)

func waitForAptosWriteSuccessLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	workflowID string,
	timeout time.Duration,
) string {
	t.Helper()
	return waitForAptosLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, workflowID, "Aptos write capability succeeded", timeout)
}

func waitForAptosWriteExpectedFailureLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	workflowID string,
	timeout time.Duration,
) string {
	t.Helper()
	return waitForAptosLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, workflowID, "Aptos write failure observed as expected", timeout)
}

func waitForAptosLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	workflowID string,
	expectedLog string,
	timeout time.Duration,
) string {
	t.Helper()

	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, fmt.Errorf("failed to find Aptos workflow log with non-empty tx hash: %s", expectedLog))
	defer cancelFn()

	cancelCtx, cancelCauseFn := context.WithCancelCause(ctx)
	defer cancelCauseFn(nil)

	go func() {
		if workflowID != "" {
			t_helpers.FailOnBaseMessage(cancelCtx, cancelCauseFn, t, lggr, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, t_helpers.WithBaseMessageWorkflowID(workflowID))
			return
		}
		t_helpers.FailOnBaseMessage(cancelCtx, cancelCauseFn, t, lggr, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)
	}()

	mismatchCount := 0
	for {
		select {
		case <-cancelCtx.Done():
			require.NoError(t, context.Cause(cancelCtx), "failed to observe Aptos log with non-empty tx hash: %s", expectedLog)
			return ""
		case logs := <-userLogsCh:
			if workflowID != "" && !aptosUserLogsHaveWorkflowID(logs, workflowID) {
				continue
			}
			for _, line := range logs.LogLines {
				if !strings.Contains(line.Message, expectedLog) {
					mismatchCount++
					if mismatchCount%20 == 0 {
						lggr.Warn().
							Str("expected_log", expectedLog).
							Str("found_message", strings.TrimSpace(line.Message)).
							Int("mismatch_count", mismatchCount).
							Msg("[soft assertion] Received UserLogs messages, but none match expected log yet")
					}
					continue
				}

				matches := aptosTxHashInLogRe.FindStringSubmatch(line.Message)
				if len(matches) == 2 {
					txHash := normalizeTxHash(matches[1])
					if txHash != "" {
						return txHash
					}
				}

				lggr.Warn().
					Str("message", strings.TrimSpace(line.Message)).
					Str("expected_log", expectedLog).
					Msg("[soft assertion] Matched Aptos log without non-empty tx hash; waiting for another match")
			}
		}
	}
}

func aptosUserLogsHaveWorkflowID(logs *workflowevents.UserLogs, workflowID string) bool {
	if logs == nil || logs.M == nil || workflowID == "" {
		return false
	}
	return normalizeHexValue(logs.M.WorkflowID) == normalizeHexValue(workflowID)
}

func assertAptosWriteFailureTxOnChain(t *testing.T, aptosChain blockchains.Blockchain, txHash string) {
	t.Helper()

	bc, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")

	nodeURL := bc.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for onchain verification")
	nodeURL, err := blockchains_aptos.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := blockchains_aptos.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")

	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	tx, err := client.WaitForTransaction(txHash)
	require.NoError(t, err, "failed waiting for Aptos tx by hash")
	require.False(t, tx.Success, "Aptos tx must fail in expected-failure workflow; vm_status=%s", tx.VmStatus)
}

func assertAptosWriteTxOnChain(t *testing.T, aptosChain blockchains.Blockchain, txHash string, expectedReceiver string) {
	t.Helper()

	bc, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")

	nodeURL := bc.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for onchain verification")
	nodeURL, err := blockchains_aptos.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := blockchains_aptos.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")

	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	tx, err := client.WaitForTransaction(txHash)
	require.NoError(t, err, "failed waiting for Aptos tx by hash")
	require.True(t, tx.Success, "Aptos tx must be successful; vm_status=%s", tx.VmStatus)

	expectedReceiverNorm := normalizeTxHashLikeHex(expectedReceiver)
	found := false
	for _, evt := range tx.Events {
		if !strings.HasSuffix(evt.Type, "::forwarder::ReportProcessed") {
			continue
		}
		receiverVal, ok := evt.Data["receiver"].(string)
		require.True(t, ok, "ReportProcessed event receiver field must be a string")
		if normalizeTxHashLikeHex(receiverVal) != expectedReceiverNorm {
			continue
		}
		_, hasExecutionID := evt.Data["workflow_execution_id"]
		_, hasReportID := evt.Data["report_id"]
		require.True(t, hasExecutionID, "ReportProcessed must include workflow_execution_id")
		require.True(t, hasReportID, "ReportProcessed must include report_id")
		found = true
		break
	}
	require.True(t, found, "expected ReportProcessed event for receiver %s in tx %s", expectedReceiverNorm, txHash)
}

func assertAptosReceiverUpdatedOnChain(
	t *testing.T,
	aptosChain blockchains.Blockchain,
	receiverHex string,
	expectedBenchmark uint64,
) {
	t.Helper()

	aptosBC, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")
	nodeURL := aptosBC.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for onchain verification")
	nodeURL, err := blockchains_aptos.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := blockchains_aptos.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")
	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	var receiverAddr aptoslib.AccountAddress
	err = receiverAddr.ParseStringRelaxed(receiverHex)
	require.NoError(t, err, "failed to parse Aptos receiver address")

	dataFeeds := aptosdatafeeds.Bind(receiverAddr, client)
	feedID := aptosTestFeedID()
	feedIDHex := hex.EncodeToString(feedID)

	require.Eventually(t, func() bool {
		feeds, bErr := dataFeeds.Registry().GetFeeds(&aptosbind.CallOpts{})
		if bErr != nil || len(feeds) == 0 {
			return false
		}
		for _, feed := range feeds {
			if hex.EncodeToString(feed.FeedId) != feedIDHex {
				continue
			}
			if feed.Feed.Benchmark == nil {
				return false
			}
			return feed.Feed.Benchmark.Uint64() == expectedBenchmark
		}
		return false
	}, aptosOnchainBenchmarkTimeout, aptosOnchainBenchmarkPollInterval, "expected benchmark value %d not observed onchain for receiver %s", expectedBenchmark, receiverHex)
}

func normalizeTxHash(input string) string {
	s := strings.TrimSpace(strings.ToLower(input))
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return "0x" + s
}

func normalizeTxHashLikeHex(input string) string {
	s := strings.TrimSpace(strings.ToLower(input))
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "0x0"
	}
	return "0x" + s
}

func normalizeHexValue(input string) string {
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(input)), "0x")
}

func deployAptosDataFeedsReceiverForWrite(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	primaryForwarderHex string,
	feedID []byte,
) string {
	t.Helper()

	aptosBC, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")
	nodeURL := aptosBC.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for receiver deployment")
	nodeURL, err := blockchains_aptos.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for receiver deployment")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := blockchains_aptos.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")
	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	deployer, err := aptosDeployerAccount()
	require.NoError(t, err, "failed to create Aptos deployer account")
	deployerAddress := deployer.AccountAddress()
	require.NoError(t, aptosBC.Fund(t.Context(), deployerAddress.StringLong(), aptosWorkerFundingAmountOctas), "failed to fund Aptos deployer account")

	var primaryForwarderAddr aptoslib.AccountAddress
	err = primaryForwarderAddr.ParseStringRelaxed(primaryForwarderHex)
	require.NoError(t, err, "failed to parse primary forwarder address")

	owner := deployerAddress
	secondaryAddress, secondaryTx, _, err := aptosplatformsecondary.DeployToObject(deployer, client, owner)
	require.NoError(t, err, "failed to deploy Aptos secondary platform package")
	require.NoError(t, blockchains_aptos.WaitForTransactionSuccess(client, secondaryTx.Hash, "platform_secondary deployment"))

	dataFeedsAddress, dataFeedsTx, dataFeeds, err := aptosdatafeeds.DeployToObject(
		deployer,
		client,
		owner,
		primaryForwarderAddr,
		owner,
		secondaryAddress,
	)
	require.NoError(t, err, "failed to deploy Aptos data feeds receiver package")
	require.NoError(t, blockchains_aptos.WaitForTransactionSuccess(client, dataFeedsTx.Hash, "data_feeds deployment"))

	workflowOwner := workflowRegistryOwnerBytes(t, tenv)
	tx, err := dataFeeds.Registry().SetWorkflowConfig(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{workflowOwner},
		[][]byte{},
	)
	require.NoError(t, err, "failed to set data feeds workflow config")
	require.NoError(t, blockchains_aptos.WaitForTransactionSuccess(client, tx.Hash, "data_feeds set_workflow_config"))

	// Configure the feed that the write workflow will update.
	// Without this, registry::perform_update emits WriteSkippedFeedNotSet and benchmark remains unchanged.
	tx, err = dataFeeds.Registry().SetFeeds(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{feedID},
		[]string{"CRE-BENCHMARK"},
		[]byte{0x99},
	)
	require.NoError(t, err, "failed to set data feeds feed config")
	require.NoError(t, blockchains_aptos.WaitForTransactionSuccess(client, tx.Hash, "data_feeds set_feeds"))

	return dataFeedsAddress.StringLong()
}

func aptosDeployerAccount() (*aptoslib.Account, error) {
	const defaultAptosDeployerKey = "d477c65f88ed9e6d4ec6e2014755c3cfa3e0c44e521d0111a02868c5f04c41d4"
	keyHex := strings.TrimSpace(os.Getenv("CRE_APTOS_DEPLOYER_PRIVATE_KEY"))
	if keyHex == "" {
		keyHex = defaultAptosDeployerKey
	}
	if keyHex == "" {
		return nil, errors.New("empty Aptos deployer key")
	}
	keyHex = strings.TrimPrefix(keyHex, "0x")
	var privateKey aptoscrypto.Ed25519PrivateKey
	if err := privateKey.FromHex(keyHex); err != nil {
		return nil, fmt.Errorf("parse Aptos deployer private key: %w", err)
	}
	return aptoslib.NewAccountFromSigner(&privateKey)
}

func ensureAptosWriteWorkersFunded(t *testing.T, aptosChain blockchains.Blockchain, writeDon *crelib.Don) {
	t.Helper()

	aptosBC, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")
	workers, workerErr := writeDon.Workers()
	require.NoError(t, workerErr, "failed to list Aptos write DON workers for funding")
	require.NotEmpty(t, workers, "Aptos write DON workers list is empty")

	for _, worker := range workers {
		require.NotNil(t, worker.Keys, "worker %q is missing metadata keys", worker.Name)
		require.NotNil(t, worker.Keys.Aptos, "worker %q is missing metadata Aptos key", worker.Name)

		var account aptoslib.AccountAddress
		parseErr := account.ParseStringRelaxed(worker.Keys.Aptos.Account)
		require.NoError(t, parseErr, "failed to parse Aptos worker account for worker %q", worker.Name)

		require.NoError(t, aptosBC.Fund(t.Context(), account.StringLong(), aptosWorkerFundingAmountOctas), "failed to fund Aptos worker account %s for worker %q", account.StringLong(), worker.Name)
	}
}

func workflowRegistryOwnerBytes(t *testing.T, tenv *configuration.TestEnvironment) []byte {
	t.Helper()
	registryChain, ok := tenv.CreEnvironment.Blockchains[0].(*blockchains_evm.Blockchain)
	require.True(t, ok, "registry chain must be EVM")
	rootOwner := registryChain.SethClient.MustGetRootKeyAddress()
	return common.HexToAddress(rootOwner.Hex()).Bytes()
}

func buildAptosDataFeedsBenchmarkPayloadFor(feedID []byte, benchmark uint64) []byte {
	// ABI-like benchmark payload expected by data_feeds::registry::parse_raw_report
	// [offset=32][count=1][feed_id(32)][report(64)]
	const (
		offsetToArray = uint64(32)
		reportCount   = uint64(1)
		timestamp     = uint64(1700000000)
	)

	report := make([]byte, 64)
	writeU256BE(report[0:32], timestamp)
	writeU256BE(report[32:64], benchmark)

	out := make([]byte, 0, 160)
	out = appendU256BE(out, offsetToArray)
	out = appendU256BE(out, reportCount)
	out = append(out, feedID...)
	out = append(out, report...)
	return out
}

func aptosTestFeedID() []byte {
	feedID := make([]byte, 32)
	feedID[len(feedID)-1] = aptosTestFeedIDSuffix
	return feedID
}

func appendU256BE(dst []byte, v uint64) []byte {
	buf := make([]byte, 32)
	binary.BigEndian.PutUint64(buf[24:], v)
	return append(dst, buf...)
}

func writeU256BE(dst []byte, v uint64) {
	binary.BigEndian.PutUint64(dst[24:], v)
}
