package ccip

import (
	"context"
	"net/http"
	"testing"
	"time"

	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
)

// WaitForSuiRPCReady waits for SUI RPC to be ready and responsive
// This should be called after SUI containers are started but before test logic
func WaitForSuiRPCReady(t *testing.T, suiChain cldf_sui.Chain, timeout time.Duration) {
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
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Try a simple HTTP request to the SUI RPC endpoint
	resp, err := client.Get(suiURL)
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
		WaitForSuiRPCReady(t, suiChain, timeout)
	}

	t.Logf("✅ All %d SUI chain(s) are ready", len(suiChains))
}

// EnsureSuiRPCHealth performs a comprehensive health check on SUI RPC
func EnsureSuiRPCHealth(t *testing.T, suiChain cldf_sui.Chain) {
	t.Helper()

	// Basic connectivity check
	WaitForSuiRPCReady(t, suiChain, 30*time.Second)

	// Additional SUI-specific health checks can be added here
	// For example, checking if we can get the SUI address
	if suiChain.Signer != nil {
		addr, err := suiChain.Signer.GetAddress()
		if err != nil {
			t.Logf("⚠️  Warning: Could not get SUI signer address: %v", err)
		} else {
			t.Logf("✅ SUI signer address: %s", addr)
		}
	}

	// Check if client is responsive
	if suiChain.Client != nil {
		// We could add more specific SUI client health checks here
		t.Logf("✅ SUI client is available")
	}
}
