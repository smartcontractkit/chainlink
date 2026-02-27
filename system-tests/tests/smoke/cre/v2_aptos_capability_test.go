package cre

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

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
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	blockchains_aptos "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	blockchains_evm "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	aptoswrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite/config"
	aptoswriteroundtrip_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswriteroundtrip/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// ExecuteAptosTest runs the Aptos CRE suite (read consensus test and any future scenarios).
func ExecuteAptosTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, true, true, true, true)
}

func ExecuteAptosReadOnlyTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, true, false, false, false)
}

func ExecuteAptosWriteOnlyTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, false, true, false, false)
}

func ExecuteAptosWriteReadRoundtripOnlyTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, false, false, true, false)
}

func ExecuteAptosWriteExpectedFailureOnlyTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, false, false, false, true)
}

func executeAptosScenarios(t *testing.T, tenv *configuration.TestEnvironment, runRead bool, runWrite bool, runRoundtrip bool, runWriteExpectedFailure bool) {
	creEnv := tenv.CreEnvironment
	require.NotEmpty(t, creEnv.Blockchains, "Aptos suite expects at least one blockchain in the environment")

	var aptosChain blockchains.Blockchain
	for _, bc := range creEnv.Blockchains {
		if bc.IsFamily(blockchain.FamilyAptos) {
			aptosChain = bc
			break
		}
	}
	require.NotNil(t, aptosChain, "Aptos suite expects an Aptos chain in the environment (use config workflow-gateway-don-aptos.toml)")

	lggr := framework.L
	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetLoggingPublishFn(lggr, userLogsCh, baseMessageCh, "./logs/aptos_capability_workflow_test.log"))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	if runRead {
		t.Run("Aptos Read", func(t *testing.T) {
			ExecuteAptosReadTest(t, tenv, aptosChain, userLogsCh, baseMessageCh)
		})
	}
	if runWrite {
		t.Run("Aptos Write", func(t *testing.T) {
			ExecuteAptosWriteTest(t, tenv, aptosChain, userLogsCh, baseMessageCh)
		})
	}
	if runRoundtrip {
		t.Run("Aptos Write Read Roundtrip", func(t *testing.T) {
			ExecuteAptosWriteReadRoundtripTest(t, tenv, aptosChain, userLogsCh, baseMessageCh)
		})
	}
	if runWriteExpectedFailure {
		t.Run("Aptos Write Expected Failure", func(t *testing.T) {
			ExecuteAptosWriteExpectedFailureTest(t, tenv, aptosChain, userLogsCh, baseMessageCh)
		})
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

	// Fixed name so re-runs against the same DON overwrite the same workflow instead of accumulating multiple (e.g. aptos-read-workflow-4838 and aptos-read-workflow-5736).
	const workflowName = "aptos-read-workflow"
	workflowConfig := t_helpers.AptosReadWorkflowConfig{
		ChainSelector:    aptosChain.ChainSelector(),
		WorkflowName:     workflowName,
		ExpectedCoinName: "Aptos Coin", // 0x1::coin::name<0x1::aptos_coin::AptosCoin>() on local devnet
	}

	const workflowFileLocation = "./aptos/aptosread/main.go"
	t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	expectedLog := "Aptos read consensus succeeded"
	t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, 4*time.Minute)
	lggr.Info().Str("expected_log", expectedLog).Msg("Aptos read capability test passed")
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

	const workflowName = "aptos-write-workflow"
	workflowConfig := aptoswrite_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        scenario.receiverHex,
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		// Keep within local Aptos devnet transaction max-gas bound.
		MaxGasAmount: 1_000_000,
		GasUnitPrice: 100,
	}

	const workflowFileLocation = "./aptos/aptoswrite/main.go"
	ensureAptosWriteWorkersFunded(t, aptosChain, scenario.writeDon)
	t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	txHash := waitForAptosWriteSuccessLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, 4*time.Minute)
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

	const workflowName = "aptos-write-read-roundtrip-workflow"
	roundtripCfg := aptoswriteroundtrip_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        scenario.receiverHex,
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		MaxGasAmount:       1_000_000,
		GasUnitPrice:       100,
		FeedIDHex:          scenario.feedIDHex,
		ExpectedBenchmark:  scenario.expectedBenchmarkValue,
	}

	ensureAptosWriteWorkersFunded(t, aptosChain, scenario.writeDon)
	t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &roundtripCfg, "./aptos/aptoswriteroundtrip/main.go")
	t_helpers.WatchWorkflowLogs(
		t,
		lggr,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		"Aptos write/read consensus succeeded",
		4*time.Minute,
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
	scenario := prepareAptosWriteScenario(t, tenv, aptosChain)

	const workflowName = "aptos-write-expected-failure-workflow"
	workflowConfig := aptoswrite_config.Config{
		ChainSelector:      scenario.chainSelector,
		WorkflowName:       workflowName,
		ReceiverHex:        "0x0", // Intentionally invalid write receiver to force onchain failure path.
		RequiredSignatures: scenario.requiredSignatures,
		ReportPayloadHex:   scenario.reportPayloadHex,
		MaxGasAmount:       1_000_000,
		GasUnitPrice:       100,
		ExpectFailure:      true,
	}

	const workflowFileLocation = "./aptos/aptoswrite/main.go"
	ensureAptosWriteWorkersFunded(t, aptosChain, scenario.writeDon)
	t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	txHash := waitForAptosWriteExpectedFailureLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, 4*time.Minute)
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
	return prepareAptosWriteScenarioWithBenchmark(t, tenv, aptosChain, aptosBenchmarkFeedID(), 123456789)
}

