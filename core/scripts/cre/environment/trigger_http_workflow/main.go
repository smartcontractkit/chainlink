// Package main implements a standalone CLI tool that sends JWT-signed HTTP
// trigger requests to a CRE gateway endpoint. Supports both single-shot
// triggering and RPS-based load testing with concurrent workers.
//
// Single trigger:
//
//	go run ./cre/environment/trigger_http_workflow/ \
//	  --gateway-url http://localhost:5002 \
//	  --workflow-name http_simple \
//	  --input '{"customer":"test","size":"large"}'
//
// Load test at 400 RPM for 60 seconds:
//
//	go run ./cre/environment/trigger_http_workflow/ \
//	  --gateway-url http://localhost:5002 \
//	  --workflow-name http_simple \
//	  --workflow-id <id> \
//	  --rps 6.67 --duration 60s --workers 20
package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	gatewayURL := flag.String("gateway-url", "http://localhost:5002", "Gateway URL")
	privateKeyHex := flag.String("private-key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private key (hex, no 0x prefix). Default: Anvil account #0")
	workflowName := flag.String("workflow-name", "", "Workflow name (required)")
	workflowOwner := flag.String("workflow-owner", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address (must be lowercase)")
	workflowID := flag.String("workflow-id", "", "Workflow ID (hex)")
	workflowTag := flag.String("workflow-tag", "some-tag", "Workflow tag")
	input := flag.String("input", `{"customer":"test","size":"large"}`, "JSON input for the workflow trigger")

	// Single-shot mode
	count := flag.Int("count", 1, "Number of sequential trigger requests (single-shot mode)")
	delay := flag.Duration("delay", time.Second, "Delay between requests in single-shot mode")

	// Load test mode
	rps := flag.Float64("rps", 0, "Target requests per second (enables load test mode). 400 RPM = 6.67 RPS")
	duration := flag.Duration("duration", time.Minute, "Load test duration")
	workers := flag.Int("workers", 20, "Number of concurrent worker goroutines for load test")
	reportInterval := flag.Duration("report-interval", 10*time.Second, "How often to print live stats during load test")

	flag.Parse()

	if *workflowName == "" {
		fmt.Fprintf(os.Stderr, "Error: --workflow-name is required\n")
		flag.Usage()
		os.Exit(1)
	}

	signingKey, err := crypto.HexToECDSA(*privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	address := crypto.PubkeyToAddress(signingKey.PublicKey)
	fmt.Printf("Signer address: %s\n", address.Hex())
	fmt.Printf("Gateway URL:    %s\n", *gatewayURL)
	fmt.Printf("Workflow:       name=%s owner=%s tag=%s\n", *workflowName, *workflowOwner, *workflowTag)
	if *workflowID != "" {
		fmt.Printf("Workflow ID:    %s\n", *workflowID)
	}

	cfg := &triggerConfig{
		gatewayURL:    *gatewayURL,
		signingKey:    signingKey,
		workflowName:  *workflowName,
		workflowOwner: *workflowOwner,
		workflowID:    *workflowID,
		workflowTag:   *workflowTag,
		inputJSON:     *input,
	}

	if *rps > 0 {
		// Load test mode
		fmt.Printf("Mode:           LOAD TEST\n")
		fmt.Printf("Target RPS:     %.2f (%.0f RPM)\n", *rps, *rps*60)
		fmt.Printf("Duration:       %s\n", *duration)
		fmt.Printf("Workers:        %d\n", *workers)
		fmt.Println()
		runLoadTest(cfg, *rps, *duration, *workers, *reportInterval)
	} else {
		// Single-shot mode
		fmt.Printf("Input:          %s\n", *input)
		fmt.Println()
		client := &http.Client{Timeout: 2 * time.Minute}
		for i := 0; i < *count; i++ {
			if *count > 1 {
				fmt.Printf("--- Request %d/%d ---\n", i+1, *count)
			}
			_, err := sendTrigger(client, cfg)
			if err != nil {
				log.Printf("Request %d failed: %v", i+1, err)
			}
			if i < *count-1 {
				time.Sleep(*delay)
			}
		}
	}
}

type triggerConfig struct {
	gatewayURL    string
	signingKey    *ecdsa.PrivateKey
	workflowName  string
	workflowOwner string
	workflowID    string
	workflowTag   string
	inputJSON     string
}

type result struct {
	latency    time.Duration
	statusCode int
	err        error
}

// runLoadTest fires requests at the target RPS using a ticker and worker pool.
func runLoadTest(cfg *triggerConfig, targetRPS float64, duration time.Duration, numWorkers int, reportInterval time.Duration) {
	var (
		totalSent    atomic.Int64
		totalSuccess atomic.Int64
		totalFailed  atomic.Int64
		mu           sync.Mutex
		latencies    []time.Duration
	)

	// Results channel and worker pool
	workCh := make(chan int64, numWorkers*2)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 2 * time.Minute}
			for callNum := range workCh {
				start := time.Now()
				statusCode, err := sendTriggerQuiet(client, cfg, callNum)
				elapsed := time.Since(start)

				totalSent.Add(1)
				if err != nil || statusCode != http.StatusOK {
					totalFailed.Add(1)
					if err != nil {
						log.Printf("[#%d] FAIL (%.0fms): %v", callNum, float64(elapsed.Milliseconds()), err)
					} else {
						log.Printf("[#%d] FAIL (%.0fms): status %d", callNum, float64(elapsed.Milliseconds()), statusCode)
					}
				} else {
					totalSuccess.Add(1)
				}

				mu.Lock()
				latencies = append(latencies, elapsed)
				mu.Unlock()
			}
		}()
	}

	// Ticker to emit work at the target RPS
	interval := time.Duration(float64(time.Second) / targetRPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Report ticker
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	startTime := time.Now()
	deadline := startTime.Add(duration)
	var callCounter int64

	fmt.Printf("Load test started at %s\n", startTime.Format(time.RFC3339))
	fmt.Printf("Sending 1 request every %s (target %.2f RPS / %.0f RPM)\n\n", interval, targetRPS, targetRPS*60)

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				goto done
			}
			callCounter++
			workCh <- callCounter

		case <-reportTicker.C:
			elapsed := time.Since(startTime)
			sent := totalSent.Load()
			success := totalSuccess.Load()
			failed := totalFailed.Load()
			actualRPS := float64(sent) / elapsed.Seconds()
			fmt.Printf("[%s] sent=%d success=%d failed=%d actual_rps=%.1f target_rps=%.1f\n",
				elapsed.Round(time.Second), sent, success, failed, actualRPS, targetRPS)
		}
	}

