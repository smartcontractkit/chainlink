package cre

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	libcrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	evmreadcontracts "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/contracts"
	evm_logTrigger_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/logtrigger/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// -------------------------------------------------------
// LogTriggerGun – emits on-chain events at target RPM
// -------------------------------------------------------

var _ wasp.Gun = (*LogTriggerGun)(nil)

// LogTriggerGun is a wasp Gun that emits EVM log events by calling
// the MessageEmitter contract's EmitMessage function at the configured rate.
type LogTriggerGun struct {
	chain       *evm.Blockchain
	msgEmitter  *evmreadcontracts.MessageEmitter
	logger      zerolog.Logger
	callCounter atomic.Int64
}

func NewLogTriggerGun(
	chain *evm.Blockchain,
	msgEmitter *evmreadcontracts.MessageEmitter,
	logger zerolog.Logger,
) *LogTriggerGun {
	return &LogTriggerGun{
		chain:      chain,
		msgEmitter: msgEmitter,
		logger:     logger,
	}
}

func (g *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	callNum := g.callCounter.Add(1)
	message := fmt.Sprintf("load-test-event-%d-%d", callNum, time.Now().UnixNano())

	tx, err := g.msgEmitter.EmitMessage(g.chain.SethClient.NewTXOpts(), message)
	if err != nil {
		g.logger.Error().Err(err).Msgf("LogTriggerGun: failed to emit message #%d", callNum)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("emit failed: %v", err)}
	}

	_, err = g.chain.SethClient.WaitMined(context.Background(), g.logger, g.chain.SethClient.Client, tx)
	if err != nil {
		g.logger.Error().Err(err).Msgf("LogTriggerGun: tx not mined for call #%d", callNum)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("mine failed: %v", err)}
	}

	g.logger.Info().Msgf("LogTriggerGun: event #%d emitted and mined (tx: %s)", callNum, tx.Hash().Hex())
	return &wasp.Response{}
}

// -------------------------------------------------------
// HTTPTriggerGun – sends HTTP trigger requests at target RPM
// -------------------------------------------------------

var _ wasp.Gun = (*HTTPTriggerGun)(nil)

// HTTPTriggerGun is a wasp Gun that sends HTTP trigger requests to the
// CRE gateway at the configured rate. Each call creates a signed JSON-RPC
// request and posts it to the gateway endpoint.
type HTTPTriggerGun struct {
	gatewayURL           *url.URL
	workflowName         string
	workflowID           string
	workflowOwnerAddress string
	signingKey           *ecdsa.PrivateKey
	logger               zerolog.Logger
	httpClient           *http.Client
	callCounter          atomic.Int64
	mu                   sync.Mutex
}

func NewHTTPTriggerGun(
	gatewayURL *url.URL,
	workflowName string,
	workflowID string,
	workflowOwnerAddress string,
	signingKey *ecdsa.PrivateKey,
	logger zerolog.Logger,
) *HTTPTriggerGun {
	return &HTTPTriggerGun{
		gatewayURL:           gatewayURL,
		workflowName:         workflowName,
		workflowID:           workflowID,
		workflowOwnerAddress: workflowOwnerAddress,
		signingKey:           signingKey,
		logger:               logger,
		httpClient:           &http.Client{Timeout: 4 * time.Minute}, // Must exceed gateway RequestTimeoutMillis (180s) with headroom
	}
}

func (g *HTTPTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	callNum := g.callCounter.Add(1)

	triggerRequest := g.createHTTPTriggerRequest(callNum)
	triggerRequestBody, err := json.Marshal(triggerRequest)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("marshal failed: %v", err)}
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", g.gatewayURL.String(), bytes.NewBuffer(triggerRequestBody))
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("request creation failed: %v", err)}
	}
	req.Header.Set("Content-Type", "application/jsonrpc")
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		g.logger.Warn().Err(err).Msgf("HTTPTriggerGun: request #%d failed", callNum)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("http request failed: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("read body failed: %v", err)}
	}

	if resp.StatusCode != http.StatusOK {
		g.logger.Warn().Msgf("HTTPTriggerGun: request #%d returned status %d: %s", callNum, resp.StatusCode, string(body))
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("status %d: %s", resp.StatusCode, string(body))}
	}

	var jsonResp jsonrpc.Response[json.RawMessage]
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("unmarshal response failed: %v", err)}
	}

	if jsonResp.Error != nil {
		g.logger.Warn().Msgf("HTTPTriggerGun: request #%d JSON-RPC error: %v", callNum, jsonResp.Error)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("jsonrpc error: %v", jsonResp.Error)}
	}

	g.logger.Info().Msgf("HTTPTriggerGun: request #%d succeeded", callNum)
	return &wasp.Response{}
}

func (g *HTTPTriggerGun) createHTTPTriggerRequest(callNum int64) jsonrpc.Request[json.RawMessage] {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: g.workflowOwnerAddress,
			WorkflowName:  g.workflowName,
			WorkflowTag:   "some-tag",
			WorkflowID:    g.workflowID,
		},
		Input: json.RawMessage(fmt.Sprintf(`{
			"customer": "load-test-customer-%d",
			"size": "large",
			"toppings": ["cheese", "pepperoni"],
			"dedupe": false
		}`, callNum)),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal trigger payload: %v", err))
	}
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      fmt.Sprintf("load-test-%d-%s", callNum, uuid.New().String()[0:8]),
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		panic(fmt.Sprintf("failed to create JWT: %v", err))
	}

	tokenString, err := token.SignedString(g.signingKey)
	if err != nil {
		panic(fmt.Sprintf("failed to sign JWT: %v", err))
	}
	req.Auth = tokenString

	return req
}

// -------------------------------------------------------
// Network latency injection (Pumba + tc netem)
// -------------------------------------------------------

const (
	// netemContainerPattern is a Pumba-compatible RE2 regex matching the CRE
	// node containers to apply network delay to.
	netemContainerPattern = `re2:(workflow-node|bootstrap-gateway-node).*`
)

