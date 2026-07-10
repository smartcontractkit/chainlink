package cre

// Standalone load driver for the confidential-workflows enclaves on a live
// (staging) DON. It fires the HTTP-triggered confidential load-test workflow
// (chainlink-deployments confidential_workflows_load, cn_confidential_workflows_load_a)
// at a controlled, capped, stepped request rate by POSTing workflows.execute
// JSON-RPC to the gateway user server, and reports round-trip latency
// percentiles + success/failure.
//
// It is deliberately NOT built on the wasp harness: that harness stands up a
// local CTF env with mock capabilities and lives on an unmerged branch. This
// driver points at an already-deployed DON and needs nothing but the gateway
// URL and the load-test signing key, so it can run against staging as-is.
//
// SAFETY: staging DON/gateway/vault are shared CRE infra. The rate is capped
// (LOADTEST_MAX_RPS) and ramps in steps, and the run self-aborts if the rolling
// failure rate exceeds a threshold. Keep the cap low and watch the enclave
// metrics (enclave_execution_time_ms / enclave_execution_failures) alongside.
//
// The driver cannot see which enclave served each request (that is not in the
// HTTP response); per-enclave spread is confirmed out-of-band via the enclave
// metrics/dashboards.
//
// Run (skips unless the required env is set):
//
//	LOADTEST_GATEWAY_URL="https://cre-gateway.cre.stage.external.griddle.sh:5002/" \
//	CRE_LOADTEST_PRIVATE_KEY="0x..." \
//	LOADTEST_WORKFLOW_OWNER="0x..." LOADTEST_WORKFLOW_NAME="cn_confidential_workflows_load_a" \
//	go test -run TestLoad_ConfidentialWorkflows_HTTPTrigger -timeout 30m ./load/cre/

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type loadConfig struct {
	gatewayURL  string
	privKey     *ecdsa.PrivateKey
	owner       string
	name        string
	tag         string
	id          string
	input       string
	callTimeout time.Duration

	steps        []rpsStep
	abortFailPct float64
	abortWindow  int
}

type rpsStep struct {
	rps int
	dur time.Duration
}

func TestLoad_ConfidentialWorkflows_HTTPTrigger(t *testing.T) {
	cfg := loadConfigFromEnv(t)

	signer := crypto.PubkeyToAddress(cfg.privKey.PublicKey).Hex()
	t.Logf("confidential-workflows load test")
	t.Logf("  gateway     : %s", cfg.gatewayURL)
	t.Logf("  workflow    : owner=%s name=%q tag=%q id=%q", cfg.owner, cfg.name, cfg.tag, cfg.id)
	t.Logf("  signer addr : %s (must match the workflow authKey)", signer)
	t.Logf("  schedule    : %s", describeSteps(cfg.steps))
	t.Logf("  callTimeout : %s   abort if >%.0f%% of last %d fail", cfg.callTimeout, cfg.abortFailPct, cfg.abortWindow)

	client := &http.Client{Timeout: cfg.callTimeout}

	var (
		mu        sync.Mutex
		latencies []time.Duration
		success   int
		failed    int
		window    []bool
		aborted   bool
		wg        sync.WaitGroup
	)

	// record returns true once the run should abort (shared-infra protection).
	record := func(ok bool, d time.Duration) {
		mu.Lock()
		defer mu.Unlock()
		if ok {
			success++
			latencies = append(latencies, d)
		} else {
			failed++
		}
		window = append(window, ok)
		if len(window) > cfg.abortWindow {
			window = window[len(window)-cfg.abortWindow:]
		}
		if !aborted && len(window) >= cfg.abortWindow {
			fails := 0
			for _, ok := range window {
				if !ok {
					fails++
				}
			}
			if 100*float64(fails)/float64(len(window)) > cfg.abortFailPct {
				aborted = true
			}
		}
	}
	isAborted := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return aborted
	}

	fire := func() {
		defer wg.Done()
		start := time.Now()
		ok := doExecute(t.Context(), client, cfg)
		record(ok, time.Since(start))
	}

	overallStart := time.Now()
steps:
	for _, step := range cfg.steps {
		t.Logf("--> %d rps for %s", step.rps, step.dur)
		interval := time.Second / time.Duration(step.rps)
		ticker := time.NewTicker(interval)
		deadline := time.Now().Add(step.dur)
		for time.Now().Before(deadline) {
			<-ticker.C
			if isAborted() {
				ticker.Stop()
				t.Logf("ABORT: rolling failure rate exceeded %.0f%%, stopping ramp", cfg.abortFailPct)
				break steps
			}
			wg.Add(1)
			go fire()
		}
		ticker.Stop()
	}
	wg.Wait()
	elapsed := time.Since(overallStart)

	total := success + failed
	t.Logf("==================== results ====================")
	t.Logf("  duration    : %s", elapsed.Round(time.Millisecond))
	t.Logf("  requests    : %d  (success=%d failed=%d)", total, success, failed)
	if total > 0 {
		t.Logf("  success rate: %.1f%%", 100*float64(success)/float64(total))
		t.Logf("  throughput  : %.2f req/s (completed)", float64(total)/elapsed.Seconds())
	}
	if len(latencies) > 0 {
		t.Logf("  latency p50 : %s", pct(latencies, 50).Round(time.Millisecond))
		t.Logf("  latency p90 : %s", pct(latencies, 90).Round(time.Millisecond))
		t.Logf("  latency p95 : %s", pct(latencies, 95).Round(time.Millisecond))
		t.Logf("  latency p99 : %s", pct(latencies, 99).Round(time.Millisecond))
		t.Logf("  latency max : %s", pct(latencies, 100).Round(time.Millisecond))
	}
	t.Logf("================================================")

	require.NotZero(t, total, "no requests were sent")
	require.NotZero(t, success, "every request failed, check gateway URL, workflow selector, and that the signer matches authKey")
}

