package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

// stubDON implements handlers.DON. It tracks per-request latency by recording
// the time between when a request is submitted and when SendToNode is called
// with the response.
type stubDON struct {
	mu        sync.Mutex
	wg        sync.WaitGroup
	starts    map[string]time.Time
	latencies []time.Duration
	errors    int
}

func newStubDON() *stubDON {
	return &stubDON{
		starts: make(map[string]time.Time),
	}
}

func (d *stubDON) RegisterRequest(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.starts[id] = time.Now()
	d.wg.Add(1)
}

func (d *stubDON) SendToNode(_ context.Context, _ string, req *jsonrpc.Request[json.RawMessage]) error {
	defer d.wg.Done()
	d.mu.Lock()
	defer d.mu.Unlock()

	elapsed := time.Duration(0)
	if start, ok := d.starts[req.ID]; ok {
		elapsed = time.Since(start)
		delete(d.starts, req.ID)
	}

	// Check whether the handler response indicates an HTTP-level error
	if req.Params != nil {
		var resp gateway_common.OutboundHTTPResponse
		if err := json.Unmarshal(*req.Params, &resp); err == nil && resp.ErrorMessage != "" {
			d.errors++
			return nil
		}
	}

	d.latencies = append(d.latencies, elapsed)
	return nil
}

func (d *stubDON) Wait() {
	d.wg.Wait()
}

func main() {
	n := flag.Int("n", 100, "number of parallel requests")
	destURL := flag.String("url", "http://127.0.0.1:8080", "destination server URL")
	timeoutMs := flag.Int("timeout", 30000, "per-request timeout in milliseconds")
	destPort := flag.Int("dest-port", 8080, "destination server port (added to AllowedPorts)")
	flag.Parse()

	lggr, err := logger.New()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	// Build HTTP client that allows localhost connections
	httpClient, err := network.NewHTTPClient(network.HTTPClientConfig{
		AllowedIPs:     []string{"127.0.0.1"},
		AllowedPorts:   []int{80, 443, *destPort},
		AllowedSchemes: []string{"http", "https"},
	}, lggr)
	if err != nil {
		log.Fatalf("failed to create HTTP client: %v", err)
	}

	// Handler config with generous rate limits
	handlerCfg := v2.ServiceConfig{
		NodeRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100000,
			GlobalBurst:    100000,
			PerSenderRPS:   100000,
			PerSenderBurst: 100000,
		},
	}
	cfgBytes, err := json.Marshal(handlerCfg)
	if err != nil {
		log.Fatalf("failed to marshal handler config: %v", err)
	}

	donConfig := &config.DONConfig{
		DonId: "loadtest-don",
	}
	don := newStubDON()
	lf := limits.Factory{Logger: lggr}

	handler, err := v2.NewGatewayHandler(cfgBytes, donConfig, don, httpClient, lggr, lf)
	if err != nil {
		log.Fatalf("failed to create gateway handler: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := handler.Start(ctx); err != nil {
		log.Fatalf("failed to start handler: %v", err)
	}
	defer func() {
		if err := handler.Close(); err != nil {
			log.Printf("error closing handler: %v", err)
		}
	}()

	log.Printf("starting load test: n=%d url=%s timeout=%dms", *n, *destURL, *timeoutMs)

	// Build all requests up front
	type request struct {
		id   string
		resp *jsonrpc.Response[json.RawMessage]
	}
	requests := make([]request, *n)
	for i := range requests {
		reqID := fmt.Sprintf("%s/loadtest-workflow/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       *destURL,
			TimeoutMs: uint32(*timeoutMs),
		}
		reqBytes, err := json.Marshal(outboundReq)
		if err != nil {
			log.Fatalf("failed to marshal outbound request: %v", err)
		}
		raw := json.RawMessage(reqBytes)
		requests[i] = request{
			id: reqID,
			resp: &jsonrpc.Response[json.RawMessage]{
				ID:     reqID,
				Result: &raw,
			},
		}
	}

	// Fire all requests in parallel
	start := time.Now()
	var submitWg sync.WaitGroup
	for i := range requests {
		don.RegisterRequest(requests[i].id)
		submitWg.Add(1)
		go func(r request) {
			defer submitWg.Done()
			if err := handler.HandleNodeMessage(ctx, r.resp, "loadtest-node"); err != nil {
				don.mu.Lock()
				don.errors++
				don.mu.Unlock()
				log.Printf("HandleNodeMessage error for %s: %v", r.id, err)
				don.wg.Done() // unblock wait since handler won't call SendToNode
			}
		}(requests[i])
	}

	// Wait for all submissions
	submitWg.Wait()

	// Wait for all responses (handler processes async via goroutines)
	don.Wait()
	totalDuration := time.Since(start)

	printStats(don, *n, totalDuration)
}

func printStats(don *stubDON, n int, total time.Duration) {
	don.mu.Lock()
	defer don.mu.Unlock()

	fmt.Println("\n=== Load Test Results ===")
	fmt.Printf("Total requests:   %d\n", n)
	fmt.Printf("Errors:           %d\n", don.errors)
	fmt.Printf("Total duration:   %v\n", total)
	fmt.Printf("Throughput:       %.1f req/s\n", float64(n)/total.Seconds())

	if len(don.latencies) == 0 {
		fmt.Println("No successful requests to report latencies.")
		return
	}

	sort.Slice(don.latencies, func(i, j int) bool {
		return don.latencies[i] < don.latencies[j]
	})

	fmt.Printf("\n--- Latency Distribution ---\n")
	fmt.Printf("Min:  %v\n", don.latencies[0])
	fmt.Printf("P50:  %v\n", percentile(don.latencies, 50))
	fmt.Printf("P90:  %v\n", percentile(don.latencies, 90))
	fmt.Printf("P95:  %v\n", percentile(don.latencies, 95))
	fmt.Printf("P99:  %v\n", percentile(don.latencies, 99))
	fmt.Printf("Max:  %v\n", don.latencies[len(don.latencies)-1])

	var sum time.Duration
	for _, l := range don.latencies {
		sum += l
	}
	fmt.Printf("Mean: %v\n", sum/time.Duration(len(don.latencies)))
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