// netemConfig reads network delay parameters from environment variables,
// falling back to defaults. Set NETEM_LATENCY_MS and NETEM_JITTER_MS to
// configure delay without recompiling. Set NETEM_ENABLED=false to disable.
func netemConfig() (enabled bool, latencyMs, jitterMs int) {
	enabled = true
	if v := os.Getenv("NETEM_ENABLED"); v == "false" || v == "0" {
		enabled = false
	}
	latencyMs = 100 // default: 100ms base delay
	if v := os.Getenv("NETEM_LATENCY_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			latencyMs = n
		}
	}
	jitterMs = 50 // default: +/- 50ms jitter
	if v := os.Getenv("NETEM_JITTER_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			jitterMs = n
		}
	}
	return
}

// startNetemDelay uses Pumba to inject tc-netem network delay on all CRE node
// containers matching netemContainerPattern. The delay persists for durationSec
// seconds. Returns a cleanup function that terminates the Pumba sidecar
// container (removing the netem rules).
//
// The function blocks for 5 seconds to allow netem rules to propagate, then
// returns. The Pumba container (and thus the delay rules) remain active for
// the full --duration.
func startNetemDelay(t *testing.T, logger zerolog.Logger, latencyMs, jitterMs, durationSec int) func() {
	t.Helper()

	command := fmt.Sprintf(
		"netem --tc-image=gaiadocker/iproute2 --duration=%ds delay --time=%d --jitter=%d %s",
		durationSec, latencyMs, jitterMs, netemContainerPattern,
	)

	logger.Info().
		Int("latency_ms", latencyMs).
		Int("jitter_ms", jitterMs).
		Int("duration_sec", durationSec).
		Str("container_pattern", netemContainerPattern).
		Str("pumba_command", command).
		Msg("Injecting network delay via Pumba netem on CRE containers")

	// ExecPumba starts a Pumba container, waits for the specified duration,
	// then returns a cleanup function. We use a short wait (5s) so the netem
	// rules have time to propagate, but we don't block the full test duration.
	cleanup, err := chaos.ExecPumba(command, 5*time.Second)
	require.NoError(t, err, "failed to start Pumba netem delay injection")

	logger.Info().Msg("Pumba netem delay injection active – network delay applied to CRE containers")

	return cleanup
}

// -------------------------------------------------------
// Main load test: 400 RPM per trigger (LogTrigger + HTTPTrigger)
// -------------------------------------------------------