// doExecute builds a signed workflows.execute request and POSTs it to the
// gateway, returning true only on HTTP 200 with no JSON-RPC error.
func doExecute(ctx context.Context, client *http.Client, cfg loadConfig) bool {
	req, err := buildExecuteRequest(cfg)
	if err != nil {
		return false
	}
	body, err := json.Marshal(req)
	if err != nil {
		return false
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.gatewayURL, bytes.NewReader(body))
	if err != nil {
		return false
	}
	httpReq.Header.Set("Content-Type", "application/jsonrpc")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	var out jsonrpc.Response[json.RawMessage]
	if err := json.Unmarshal(rb, &out); err != nil {
		return false
	}
	return out.Error == nil
}

// buildExecuteRequest mirrors the smoke test's createHTTPTriggerRequestWithKey:
// a workflows.execute JSON-RPC request authed by a JWT the load-test key signs.
func buildExecuteRequest(cfg loadConfig) (jsonrpc.Request[json.RawMessage], error) {
	payload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: cfg.owner,
			WorkflowName:  cfg.name,
			WorkflowTag:   cfg.tag,
			WorkflowID:    cfg.id,
		},
		Input: json.RawMessage(cfg.input),
	}
	pb, err := json.Marshal(payload)
	if err != nil {
		return jsonrpc.Request[json.RawMessage]{}, err
	}
	raw := json.RawMessage(pb)
	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &raw,
		ID:      "cwload-" + uuid.New().String()[:8],
	}
	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		return jsonrpc.Request[json.RawMessage]{}, err
	}
	signed, err := token.SignedString(cfg.privKey)
	if err != nil {
		return jsonrpc.Request[json.RawMessage]{}, err
	}
	req.Auth = signed
	return req, nil
}

// pct returns the p-th percentile (p in [0,100]) of a copy of ds.
func pct(ds []time.Duration, p int) time.Duration {
	s := append([]time.Duration(nil), ds...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	if p >= 100 {
		return s[len(s)-1]
	}
	idx := p * (len(s) - 1) / 100
	return s[idx]
}

func describeSteps(steps []rpsStep) string {
	parts := make([]string, len(steps))
	for i, s := range steps {
		parts[i] = fmt.Sprintf("%drps/%s", s.rps, s.dur)
	}
	return strings.Join(parts, " -> ")
}

func loadConfigFromEnv(t *testing.T) loadConfig {
	gw := os.Getenv("LOADTEST_GATEWAY_URL")
	pk := os.Getenv("CRE_LOADTEST_PRIVATE_KEY")
	if gw == "" || pk == "" {
		t.Skip("set LOADTEST_GATEWAY_URL and CRE_LOADTEST_PRIVATE_KEY to run the confidential-workflows load test")
	}

	key, err := crypto.HexToECDSA(strings.TrimPrefix(pk, "0x"))
	require.NoError(t, err, "CRE_LOADTEST_PRIVATE_KEY is not a valid hex secp256k1 key")

	owner := os.Getenv("LOADTEST_WORKFLOW_OWNER")
	name := envOr("LOADTEST_WORKFLOW_NAME", "cn_confidential_workflows_load_a")
	id := os.Getenv("LOADTEST_WORKFLOW_ID")
	require.True(t, id != "" || (owner != "" && name != ""),
		"set LOADTEST_WORKFLOW_ID, or both LOADTEST_WORKFLOW_OWNER and LOADTEST_WORKFLOW_NAME")

	start := envIntOr(t, "LOADTEST_RPS_START", 1)
	stepBy := envIntOr(t, "LOADTEST_RPS_STEP", 1)
	nSteps := envIntOr(t, "LOADTEST_RPS_STEPS", 3)
	stepSecs := envIntOr(t, "LOADTEST_STEP_SECONDS", 60)
	maxRPS := envIntOr(t, "LOADTEST_MAX_RPS", 5)
	require.Positive(t, start, "LOADTEST_RPS_START must be > 0")
	require.Positive(t, nSteps, "LOADTEST_RPS_STEPS must be > 0")

	steps := make([]rpsStep, 0, nSteps)
	for i := 0; i < nSteps; i++ {
		rps := start + i*stepBy
		require.LessOrEqualf(t, rps, maxRPS,
			"step %d rps=%d exceeds LOADTEST_MAX_RPS=%d; raise the cap deliberately or lower the ramp", i, rps, maxRPS)
		steps = append(steps, rpsStep{rps: rps, dur: time.Duration(stepSecs) * time.Second})
	}

	return loadConfig{
		gatewayURL:   gw,
		privKey:      key,
		owner:        owner,
		name:         name,
		tag:          os.Getenv("LOADTEST_WORKFLOW_TAG"),
		id:           id,
		input:        envOr("LOADTEST_INPUT", `{"n":1}`),
		callTimeout:  time.Duration(envIntOr(t, "LOADTEST_CALL_TIMEOUT_SECONDS", 120)) * time.Second,
		steps:        steps,
		abortFailPct: 50,
		abortWindow:  envIntOr(t, "LOADTEST_ABORT_WINDOW", 20),
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envIntOr(t *testing.T, k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	require.NoErrorf(t, err, "%s must be an integer", k)
	return n
}