func prepareAptosRoundtripScenario(t *testing.T, tenv *configuration.TestEnvironment, aptosChain blockchains.Blockchain) aptosWriteScenario {
	return prepareAptosWriteScenarioWithBenchmark(t, tenv, aptosChain, aptosRoundtripFeedID(), 987654321)
}

func prepareAptosWriteScenarioWithBenchmark(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	feedID []byte,
	expectedBenchmark uint64,
) aptosWriteScenario {
	t.Helper()

	forwarderHex := ""
	if tenv.CreEnvironment.AptosForwarderAddresses != nil {
		forwarderHex = tenv.CreEnvironment.AptosForwarderAddresses[aptosChain.ChainSelector()]
	}
	require.NotEmpty(t, forwarderHex, "Aptos write test requires forwarder address for chainSelector=%d", aptosChain.ChainSelector())
	require.False(t, isZeroAptosAddress(forwarderHex), "Aptos write test requires non-zero forwarder address for chainSelector=%d", aptosChain.ChainSelector())

	writeDon := findWriteAptosDonForChain(t, tenv, aptosChain.ChainID())
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

func findWriteAptosDonForChain(t *testing.T, tenv *configuration.TestEnvironment, chainID uint64) *crelib.Don {
	t.Helper()
	require.NotNil(t, tenv.Dons, "test environment DON metadata is required")

	for _, don := range tenv.Dons.List() {
		if !don.HasFlag("write-aptos") {
			continue
		}
		chainIDs, err := don.GetEnabledChainIDsForCapability("write-aptos")
		require.NoError(t, err, "failed to read enabled chain ids for DON %q", don.Name)
		for _, id := range chainIDs {
			if id == chainID {
				return don
			}
		}
	}

	require.FailNowf(t, "missing Aptos write DON", "could not find write-aptos DON for chainID=%d", chainID)
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

var aptosTxHashInLogRe = regexp.MustCompile(`txHash=([^\s"]+)`)

func waitForAptosWriteSuccessLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	timeout time.Duration,
) string {
	t.Helper()
	return waitForAptosLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, "Aptos write capability succeeded", timeout)
}

func waitForAptosWriteExpectedFailureLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	timeout time.Duration,
) string {
	t.Helper()
	return waitForAptosLogAndTxHash(t, lggr, userLogsCh, baseMessageCh, "Aptos write failure observed as expected", timeout)
}

