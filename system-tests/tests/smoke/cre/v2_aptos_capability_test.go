package cre

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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
	aptosfeature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/aptos"
	aptoswrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite/config"
	aptoswriteroundtrip_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswriteroundtrip/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const aptosLocalMaxGasAmount uint64 = 200_000
const aptosWorkerFundingAmountOctas uint64 = 1_000_000_000_000
const aptosWorkerMinBalanceOctas uint64 = 100_000_000

// ExecuteAptosTest runs the Aptos CRE suite with the minimum CI scenarios that
// still cover the end-to-end happy path and the expected-failure path.
func ExecuteAptosTest(t *testing.T, tenv *configuration.TestEnvironment) {
	executeAptosScenarios(t, tenv, aptosDefaultScenarios())
}

type aptosScenario struct {
	name string
	run  func(
		t *testing.T,
		tenv *configuration.TestEnvironment,
		aptosChain blockchains.Blockchain,
		userLogsCh <-chan *workflowevents.UserLogs,
		baseMessageCh <-chan *commonevents.BaseMessage,
	)
}

func aptosDefaultScenarios() []aptosScenario {
	return []aptosScenario{
		{name: "Aptos Write Read Roundtrip", run: ExecuteAptosWriteReadRoundtripTest},
		{name: "Aptos Write Expected Failure", run: ExecuteAptosWriteExpectedFailureTest},
	}
}

func executeAptosScenarios(t *testing.T, tenv *configuration.TestEnvironment, scenarios []aptosScenario) {
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

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.run(t, tenv, aptosChain, userLogsCh, baseMessageCh)
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
		// Keep within the current local Aptos transaction max-gas bound.
		MaxGasAmount: aptosLocalMaxGasAmount,
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
		MaxGasAmount:       aptosLocalMaxGasAmount,
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
		MaxGasAmount:       aptosLocalMaxGasAmount,
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

	forwarderHex, ok := aptosfeature.ForwarderAddress(tenv.CreEnvironment.CldfEnvironment.DataStore, aptosChain.ChainSelector())
	require.True(t, ok, "Aptos write test requires forwarder address in datastore for chainSelector=%d", aptosChain.ChainSelector())
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
	nodeURL, err := aptosfeature.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := aptosfeature.ChainIDUint8(chainID)
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
	nodeURL, err := aptosfeature.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := bc.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := aptosfeature.ChainIDUint8(chainID)
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
	nodeURL, err := aptosfeature.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for onchain verification")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := aptosfeature.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")
	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
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
	nodeURL, err := aptosfeature.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for receiver deployment")
	containerName := aptosBC.CtfOutput().ContainerName

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := aptosfeature.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")
	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	deployer, err := aptosDeployerAccount()
	require.NoError(t, err, "failed to create Aptos deployer account")

	require.NoError(
		t,
		aptosfeature.FundAccountBestEffort(t.Context(), framework.L, client, nodeURL, containerName, deployer.AccountAddress(), aptosWorkerMinBalanceOctas, aptosWorkerFundingAmountOctas),
		"failed to fund Aptos deployer account",
	)
	require.NoError(
		t,
		aptosfeature.WaitForAccountVisible(t.Context(), client, deployer.AccountAddress(), 45*time.Second),
		"Aptos deployer account must be visible before deploy",
	)

	var primaryForwarderAddr aptoslib.AccountAddress
	err = primaryForwarderAddr.ParseStringRelaxed(primaryForwarderHex)
	require.NoError(t, err, "failed to parse primary forwarder address")

	owner := deployer.AccountAddress()
	secondaryAddress, secondaryTx, _, err := aptosplatformsecondary.DeployToObject(deployer, client, owner)
	require.NoError(t, err, "failed to deploy Aptos secondary platform package")
	require.NoError(t, aptosfeature.WaitForTransactionSuccess(client, secondaryTx.Hash, "platform_secondary deployment"))

	dataFeedsAddress, dataFeedsTx, dataFeeds, err := aptosdatafeeds.DeployToObject(
		deployer,
		client,
		owner,
		primaryForwarderAddr,
		owner,
		secondaryAddress,
	)
	require.NoError(t, err, "failed to deploy Aptos data feeds receiver package")
	require.NoError(t, aptosfeature.WaitForTransactionSuccess(client, dataFeedsTx.Hash, "data_feeds deployment"))

	workflowOwner := workflowRegistryOwnerBytes(t, tenv)
	tx, err := dataFeeds.Registry().SetWorkflowConfig(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{workflowOwner},
		[][]byte{},
	)
	require.NoError(t, err, "failed to set data feeds workflow config")
	require.NoError(t, aptosfeature.WaitForTransactionSuccess(client, tx.Hash, "data_feeds set_workflow_config"))

	// Configure the feed that the write workflow will update.
	// Without this, registry::perform_update emits WriteSkippedFeedNotSet and benchmark remains unchanged.
	tx, err = dataFeeds.Registry().SetFeeds(
		&aptosbind.TransactOpts{Signer: deployer},
		[][]byte{feedID},
		[]string{"CRE-BENCHMARK"},
		[]byte{0x99},
	)
	require.NoError(t, err, "failed to set data feeds feed config")
	require.NoError(t, aptosfeature.WaitForTransactionSuccess(client, tx.Hash, "data_feeds set_feeds"))

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

	nodeURL := aptosBC.CtfOutput().Nodes[0].ExternalHTTPUrl
	require.NotEmpty(t, nodeURL, "Aptos node URL is required for worker funding")
	nodeURL, err := aptosfeature.NormalizeNodeURL(nodeURL)
	require.NoError(t, err, "failed to normalize Aptos node URL for worker funding")

	chainID := aptosBC.ChainID()
	require.LessOrEqual(t, chainID, uint64(255), "Aptos chain id must fit in uint8")
	chainIDUint8, err := aptosfeature.ChainIDUint8(chainID)
	require.NoError(t, err, "failed to convert Aptos chain id")
	client, err := aptoslib.NewNodeClient(nodeURL, chainIDUint8)
	require.NoError(t, err, "failed to create Aptos client")

	containerName := aptosBC.CtfOutput().ContainerName
	workers, workerErr := writeDon.Workers()
	require.NoError(t, workerErr, "failed to list Aptos write DON workers for funding")
	require.NotEmpty(t, workers, "Aptos write DON workers list is empty")

	for _, worker := range workers {
		addresses, fetchErr := crelib.AptosAccountsForNode(t.Context(), worker)
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

			require.NoError(
				t,
				aptosfeature.FundAccountBestEffort(t.Context(), framework.L, client, nodeURL, containerName, account, aptosWorkerMinBalanceOctas, aptosWorkerFundingAmountOctas),
				"failed to fund Aptos worker account %s for worker %q",
				account.StringLong(),
				worker.Name,
			)
			require.NoError(
				t,
				aptosfeature.WaitForAccountVisible(t.Context(), client, account, 45*time.Second),
				"Aptos worker account %s must be visible/funded before write workflow for worker %q",
				account.StringLong(),
				worker.Name,
			)
		}
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