// TestLoad_LogTrigger_HTTPTrigger_400RPM validates that the CRE platform can sustain
// 400 RPM per trigger type (400 RPM log trigger + 400 RPM HTTP trigger = 800 RPM total),
// where each individual workflow execution is expected to take ~77 seconds to complete.
//
// The test duration is configurable via durationSeconds. Use shorter durations (e.g. 180s)
// for iterative rate-limit tuning and longer durations (e.g. 600s) for final validation.
//
// Prerequisites:
//   - CRE environment must be started with log-event-trigger-1337 capability enabled.
//     Use the TOML config: workflow-don-load-test-logtrigger-http-kind.toml
//   - Required capability binaries (log-event-trigger, http_trigger, http_action) must be
//     present in core/scripts/cre/environment/binaries/
//   - Job Distributor Docker image must be available (job-distributor:0.22.1)
//   - CL_CRE_SETTINGS_DEFAULT env var must be set on all nodes to raise per-workflow
//     rate limits above the defaults (2 RPM for HTTPTrigger, 10 RPM for LogTrigger).
//
// Run with:
//
//	CTF_CONFIGS=./workflow-don-load-test-logtrigger-http-kind.toml \
//	  go test -v -run TestLoad_LogTrigger_HTTPTrigger_400RPM -timeout 30m
func TestLoad_LogTrigger_HTTPTrigger_400RPM(t *testing.T) {
	testLogger := framework.L

	// Target load parameters – each trigger fires independently at targetRPMPerTrigger
	const (
		targetRPMPerTrigger        = 400 // RPM per individual trigger type
		expectedExecSeconds        = 77  // expected wall-clock time per workflow execution
		durationSeconds            = 180 // test duration (3 min)
		callTimeoutSec             = 300 // must exceed gateway timeout (180s) with headroom for consensus overhead under load
		expectedConcurrentPerTrigger = targetRPMPerTrigger * expectedExecSeconds / 60 // ~513
	)
	duration := time.Duration(durationSeconds) * time.Second
	callTimeout := time.Duration(callTimeoutSec) * time.Second

	testLogger.Info().
		Int("target_rpm_per_trigger", targetRPMPerTrigger).
		Int("expected_exec_seconds", expectedExecSeconds).
		Int("duration_seconds", durationSeconds).
		Int("expected_concurrent_per_trigger", expectedConcurrentPerTrigger).
		Msg("Starting CRE load test: 400 RPM per trigger (EVM LogTrigger + HTTP Trigger)")

	// --------------------------------------------------
	// Phase 1: Set up the CRE environment via CLI
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 1: Setting up CRE test environment...")

	testConfig := t_helpers.GetTestConfig(t, "/configs/workflow-don-load-test-logtrigger-http-kind.toml")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, testConfig, "--with-contracts-version", "v2")

	require.NotEmpty(t, testEnv.CreEnvironment.Blockchains, "expected at least one blockchain")
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	evmChain := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)

	testLogger.Info().Msg("CRE environment setup complete")

	// --------------------------------------------------
	// Phase 2: Deploy MessageEmitter contract for log trigger
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 2: Deploying MessageEmitter contract...")

	msgEmitterAddr, deployTx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(
		evmChain.SethClient.NewTXOpts(), evmChain.SethClient.Client,
	)
	require.NoError(t, err, "failed to deploy MessageEmitter contract")
	_, err = evmChain.SethClient.WaitMined(t.Context(), testLogger, evmChain.SethClient.Client, deployTx)
	require.NoError(t, err, "failed to mine MessageEmitter deployment tx")

	testLogger.Info().Str("address", msgEmitterAddr.Hex()).Msg("MessageEmitter deployed")

	// Get event topic for log trigger configuration
	abiDef, err := evmreadcontracts.MessageEmitterMetaData.GetAbi()
	require.NoError(t, err, "failed to get MessageEmitter ABI")
	eventSigMessageEmitted := abiDef.Events["MessageEmitted"].ID.Hex()

	// --------------------------------------------------
	// Phase 3: Start delayed HTTP server (77s delay per response)
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 3: Starting delayed HTTP server (77s response delay)...")

	publicKeyAddr, signingKey, err := libcrypto.GenerateNewKeyPair()
	require.NoError(t, err, "failed to generate signing key pair")

	backendPort := testEnv.Config.Fake.Port
	// The 77-second delay happens inside the WASM workflow (time.Sleep), not in the backend server.
	// Use a very short backend delay (just enough to be non-trivial) – the real concurrency test
	// comes from the WASM sleep that keeps each workflow execution alive for 77 seconds.
	backendServer := startDelayedHTTPServer(t, backendPort, 0)

	// --------------------------------------------------
	// Phase 4: Deploy EVM LogTrigger + HTTP Action workflow
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 4: Deploying EVM LogTrigger + HTTP Action workflow...")

	logTriggerWorkflowName := fmt.Sprintf("lt-load-%04d", rand.IntN(10000))
	logTriggerWorkflowConfig := evm_logTrigger_config.Config{
		ChainSelector: evmChain.ChainSelector(),
		Addresses:     []string{msgEmitterAddr.Hex()},
		Topics: []struct {
			Values []string `yaml:"values"`
		}{
			{Values: []string{eventSigMessageEmitted}},
		},
		Event: "MessageEmitted",
		Abi:   evmreadcontracts.MessageEmitterMetaData.ABI,
		URL:   backendServer.baseURLDocker + "/events", // HTTP action target – 77s delayed response
	}

	// Use the new log trigger + HTTP action workflow for load testing
	logTriggerWorkflowFile := "../../../../system-tests/tests/load/cre/workflows/logtrigger_http/main.go"
	logTriggerWorkflowID := t_helpers.CompileAndDeployWorkflow(
		t, testEnv, testLogger, logTriggerWorkflowName,
		&logTriggerWorkflowConfig, logTriggerWorkflowFile,
	)
	testLogger.Info().
		Str("name", logTriggerWorkflowName).
		Str("id", logTriggerWorkflowID).
		Msg("LogTrigger + HTTP Action workflow deployed")

	// --------------------------------------------------
	// Phase 5: Deploy HTTP Trigger workflow
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 5: Deploying HTTP Trigger workflow...")

	httpWorkflowName := fmt.Sprintf("ht-load-%04d", rand.IntN(10000))
	httpWorkflowConfig := t_helpers.HTTPWorkflowConfig{
		AuthorizedKey: publicKeyAddr,
		URL:           backendServer.baseURLDocker + "/orders", // HTTP action target – 77s delayed response
	}

	// Use the existing simple HTTP workflow
	httpWorkflowFile := "../../../../core/scripts/cre/environment/examples/workflows/v2/http_simple/main.go"
	httpWorkflowID := t_helpers.CompileAndDeployWorkflow(
		t, testEnv, testLogger, httpWorkflowName,
		&httpWorkflowConfig, httpWorkflowFile,
	)
	testLogger.Info().
		Str("name", httpWorkflowName).
		Str("id", httpWorkflowID).
		Msg("HTTP Trigger workflow deployed")

	// Resolve gateway URL for HTTP trigger requests
	require.NotNil(t, testEnv.Dons, "DON topology must be available")
	gatewayURL := resolveGatewayURL(t, testEnv)
	testLogger.Info().Str("gateway_url", gatewayURL.String()).Msg("Gateway URL resolved")

	workflowOwner := evmChain.SethClient.MustGetRootPrivateKey()
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	// --------------------------------------------------
	// Phase 6: Wait for workflows to be ready
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 6: Waiting for workflows to initialize (up to 5 min)...")

	// Emit a warmup event to ensure the log trigger pipeline is active
	warmupTx, warmupErr := msgEmitter.EmitMessage(evmChain.SethClient.NewTXOpts(), "warmup-event")
	if warmupErr == nil {
		_, _ = evmChain.SethClient.WaitMined(t.Context(), testLogger, evmChain.SethClient.Client, warmupTx)
		testLogger.Info().Msg("Warmup event emitted for log trigger")
	}

	// Wait until the gateway accepts HTTP trigger requests (workflow must be loaded)
	waitForHTTPTriggerReady(t, testLogger, gatewayURL, httpWorkflowName, httpWorkflowID, workflowOwnerAddress, signingKey)

	// --------------------------------------------------
	// Phase 6.5: Inject network latency (optional)
	// --------------------------------------------------
	if enabled, latMs, jitMs := netemConfig(); enabled {
		testLogger.Info().Int("latency_ms", latMs).Int("jitter_ms", jitMs).Msg("Phase 6.5: Injecting network latency via Pumba netem...")
		// Duration must exceed load duration + drain: 180s + 107s + buffer = 600s
		netemCleanup := startNetemDelay(t, testLogger, latMs, jitMs, 600)
		t.Cleanup(netemCleanup)
	} else {
		testLogger.Info().Msg("Phase 6.5: Network latency injection DISABLED (NETEM_ENABLED=false)")
	}

	// --------------------------------------------------
	// Phase 7: Run load test with wasp
	// --------------------------------------------------
	// Each trigger fires at the full target RPM independently
	logTriggerRPM := int64(targetRPMPerTrigger)
	httpTriggerRPM := int64(targetRPMPerTrigger)

	testLogger.Info().
		Int64("log_trigger_rpm", logTriggerRPM).
		Int64("http_trigger_rpm", httpTriggerRPM).
		Int64("total_rpm", logTriggerRPM+httpTriggerRPM).
		Str("duration", duration.String()).
		Msg("Phase 7: Starting load generation (each trigger at 400 RPM, each execution sleeps 77s)...")

	logTriggerGun := NewLogTriggerGun(evmChain, msgEmitter, testLogger)
	httpTriggerGun := NewHTTPTriggerGun(
		gatewayURL, httpWorkflowName, httpWorkflowID,
		workflowOwnerAddress, signingKey, testLogger,
	)

	labels := map[string]string{
		"go_test_name": "logtrigger-http-400rpm",
		"branch":       "load-test",
		"commit":       "load-test",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			CallTimeout: callTimeout,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(logTriggerRPM, duration),
			),
			Gun:                   logTriggerGun,
			Labels:                labels,
			RateLimitUnitDuration: time.Minute,
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			CallTimeout: callTimeout,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(httpTriggerRPM, duration),
			),
			Gun:                   httpTriggerGun,
			Labels:                labels,
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "load test did not finish successfully")

	// --------------------------------------------------
	// Phase 8: Drain – wait for in-flight workflow executions to complete
	// --------------------------------------------------
	// Triggers fired in the last 77 seconds are still executing their HTTP calls.
	// Wait for them to finish so we can count backend hits accurately.
	drainDuration := time.Duration(expectedExecSeconds+30) * time.Second // 77s + 30s buffer
	testLogger.Info().
		Str("drain_duration", drainDuration.String()).
		Int64("backend_hits_before_drain", backendServer.totalHits()).
		Msg("Phase 8: Waiting for in-flight workflow executions to complete...")

	// Poll every 15s and log hit count progress during the drain period
	drainDeadline := time.Now().Add(drainDuration)
	for time.Now().Before(drainDeadline) {
		sleepFor := 15 * time.Second
		if remaining := time.Until(drainDeadline); remaining < sleepFor {
			sleepFor = remaining
		}
		time.Sleep(sleepFor)
		testLogger.Info().
			Int64("orders_hits", backendServer.ordersHits.Load()).
			Int64("events_hits", backendServer.eventsHits.Load()).
			Int64("other_hits", backendServer.otherHits.Load()).
			Int64("total_hits", backendServer.totalHits()).
			Str("remaining", time.Until(drainDeadline).Round(time.Second).String()).
			Msg("Drain progress – backend hit counts")
	}

	// --------------------------------------------------
	// Phase 9: Report results
	// --------------------------------------------------
	totalLogCalls := logTriggerGun.callCounter.Load()
	totalHTTPCalls := httpTriggerGun.callCounter.Load()
	totalCalls := totalLogCalls + totalHTTPCalls
	expectedPerTrigger := int64(targetRPMPerTrigger) * int64(durationSeconds) / 60

	testLogger.Info().
		Int64("log_trigger_calls", totalLogCalls).
		Int64("http_trigger_calls", totalHTTPCalls).
		Int64("total_calls", totalCalls).
		Int64("expected_per_trigger", expectedPerTrigger).
		Int("target_rpm_per_trigger", targetRPMPerTrigger).
		Int("expected_exec_seconds", expectedExecSeconds).
		Int("duration_seconds", durationSeconds).
		Int("expected_concurrent_per_trigger", expectedConcurrentPerTrigger).
		Int64("backend_orders_hits", backendServer.ordersHits.Load()).
		Int64("backend_events_hits", backendServer.eventsHits.Load()).
		Int64("backend_other_hits", backendServer.otherHits.Load()).
		Int64("backend_total_hits", backendServer.totalHits()).
		Msg("Load test completed (77s execution per workflow) – backend hit counts show actual completions")
}