func waitForAptosLogAndTxHash(
	t *testing.T,
	lggr zerolog.Logger,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	expectedLog string,
	timeout time.Duration,
) string {
	t.Helper()

	ctx, cancelFn := context.WithTimeoutCause(t.Context(), timeout, fmt.Errorf("failed to find Aptos workflow log with non-empty tx hash: %s", expectedLog))
	defer cancelFn()

	cancelCtx, cancelCauseFn := context.WithCancelCause(ctx)
	defer cancelCauseFn(nil)

	go func() {
		t_helpers.FailOnBaseMessage(cancelCtx, cancelCauseFn, t, lggr, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog)
	}()

	mismatchCount := 0
	for {
		select {
		case <-cancelCtx.Done():
			require.NoError(t, context.Cause(cancelCtx), "failed to observe Aptos log with non-empty tx hash: %s", expectedLog)
			return ""
		case logs := <-userLogsCh:
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

func assertAptosWriteFailureTxOnChain(t *testing.T, aptosChain blockchains.Blockchain, txHash string) {
	t.Helper()

	bc, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")

	nodeURL := bc.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for onchain verification")
	nodeURL, err := normalizeAptosNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")

	client, err := aptoslib.NewNodeClient(nodeURL, uint8(chainID))
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
	nodeURL, err := normalizeAptosNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")

	client, err := aptoslib.NewNodeClient(nodeURL, uint8(chainID))
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
	nodeURL, err := normalizeAptosNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	client, err := aptoslib.NewNodeClient(nodeURL, uint8(chainID))
	require.NoError(t, err, "failed to create Aptos client")

	var receiverAddr aptoslib.AccountAddress
	err = receiverAddr.ParseStringRelaxed(receiverHex)
	require.NoError(t, err, "failed to parse Aptos receiver address")

	dataFeeds := aptosdatafeeds.Bind(receiverAddr, client)
	feedID := aptosBenchmarkFeedID()
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
	}, 2*time.Minute, 3*time.Second, "expected benchmark value %d not observed onchain for receiver %s", expectedBenchmark, receiverHex)

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
	nodeURL, err := normalizeAptosNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for receiver deployment")
	containerName := aptosBC.CtfOutput().ContainerName

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	client, err := aptoslib.NewNodeClient(nodeURL, uint8(chainID))
	require.NoError(t, err, "failed to create Aptos client")

	deployer, err := aptosDeployerAccount()
	require.NoError(t, err, "failed to create Aptos deployer account")

	fundAptosAccountBestEffort(t, client, nodeURL, containerName, deployer.AccountAddress())
	require.NoError(t, waitForAptosAccountVisible(client, deployer.AccountAddress(), 45*time.Second), "Aptos deployer account must be visible before deploy")

	var primaryForwarderAddr aptoslib.AccountAddress
	err = primaryForwarderAddr.ParseStringRelaxed(primaryForwarderHex)
	require.NoError(t, err, "failed to parse primary forwarder address")

	owner := deployer.AccountAddress()
	secondaryAddress, secondaryTx, _, err := aptosplatformsecondary.DeployToObject(deployer, client, owner)
	require.NoError(t, err, "failed to deploy Aptos secondary platform package")
	waitForAptosTransactionSuccess(t, client, secondaryTx.Hash, "platform_secondary deployment")

	dataFeedsAddress, dataFeedsTx, dataFeeds, err := aptosdatafeeds.DeployToObject(
		deployer,
		client,
		owner,
		primaryForwarderAddr,
		owner,
		secondaryAddress,
	)
	require.NoError(t, err, "failed to deploy Aptos data feeds receiver package")
	waitForAptosTransactionSuccess(t, client, dataFeedsTx.Hash, "data_feeds deployment")

	workflowOwner := workflowRegistryOwnerBytes(t, tenv)
	tx, err := dataFeeds.Registry().SetWorkflowConfig(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{workflowOwner},
		[][]byte{},
	)
	require.NoError(t, err, "failed to set data feeds workflow config")
	waitForAptosTransactionSuccess(t, client, tx.Hash, "data_feeds set_workflow_config")

	// Configure the feed that the write workflow will update.
	// Without this, registry::perform_update emits WriteSkippedFeedNotSet and benchmark remains unchanged.
	tx, err = dataFeeds.Registry().SetFeeds(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{feedID},
		[]string{"CRE-BENCHMARK"},
		[]byte{0x99},
	)
	require.NoError(t, err, "failed to set data feeds feed config")
	waitForAptosTransactionSuccess(t, client, tx.Hash, "data_feeds set_feeds")

	return dataFeedsAddress.StringLong()
}

func aptosDeployerAccount() (*aptoslib.Account, error) {
	const defaultAptosDeployerKey = "d477c65f88ed9e6d4ec6e2014755c3cfa3e0c44e521d0111a02868c5f04c41d4"
	keyHex := strings.TrimSpace(os.Getenv("CRE_APTOS_DEPLOYER_PRIVATE_KEY"))
	if keyHex == "" {
		keyHex = defaultAptosDeployerKey
	}
	if keyHex == "" {
		return nil, fmt.Errorf("empty Aptos deployer key")
	}
	keyHex = strings.TrimPrefix(keyHex, "0x")
	var privateKey aptoscrypto.Ed25519PrivateKey
	if err := privateKey.FromHex(keyHex); err != nil {
		return nil, fmt.Errorf("parse Aptos deployer private key: %w", err)
	}
	return aptoslib.NewAccountFromSigner(&privateKey)
}

func fundAptosAccountBestEffort(
	t *testing.T,
	client *aptoslib.NodeClient,
	nodeURL string,
	containerName string,
	account aptoslib.AccountAddress,
) {
	t.Helper()
	// Fast path: account already exists.
	if _, err := client.Account(account); err == nil {
		return
	}

	faucetURL, err := aptosFaucetURLFromNodeURL(nodeURL)
	if err == nil {
		if faucetClient, cErr := aptoslib.NewFaucetClient(client, faucetURL); cErr != nil {
			framework.L.Warn().
				Err(cErr).
				Str("faucet_url", faucetURL).
				Str("account", account.StringLong()).
				Msg("Aptos faucet client init failed; trying container fallback")
		} else {
			fundErr := faucetClient.Fund(account, 1_000_000_000_000)
			if fundErr != nil {
				framework.L.Warn().
					Err(fundErr).
					Str("faucet_url", faucetURL).
					Str("account", account.StringLong()).
					Msg("Aptos host faucet fund failed")
			}
			if waitForAptosAccountVisible(client, account, 8*time.Second) == nil {
				framework.L.Info().
					Str("faucet_url", faucetURL).
					Str("account", account.StringLong()).
					Msg("Aptos account funded/visible via host faucet")
				return
			}
			framework.L.Warn().
				Str("faucet_url", faucetURL).
				Str("account", account.StringLong()).
				Msg("Aptos host faucet path did not make account visible; trying container fallback")
		}
	} else {
		framework.L.Warn().
			Err(err).
			Str("node_url", nodeURL).
			Str("account", account.StringLong()).
			Msg("failed to derive Aptos faucet URL; trying container fallback")
	}

	if containerName != "" {
		if cErr := fundAptosAccountInContainer(containerName, account.StringLong()); cErr != nil {
			framework.L.Warn().
				Err(cErr).
				Str("container", containerName).
				Str("account", account.StringLong()).
				Msg("Aptos container faucet fund failed")
		}
		if waitForAptosAccountVisible(client, account, 8*time.Second) == nil {
			framework.L.Info().
				Str("container", containerName).
				Str("account", account.StringLong()).
				Msg("Aptos account funded/visible via container faucet path")
			return
		}
	}

	framework.L.Warn().
		Str("account", account.StringLong()).
		Msg("Aptos account still not visible after host and container funding paths")
}

func ensureAptosWriteWorkersFunded(t *testing.T, aptosChain blockchains.Blockchain, writeDon *crelib.Don) {
	t.Helper()

	aptosBC, ok := aptosChain.(*blockchains_aptos.Blockchain)
	require.True(t, ok, "expected aptos blockchain type")

	nodeURL := aptosBC.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for worker funding")
	nodeURL, err := normalizeAptosNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for worker funding")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	client, err := aptoslib.NewNodeClient(nodeURL, uint8(chainID))
	require.NoError(t, err, "failed to create Aptos client")

	containerName := aptosBC.CtfOutput().ContainerName
	workers, workerErr := writeDon.Workers()
	require.NoError(t, workerErr, "failed to list Aptos write DON workers for funding")
	require.NotEmpty(t, workers, "Aptos write DON workers list is empty")

	for _, worker := range workers {
		addresses, fetchErr := aptosAccountsForWorker(t, worker)
		require.NoError(t, fetchErr, "failed to fetch Aptos key for worker %q", worker.Name)
		require.NotEmpty(t, addresses, "missing Aptos key for worker %q", worker.Name)
		for _, rawAddress := range addresses {
			rawAddress = strings.TrimSpace(rawAddress)
			if rawAddress == "" {
				continue
			}

			var account aptoslib.AccountAddress
			parseErr := account.ParseStringRelaxed(rawAddress)
			require.NoError(t, parseErr, "failed to parse Aptos worker account for worker %q", worker.Name)

			fundAptosAccountBestEffort(t, client, nodeURL, containerName, account)
			require.NoError(
				t,
				waitForAptosAccountVisible(client, account, 45*time.Second),
				"Aptos worker account %s must be visible/funded before write workflow for worker %q",
				account.StringLong(),
				worker.Name,
			)
		}
	}
}

func aptosAccountsForWorker(t *testing.T, worker *crelib.Node) ([]string, error) {
	t.Helper()

	seen := make(map[string]struct{})
	out := make([]string, 0)
	add := func(raw string) {
		addr := strings.TrimSpace(raw)
		if addr == "" {
			return
		}
		if _, ok := seen[addr]; ok {
			return
		}
		seen[addr] = struct{}{}
		out = append(out, addr)
	}

	gqlKeys, gqlErr := worker.Clients.GQLClient.FetchKeys(t.Context(), "APTOS")
	if gqlErr == nil {
		for _, key := range gqlKeys {
			add(key)
		}
	}

	var raw struct {
		Data []struct {
			Attributes struct {
				Account string `json:"account"`
				Address string `json:"address"`
			} `json:"attributes"`
		} `json:"data"`
	}
	restResp, restErr := worker.Clients.RestClient.APIClient.R().
		SetContext(t.Context()).
		SetResult(&raw).
		Get("/v2/keys/aptos")
	if restErr == nil && restResp != nil && restResp.IsSuccess() {
		for _, entry := range raw.Data {
			add(entry.Attributes.Account)
			add(entry.Attributes.Address)
		}
	}

	if len(out) == 0 {
		if gqlErr != nil && restErr != nil {
			return nil, fmt.Errorf("graphql and rest aptos key lookups failed (gql=%v, rest=%v)", gqlErr, restErr)
		}
		if gqlErr != nil {
			return nil, gqlErr
		}
		if restErr != nil {
			return nil, restErr
		}
	}

	return out, nil
}

func aptosFaucetURLFromNodeURL(nodeURL string) (string, error) {
	u, err := url.Parse(nodeURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("empty host in node url %q", nodeURL)
	}
	u.Host = fmt.Sprintf("%s:8081", host)
	u.Path = ""
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func waitForAptosTransactionSuccess(t *testing.T, client *aptoslib.NodeClient, txHash string, label string) {
	t.Helper()
	tx, err := client.WaitForTransaction(txHash)
	require.NoError(t, err, "failed waiting for Aptos tx: %s", label)
	require.True(t, tx.Success, "Aptos tx failed: %s vm_status=%s", label, tx.VmStatus)
}

func fundAptosAccountInContainer(containerName string, account string) error {
	dc, err := framework.NewDockerClient()
	if err != nil {
		return fmt.Errorf("create docker client: %w", err)
	}
	_, err = dc.ExecContainerWithContext(context.Background(), containerName, []string{
		"aptos", "account", "fund-with-faucet",
		"--account", account,
		"--amount", "1000000000000",
	})
	if err != nil {
		return fmt.Errorf("execute aptos faucet fund in %s: %w", containerName, err)
	}
	return nil
}

func waitForAptosAccountVisible(client *aptoslib.NodeClient, account aptoslib.AccountAddress, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		if _, err := client.Account(account); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(1 * time.Second)
	}
	if lastErr != nil {
		return fmt.Errorf("account %s not visible within timeout: %w", account.StringLong(), lastErr)
	}
	return fmt.Errorf("account %s not visible within timeout", account.StringLong())
}

func normalizeAptosNodeURL(nodeURL string) (string, error) {
	parsed, err := url.Parse(nodeURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse Aptos node URL %q: %w", nodeURL, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("Aptos node URL %q must include scheme and host", nodeURL)
	}
	trimmedPath := strings.TrimRight(parsed.Path, "/")
	if trimmedPath == "" {
		parsed.Path = "/v1"
	} else if trimmedPath != "/v1" {
		parsed.Path = trimmedPath + "/v1"
	}
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func workflowRegistryOwnerBytes(t *testing.T, tenv *configuration.TestEnvironment) []byte {
	t.Helper()
	registryChain, ok := tenv.CreEnvironment.Blockchains[0].(*blockchains_evm.Blockchain)
	require.True(t, ok, "registry chain must be EVM")
	rootOwner := registryChain.SethClient.MustGetRootKeyAddress()
	return common.HexToAddress(rootOwner.Hex()).Bytes()
}

func buildAptosDataFeedsBenchmarkPayload() []byte {
	return buildAptosDataFeedsBenchmarkPayloadFor(aptosBenchmarkFeedID(), 123456789)
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

func aptosBenchmarkFeedID() []byte {
	feedID := make([]byte, 32)
	feedID[31] = 1
	return feedID
}

func aptosRoundtripFeedID() []byte {
	feedID := make([]byte, 32)
	feedID[31] = 2
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
