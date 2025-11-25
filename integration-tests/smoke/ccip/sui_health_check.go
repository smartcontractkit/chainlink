package ccip

import (
	"context"
	"net/http"
	"testing"
	"time"

	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
)

// waitForSuiRPCReady waits for SUI RPC to be ready and responsive
// This should be called after SUI containers are started but before test logic
func waitForSuiRPCReady(t *testing.T, suiChain cldf_sui.Chain, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	suiURL := suiChain.URL
	t.Logf("⏳ Waiting for SUI RPC to be ready at %s (timeout: %v)", suiURL, timeout)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("❌ SUI RPC not ready after %v timeout at %s", timeout, suiURL)
		case <-ticker.C:
			if isSuiRPCReady(suiURL) {
				elapsed := time.Since(startTime)
				t.Logf("✅ SUI RPC is ready at %s (took %v)", suiURL, elapsed.Round(100*time.Millisecond))
				return
			}
			t.Logf("🔄 SUI RPC not ready yet at %s (elapsed: %v)", suiURL, time.Since(startTime).Round(time.Second))
		}
	}
}

// isSuiRPCReady checks if SUI RPC is responding to basic HTTP requests
func isSuiRPCReady(suiURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a GET request with proper context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, suiURL, nil)
	if err != nil {
		return false
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// SUI RPC should respond with some HTTP status (even if it's an error for GET requests)
	// The important thing is that the port is open and responding
	return resp.StatusCode > 0
}

// WaitForAllSuiChainsReady waits for all SUI chains in the environment to be ready
func WaitForAllSuiChainsReady(t *testing.T, suiChains map[uint64]cldf_sui.Chain, timeout time.Duration) {
	t.Helper()

	if len(suiChains) == 0 {
		t.Log("ℹ️  No SUI chains to check")
		return
	}

	t.Logf("🔍 Checking %d SUI chain(s) for readiness", len(suiChains))

	for chainSelector, suiChain := range suiChains {
		t.Logf("🔍 Checking SUI chain %d", chainSelector)
		waitForSuiRPCReady(t, suiChain, timeout)
	}

	t.Logf("✅ All %d SUI chain(s) are ready", len(suiChains))
}