// resolveGatewayURL extracts the gateway URL from the test environment's DON configuration.
func resolveGatewayURL(t *testing.T, testEnv *ttypes.TestEnvironment) *url.URL {
	t.Helper()

	// Try to find gateway from the GatewayConnectors in the Dons topology
	if testEnv.Dons.GatewayConnectors != nil {
		for _, config := range testEnv.Dons.GatewayConnectors.Configurations {
			cfg := config.Incoming
			if cfg.Host != "" {
				rawURL := cfg.Protocol + "://" + cfg.Host + ":" + strconv.Itoa(cfg.ExternalPort) + cfg.Path
				parsed, err := url.Parse(rawURL)
				require.NoError(t, err, "failed to parse gateway URL: %s", rawURL)
				return parsed
			}
		}
	}

	// Fallback: use default gateway URL (localhost:5002)
	defaultURL, err := url.Parse("http://localhost:5002/")
	require.NoError(t, err, "failed to parse default gateway URL")
	return defaultURL
}

// waitForHTTPTriggerReady polls the gateway with HTTP trigger requests until
// the workflow is loaded and the gateway responds with 200 OK. This avoids
// starting the load generator before the gateway recognises the workflow.
func waitForHTTPTriggerReady(
	t *testing.T,
	logger zerolog.Logger,
	gatewayURL *url.URL,
	workflowName, workflowID, workflowOwnerAddress string,
	signingKey *ecdsa.PrivateKey,
) {
	t.Helper()

	gun := NewHTTPTriggerGun(gatewayURL, workflowName, workflowID, workflowOwnerAddress, signingKey, logger)
	tick := 5 * time.Second

	logger.Info().Msg("Waiting for gateway to load the HTTP trigger workflow...")
	require.Eventually(t, func() bool {
		resp := gun.Call(nil)
		if resp.Failed {
			logger.Warn().Msgf("HTTP trigger readiness probe failed: %s", resp.Error)
			return false
		}
		logger.Info().Msg("HTTP trigger workflow is ready – gateway returned 200 OK")
		return true
	}, 5*time.Minute, tick, "gateway did not become ready for HTTP trigger workflow within 5 minutes")
}

// delayedServer is an HTTP server that sleeps for a configurable duration before
// responding to every request. This simulates slow backend processing so that each
// workflow execution takes ~77 seconds end-to-end.
type delayedServer struct {
	server        *http.Server
	delay         time.Duration
	baseURLDocker string
	ordersHits    atomic.Int64
	eventsHits    atomic.Int64
	otherHits     atomic.Int64
}