done:
	close(workCh)
	fmt.Printf("\nWaiting for %d in-flight workers to finish...\n", numWorkers)
	wg.Wait()

	// Final report
	totalElapsed := time.Since(startTime)
	sent := totalSent.Load()
	success := totalSuccess.Load()
	failed := totalFailed.Load()
	actualRPS := float64(sent) / totalElapsed.Seconds()

	mu.Lock()
	allLatencies := make([]time.Duration, len(latencies))
	copy(allLatencies, latencies)
	mu.Unlock()

	sort.Slice(allLatencies, func(i, j int) bool { return allLatencies[i] < allLatencies[j] })

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  LOAD TEST RESULTS")
	fmt.Println("========================================")
	fmt.Printf("Duration:       %s\n", totalElapsed.Round(time.Millisecond))
	fmt.Printf("Total sent:     %d\n", sent)
	fmt.Printf("Successful:     %d (%.1f%%)\n", success, pct(success, sent))
	fmt.Printf("Failed:         %d (%.1f%%)\n", failed, pct(failed, sent))
	fmt.Printf("Target RPS:     %.2f (%.0f RPM)\n", targetRPS, targetRPS*60)
	fmt.Printf("Actual RPS:     %.2f (%.0f RPM)\n", actualRPS, actualRPS*60)

	if len(allLatencies) > 0 {
		fmt.Println()
		fmt.Println("Latency (gateway response time):")
		fmt.Printf("  Min:    %s\n", allLatencies[0].Round(time.Millisecond))
		fmt.Printf("  p50:    %s\n", percentile(allLatencies, 50).Round(time.Millisecond))
		fmt.Printf("  p90:    %s\n", percentile(allLatencies, 90).Round(time.Millisecond))
		fmt.Printf("  p95:    %s\n", percentile(allLatencies, 95).Round(time.Millisecond))
		fmt.Printf("  p99:    %s\n", percentile(allLatencies, 99).Round(time.Millisecond))
		fmt.Printf("  Max:    %s\n", allLatencies[len(allLatencies)-1].Round(time.Millisecond))
		fmt.Printf("  Avg:    %s\n", avg(allLatencies).Round(time.Millisecond))
	}
	fmt.Println("========================================")
}

func pct(part, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func avg(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// sendTriggerQuiet sends a trigger and returns status code + error without printing.
func sendTriggerQuiet(client *http.Client, cfg *triggerConfig, callNum int64) (int, error) {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: cfg.workflowOwner,
			WorkflowName:  cfg.workflowName,
			WorkflowTag:   cfg.workflowTag,
			WorkflowID:    cfg.workflowID,
		},
		Input: json.RawMessage(fmt.Sprintf(`{"customer":"load-test-%d","size":"large"}`, callNum)),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		return 0, fmt.Errorf("marshal failed: %w", err)
	}
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      fmt.Sprintf("load-%d-%s", callNum, uuid.New().String()[:8]),
	}

	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		return 0, fmt.Errorf("jwt create failed: %w", err)
	}

	tokenString, err := token.SignedString(cfg.signingKey)
	if err != nil {
		return 0, fmt.Errorf("jwt sign failed: %w", err)
	}
	req.Auth = tokenString

	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", cfg.gatewayURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("create request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/jsonrpc")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("http failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("status %d: %s", resp.StatusCode, truncate(string(body), 200))
	}

	// Check for JSON-RPC error
	var jsonResp struct {
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &jsonResp); err == nil && jsonResp.Error != nil {
		return resp.StatusCode, fmt.Errorf("jsonrpc error %d: %s", jsonResp.Error.Code, jsonResp.Error.Message)
	}

	return resp.StatusCode, nil
}

// sendTrigger sends a trigger and prints the full response (for single-shot mode).
func sendTrigger(client *http.Client, cfg *triggerConfig) (int, error) {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: cfg.workflowOwner,
			WorkflowName:  cfg.workflowName,
			WorkflowTag:   cfg.workflowTag,
			WorkflowID:    cfg.workflowID,
		},
		Input: json.RawMessage(cfg.inputJSON),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal trigger payload: %w", err)
	}
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      fmt.Sprintf("trigger-%s", uuid.New().String()[:8]),
	}

	token, err := utils.CreateRequestJWT(req)
	if err != nil {
		return 0, fmt.Errorf("failed to create JWT: %w", err)
	}

	tokenString, err := token.SignedString(cfg.signingKey)
	if err != nil {
		return 0, fmt.Errorf("failed to sign JWT: %w", err)
	}
	req.Auth = tokenString

	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("Sending request to %s...\n", cfg.gatewayURL)

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", cfg.gatewayURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/jsonrpc")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("Response:\n%s\n", prettyJSON.String())
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