// startDelayedHTTPServer starts an HTTP server that delays every response by the
// given duration. It tracks hit counts per endpoint (/orders, /events, other).
func startDelayedHTTPServer(t *testing.T, port int, delay time.Duration) *delayedServer {
	t.Helper()

	ds := &delayedServer{
		delay:         delay,
		baseURLDocker: fmt.Sprintf("%s:%d", framework.HostDockerInternal(), port),
	}

	mux := http.NewServeMux()

	writeJSON := func(w http.ResponseWriter, data map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(data)
	}

	// IMPORTANT: All endpoints return STATIC (deterministic) responses.
	// Because the workflow uses ConsensusIdenticalAggregation, every node must
	// observe the same HTTP response for consensus to succeed.
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		ds.ordersHits.Add(1)
		time.Sleep(delay)
		writeJSON(w, map[string]string{
			"status":  "success",
			"message": "Order processed",
		})
	})

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		ds.eventsHits.Add(1)
		time.Sleep(delay)
		writeJSON(w, map[string]string{
			"status":  "received",
			"message": "Event processed",
		})
	})

	// Catch-all for requests to the root or any other path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ds.otherHits.Add(1)
		time.Sleep(delay)
		writeJSON(w, map[string]string{
			"status":  "success",
			"message": "Order processed",
		})
	})

	ds.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := ds.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("delayed HTTP server error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	t.Cleanup(func() {
		_ = ds.server.Close()
	})

	framework.L.Info().
		Int("port", port).
		Str("delay", delay.String()).
		Str("docker_url", ds.baseURLDocker).
		Msg("Delayed HTTP server started (endpoints: /orders, /events, /)")

	return ds
}

func (ds *delayedServer) totalHits() int64 {
	return ds.ordersHits.Load() + ds.eventsHits.Load() + ds.otherHits.Load()
}

// -------------------------------------------------------
// AsyncHTTPTriggerGun – fires HTTP triggers without waiting for workflow completion
// -------------------------------------------------------

var _ wasp.Gun = (*AsyncHTTPTriggerGun)(nil)

// AsyncHTTPTriggerGun sends HTTP trigger requests with a short timeout (just enough
// for the gateway to accept and dispatch the trigger). It does NOT wait for the
// full workflow execution (77s). The workflow continues running in the DON after
// our client disconnects. Actual completions are measured via backend hit counts.
type AsyncHTTPTriggerGun struct {
	gatewayURL           *url.URL
	workflowName         string
	workflowID           string
	workflowOwnerAddress string
	signingKey           *ecdsa.PrivateKey
	logger               zerolog.Logger
	httpClient           *http.Client
	callCounter          atomic.Int64
	mu                   sync.Mutex

	// Counters for async dispatch tracking
	dispatchedCount atomic.Int64 // requests where gateway started processing (got any HTTP response or timeout)
	acceptedCount   atomic.Int64 // requests where gateway returned 200 before our short timeout
	rejectedCount   atomic.Int64 // requests where gateway returned non-200 error
	timedOutCount   atomic.Int64 // requests where our client timed out (workflow likely still running)
	sendErrorCount  atomic.Int64 // requests that failed to send at all
}

func NewAsyncHTTPTriggerGun(
	gatewayURL *url.URL,
	workflowName string,
	workflowID string,
	workflowOwnerAddress string,
	signingKey *ecdsa.PrivateKey,
	logger zerolog.Logger,
) *AsyncHTTPTriggerGun {
	return &AsyncHTTPTriggerGun{
		gatewayURL:           gatewayURL,
		workflowName:         workflowName,
		workflowID:           workflowID,
		workflowOwnerAddress: workflowOwnerAddress,
		signingKey:           signingKey,
		logger:               logger,
		// Short timeout: just enough for the gateway to accept the request.
		// The workflow executes asynchronously in the DON after this.
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (g *AsyncHTTPTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	callNum := g.callCounter.Add(1)

	triggerRequest := g.createHTTPTriggerRequest(callNum)
	triggerRequestBody, err := json.Marshal(triggerRequest)
	if err != nil {
		g.sendErrorCount.Add(1)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("marshal failed: %v", err)}
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", g.gatewayURL.String(), bytes.NewBuffer(triggerRequestBody))
	if err != nil {
		g.sendErrorCount.Add(1)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("request creation failed: %v", err)}
	}
	req.Header.Set("Content-Type", "application/jsonrpc")
	req.Header.Set("Accept", "application/json")

	g.dispatchedCount.Add(1)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		// Timeout is expected – the workflow is still running in the DON.
		// This is NOT a failure for async dispatch.
		if strings.Contains(err.Error(), "Client.Timeout") || strings.Contains(err.Error(), "context deadline exceeded") {
			g.timedOutCount.Add(1)
			if callNum%100 == 0 {
				g.logger.Debug().Msgf("AsyncHTTPTriggerGun: request #%d timed out (expected – workflow dispatched)", callNum)
			}
			// Report as success – the request was sent and likely dispatched
			return &wasp.Response{}
		}
		g.sendErrorCount.Add(1)
		g.logger.Warn().Err(err).Msgf("AsyncHTTPTriggerGun: request #%d send error", callNum)
		return &wasp.Response{Failed: true, Error: fmt.Sprintf("http error: %v", err)}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		g.acceptedCount.Add(1)
		if callNum%100 == 0 {
			g.logger.Info().Msgf("AsyncHTTPTriggerGun: request #%d completed inline (workflow finished within 5s!)", callNum)
		}
		return &wasp.Response{}
	}

	// Non-200 means the gateway rejected the request (workflow was NOT dispatched)
	g.rejectedCount.Add(1)
	if callNum%100 == 0 {
		g.logger.Warn().Msgf("AsyncHTTPTriggerGun: request #%d rejected: %d %s", callNum, resp.StatusCode, string(body))
	}
	return &wasp.Response{Failed: true, Error: fmt.Sprintf("rejected: %d", resp.StatusCode)}
}

func (g *AsyncHTTPTriggerGun) createHTTPTriggerRequest(callNum int64) jsonrpc.Request[json.RawMessage] {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: g.workflowOwnerAddress,
			WorkflowName:  g.workflowName,
			WorkflowTag:   "some-tag",
			WorkflowID:    g.workflowID,
		},
		Input: json.RawMessage(fmt.Sprintf(`{
			"customer": "async-load-test-%d",
			"size": "large",
			"toppings": ["cheese", "pepperoni"],
			"dedupe": false
		}`, callNum)),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal trigger payload: %v", err))
	}
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      fmt.Sprintf("async-load-%d-%s", callNum, uuid.New().String()[0:8]),
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		panic(fmt.Sprintf("failed to create JWT: %v", err))
	}

	tokenString, err := token.SignedString(g.signingKey)
	if err != nil {
		panic(fmt.Sprintf("failed to sign JWT: %v", err))
	}
	req.Auth = tokenString

	return req
}

func (g *AsyncHTTPTriggerGun) logStats(logger zerolog.Logger) {
	logger.Info().
		Int64("total_calls", g.callCounter.Load()).
		Int64("dispatched", g.dispatchedCount.Load()).
		Int64("accepted_inline", g.acceptedCount.Load()).
		Int64("timed_out_dispatched", g.timedOutCount.Load()).
		Int64("rejected", g.rejectedCount.Load()).
		Int64("send_errors", g.sendErrorCount.Load()).
		Msg("AsyncHTTPTriggerGun dispatch stats")
}

// -------------------------------------------------------
// Async HTTP Trigger Load Test: fire-and-forget at 400 RPM
// -------------------------------------------------------

// TestLoad_HTTPTrigger_Async_400RPM validates how many workflow executions actually
// complete when HTTP triggers are fired asynchronously at 400 RPM. Unlike the
// synchronous test, this does NOT hold connections open for 77 seconds per request.
// Instead, it fires requests with a short timeout (5s) and measures completions
// via backend echo server hit counts.
//
// This test isolates the HTTP trigger case to determine if the CRE workflow engine
// can handle the throughput when the gateway connection bottleneck is removed.
//
// Run with:
//
//	CTF_CONFIGS=./workflow-don-load-test-logtrigger-http-kind.toml \
//	  go test -v -run TestLoad_HTTPTrigger_Async_400RPM -timeout 30m
func TestLoad_HTTPTrigger_Async_400RPM(t *testing.T) {
	testLogger := framework.L

	const (
		targetRPM           = 400 // RPM for HTTP trigger
		expectedExecSeconds = 77  // each workflow sleeps 77s in WASM
		loadDurationSeconds = 180 // 3 min load generation – Experiment 6
		drainBufferSeconds  = 120 // extra buffer beyond 77s for consensus overhead
		callTimeoutSec      = 30  // short timeout for async dispatch
	)
	loadDuration := time.Duration(loadDurationSeconds) * time.Second
	callTimeout := time.Duration(callTimeoutSec) * time.Second
	drainDuration := time.Duration(expectedExecSeconds+drainBufferSeconds) * time.Second

	expectedTotalRequests := int64(targetRPM) * int64(loadDurationSeconds) / 60

	testLogger.Info().
		Int("target_rpm", targetRPM).
		Int("expected_exec_seconds", expectedExecSeconds).
		Int("load_duration_seconds", loadDurationSeconds).
		Int64("expected_total_requests", expectedTotalRequests).
		Str("drain_duration", drainDuration.String()).
		Msg("Starting ASYNC HTTP trigger load test – fire-and-forget at 400 RPM")

	// --------------------------------------------------
	// Phase 1: Set up the CRE environment
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 1: Setting up CRE test environment...")

	testConfig := t_helpers.GetTestConfig(t, "/configs/workflow-don-load-test-logtrigger-http-kind.toml")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, testConfig, "--with-contracts-version", "v2")

	require.NotEmpty(t, testEnv.CreEnvironment.Blockchains, "expected at least one blockchain")
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	evmChain := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)

	testLogger.Info().Msg("CRE environment setup complete")

	// --------------------------------------------------
	// Phase 2: Start backend HTTP server (no delay – 77s sleep is in WASM)
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 2: Starting backend HTTP server (no delay)...")

	publicKeyAddr, signingKey, err := libcrypto.GenerateNewKeyPair()
	require.NoError(t, err, "failed to generate signing key pair")

	backendPort := testEnv.Config.Fake.Port
	backendServer := startDelayedHTTPServer(t, backendPort, 0)

	// --------------------------------------------------
	// Phase 3: Deploy HTTP Trigger workflow
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 3: Deploying HTTP Trigger workflow...")

	httpWorkflowName := fmt.Sprintf("ha-load-%04d", rand.IntN(10000))
	httpWorkflowConfig := t_helpers.HTTPWorkflowConfig{
		AuthorizedKey: publicKeyAddr,
		URL:           backendServer.baseURLDocker + "/orders",
	}

	httpWorkflowFile := "../../../../core/scripts/cre/environment/examples/workflows/v2/http_simple/main.go"
	httpWorkflowID := t_helpers.CompileAndDeployWorkflow(
		t, testEnv, testLogger, httpWorkflowName,
		&httpWorkflowConfig, httpWorkflowFile,
	)
	testLogger.Info().
		Str("name", httpWorkflowName).
		Str("id", httpWorkflowID).
		Msg("HTTP Trigger workflow deployed")

	// Resolve gateway URL
	require.NotNil(t, testEnv.Dons, "DON topology must be available")
	gatewayURL := resolveGatewayURL(t, testEnv)
	testLogger.Info().Str("gateway_url", gatewayURL.String()).Msg("Gateway URL resolved")

	workflowOwner := evmChain.SethClient.MustGetRootPrivateKey()
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	// --------------------------------------------------
	// Phase 4: Wait for workflow to be ready
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 4: Waiting for workflow to initialize (up to 5 min)...")
	waitForHTTPTriggerReady(t, testLogger, gatewayURL, httpWorkflowName, httpWorkflowID, workflowOwnerAddress, signingKey)

	// --------------------------------------------------
	// Phase 5: Run async load test
	// --------------------------------------------------
	testLogger.Info().
		Int("target_rpm", targetRPM).
		Str("load_duration", loadDuration.String()).
		Msg("Phase 5: Starting ASYNC load generation (fire-and-forget at 400 RPM)...")

	asyncGun := NewAsyncHTTPTriggerGun(
		gatewayURL, httpWorkflowName, httpWorkflowID,
		workflowOwnerAddress, signingKey, testLogger,
	)

	labels := map[string]string{
		"go_test_name": "http-trigger-async-400rpm",
		"branch":       "load-test",
		"commit":       "load-test",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			CallTimeout: callTimeout,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(int64(targetRPM), loadDuration),
			),
			Gun:                   asyncGun,
			Labels:                labels,
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "async load test did not finish")

	// Log dispatch stats right after load generation
	asyncGun.logStats(testLogger)

	// --------------------------------------------------
	// Phase 6: Drain – wait for in-flight workflows to complete
	// --------------------------------------------------
	testLogger.Info().
		Str("drain_duration", drainDuration.String()).
		Int64("backend_hits_before_drain", backendServer.totalHits()).
		Msg("Phase 6: Waiting for in-flight async workflows to complete...")

	// Poll every 15s and log hit count progress during drain
	drainDeadline := time.Now().Add(drainDuration)
	for time.Now().Before(drainDeadline) {
		sleepFor := 15 * time.Second
		if remaining := time.Until(drainDeadline); remaining < sleepFor {
			sleepFor = remaining
		}
		time.Sleep(sleepFor)
		testLogger.Info().
			Int64("orders_hits", backendServer.ordersHits.Load()).
			Int64("events_hits", backendServer.eventsHits.Load()).
			Int64("other_hits", backendServer.otherHits.Load()).
			Int64("total_hits", backendServer.totalHits()).
			Str("remaining", time.Until(drainDeadline).Round(time.Second).String()).
			Msg("Drain progress – backend hit counts (actual workflow completions)")
	}

	// --------------------------------------------------
	// Phase 7: Report results
	// --------------------------------------------------
	totalCalls := asyncGun.callCounter.Load()
	ordersHits := backendServer.ordersHits.Load()
	completionRate := float64(0)
	if totalCalls > 0 {
		completionRate = float64(ordersHits) / float64(totalCalls) * 100
	}

	asyncGun.logStats(testLogger)

	testLogger.Info().
		Int64("total_requests_sent", totalCalls).
		Int64("dispatched", asyncGun.dispatchedCount.Load()).
		Int64("timed_out_dispatched", asyncGun.timedOutCount.Load()).
		Int64("accepted_inline", asyncGun.acceptedCount.Load()).
		Int64("rejected_by_gateway", asyncGun.rejectedCount.Load()).
		Int64("send_errors", asyncGun.sendErrorCount.Load()).
		Int64("backend_orders_hits", ordersHits).
		Int64("backend_events_hits", backendServer.eventsHits.Load()).
		Int64("backend_total_hits", backendServer.totalHits()).
		Float64("completion_rate_pct", completionRate).
		Str("completion_rate", fmt.Sprintf("%.1f%% (%d/%d)", completionRate, ordersHits, totalCalls)).
		Int("target_rpm", targetRPM).
		Int("expected_exec_seconds", expectedExecSeconds).
		Int("load_duration_seconds", loadDurationSeconds).
		Msg("ASYNC HTTP Trigger Load Test COMPLETE – completion rate shows actual end-to-end workflow successes")
}

// TestLoad_LogTrigger_Only_400RPM is a focused load test that exercises ONLY the
// EVM log trigger workflow at 400 RPM for 3 minutes. This isolates the log trigger
// pipeline (LogPoller → trigger capability → workflow engine → HTTP action → consensus)
// from the HTTP trigger path to measure log trigger throughput independently.
//
// Run:
//
//	CTF_CONFIGS=./workflow-don-load-test-logtrigger-http-kind.toml \
//	  go test -v -run TestLoad_LogTrigger_Only_400RPM -timeout 30m
func TestLoad_LogTrigger_Only_400RPM(t *testing.T) {
	testLogger := framework.L

	const (
		targetRPM           = 400 // RPM for log trigger
		durationSeconds     = 180 // 3 min load generation
		callTimeoutSec      = 300 // headroom for slow tx mining under load
		drainBufferSeconds  = 120 // wait for in-flight workflows after load stops
	)
	duration := time.Duration(durationSeconds) * time.Second
	callTimeout := time.Duration(callTimeoutSec) * time.Second
	drainDuration := time.Duration(drainBufferSeconds) * time.Second

	expectedTotalEvents := int64(targetRPM) * int64(durationSeconds) / 60

	testLogger.Info().
		Int("target_rpm", targetRPM).
		Int("duration_seconds", durationSeconds).
		Int64("expected_total_events", expectedTotalEvents).
		Str("drain_duration", drainDuration.String()).
		Msg("Starting FOCUSED log trigger load test – EVM LogTrigger only at 400 RPM")

	// --------------------------------------------------
	// Phase 1: Set up the CRE environment via CLI
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 1: Setting up CRE test environment...")

	testConfig := t_helpers.GetTestConfig(t, "/configs/workflow-don-load-test-logtrigger-http-kind.toml")
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, testConfig, "--with-contracts-version", "v2")

	require.NotEmpty(t, testEnv.CreEnvironment.Blockchains, "expected at least one blockchain")
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain")
	evmChain := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain)

	testLogger.Info().Msg("CRE environment setup complete")

	// --------------------------------------------------
	// Phase 2: Deploy MessageEmitter contract for log trigger
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 2: Deploying MessageEmitter contract...")

	msgEmitterAddr, deployTx, msgEmitter, err := evmreadcontracts.DeployMessageEmitter(
		evmChain.SethClient.NewTXOpts(), evmChain.SethClient.Client,
	)
	require.NoError(t, err, "failed to deploy MessageEmitter contract")
	_, err = evmChain.SethClient.WaitMined(t.Context(), testLogger, evmChain.SethClient.Client, deployTx)
	require.NoError(t, err, "failed to mine MessageEmitter deployment tx")

	testLogger.Info().Str("address", msgEmitterAddr.Hex()).Msg("MessageEmitter deployed")

	// Get event topic for log trigger configuration
	abiDef, err := evmreadcontracts.MessageEmitterMetaData.GetAbi()
	require.NoError(t, err, "failed to get MessageEmitter ABI")
	eventSigMessageEmitted := abiDef.Events["MessageEmitted"].ID.Hex()

	// --------------------------------------------------
	// Phase 3: Start backend HTTP server (no delay)
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 3: Starting backend HTTP server (no delay)...")

	backendPort := testEnv.Config.Fake.Port
	backendServer := startDelayedHTTPServer(t, backendPort, 0)

	// --------------------------------------------------
	// Phase 4: Deploy EVM LogTrigger + HTTP Action workflow
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 4: Deploying EVM LogTrigger + HTTP Action workflow...")

	logTriggerWorkflowName := fmt.Sprintf("lt-only-%04d", rand.IntN(10000))
	logTriggerWorkflowConfig := evm_logTrigger_config.Config{
		ChainSelector: evmChain.ChainSelector(),
		Addresses:     []string{msgEmitterAddr.Hex()},
		Topics: []struct {
			Values []string `yaml:"values"`
		}{
			{Values: []string{eventSigMessageEmitted}},
		},
		Event: "MessageEmitted",
		Abi:   evmreadcontracts.MessageEmitterMetaData.ABI,
		URL:   backendServer.baseURLDocker + "/events",
	}

	logTriggerWorkflowFile := "../../../../system-tests/tests/load/cre/workflows/logtrigger_http/main.go"
	logTriggerWorkflowID := t_helpers.CompileAndDeployWorkflow(
		t, testEnv, testLogger, logTriggerWorkflowName,
		&logTriggerWorkflowConfig, logTriggerWorkflowFile,
	)
	testLogger.Info().
		Str("name", logTriggerWorkflowName).
		Str("id", logTriggerWorkflowID).
		Msg("LogTrigger + HTTP Action workflow deployed")

	// --------------------------------------------------
	// Phase 5: Warmup – emit a few events and wait for pipeline to be active
	// --------------------------------------------------
	testLogger.Info().Msg("Phase 5: Warming up log trigger pipeline...")

	for i := 0; i < 3; i++ {
		warmupTx, warmupErr := msgEmitter.EmitMessage(evmChain.SethClient.NewTXOpts(), fmt.Sprintf("warmup-%d", i))
		if warmupErr == nil {
			_, _ = evmChain.SethClient.WaitMined(t.Context(), testLogger, evmChain.SethClient.Client, warmupTx)
		}
		time.Sleep(2 * time.Second)
	}

	// Wait a bit for log poller to pick up warmup events and workflow engine to process them
	testLogger.Info().Msg("Waiting 30s for log trigger pipeline to stabilize...")
	time.Sleep(30 * time.Second)

	warmupEventsHits := backendServer.eventsHits.Load()
	testLogger.Info().
		Int64("warmup_backend_hits", warmupEventsHits).
		Msg("Warmup complete")

	// --------------------------------------------------
	// Phase 5.5: Inject network latency (optional)
	// --------------------------------------------------
	if enabled, latMs, jitMs := netemConfig(); enabled {
		testLogger.Info().Int("latency_ms", latMs).Int("jitter_ms", jitMs).Msg("Phase 5.5: Injecting network latency via Pumba netem...")
		// Duration must exceed load duration + drain: 180s + 120s + buffer = 600s
		netemCleanup := startNetemDelay(t, testLogger, latMs, jitMs, 600)
		t.Cleanup(netemCleanup)
	} else {
		testLogger.Info().Msg("Phase 5.5: Network latency injection DISABLED (NETEM_ENABLED=false)")
	}

	// --------------------------------------------------
	// Phase 6: Run load test with wasp – log trigger only
	// --------------------------------------------------
	testLogger.Info().
		Int("target_rpm", targetRPM).
		Str("duration", duration.String()).
		Msg("Phase 6: Starting load generation (EVM LogTrigger only at 400 RPM)...")

	logTriggerGun := NewLogTriggerGun(evmChain, msgEmitter, testLogger)

	labels := map[string]string{
		"go_test_name": "logtrigger-only-400rpm",
		"branch":       "load-test",
		"commit":       "load-test",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			CallTimeout: callTimeout,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(int64(targetRPM), duration),
			),
			Gun:                   logTriggerGun,
			Labels:                labels,
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "log trigger load test did not finish")

	// --------------------------------------------------
	// Phase 7: Drain – wait for in-flight workflow executions to complete
	// --------------------------------------------------
	testLogger.Info().
		Str("drain_duration", drainDuration.String()).
		Int64("backend_events_hits_before_drain", backendServer.eventsHits.Load()).
		Msg("Phase 7: Waiting for in-flight log trigger workflows to complete...")

	drainDeadline := time.Now().Add(drainDuration)
	for time.Now().Before(drainDeadline) {
		sleepFor := 15 * time.Second
		if remaining := time.Until(drainDeadline); remaining < sleepFor {
			sleepFor = remaining
		}
		time.Sleep(sleepFor)
		testLogger.Info().
			Int64("events_hits", backendServer.eventsHits.Load()).
			Int64("orders_hits", backendServer.ordersHits.Load()).
			Int64("other_hits", backendServer.otherHits.Load()).
			Int64("total_hits", backendServer.totalHits()).
			Str("remaining", time.Until(drainDeadline).Round(time.Second).String()).
			Msg("Drain progress – backend hit counts (actual workflow completions)")
	}

	// --------------------------------------------------
	// Phase 8: Report results
	// --------------------------------------------------
	totalEventsEmitted := logTriggerGun.callCounter.Load()
	eventsHits := backendServer.eventsHits.Load() - warmupEventsHits // subtract warmup
	completionRate := float64(0)
	// Each node makes its own HTTP call, so divide by 4 (N nodes) for unique completions
	uniqueCompletions := eventsHits / 4
	if totalEventsEmitted > 0 {
		completionRate = float64(uniqueCompletions) / float64(totalEventsEmitted) * 100
	}

	testLogger.Info().
		Int64("total_events_emitted", totalEventsEmitted).
		Int64("expected_total_events", expectedTotalEvents).
		Int64("backend_events_hits_raw", eventsHits).
		Int64("unique_completions_approx", uniqueCompletions).
		Float64("completion_rate_pct", completionRate).
		Str("completion_rate", fmt.Sprintf("%.1f%% (%d/%d events)", completionRate, uniqueCompletions, totalEventsEmitted)).
		Int64("backend_orders_hits", backendServer.ordersHits.Load()).
		Int64("backend_other_hits", backendServer.otherHits.Load()).
		Int64("backend_total_hits", backendServer.totalHits()).
		Int("target_rpm", targetRPM).
		Int("duration_seconds", durationSeconds).
		Msg("LOG TRIGGER ONLY Load Test COMPLETE – events_hits show actual end-to-end workflow completions")
}
