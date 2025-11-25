#!/bin/bash
set -euo pipefail

# CCIP SUI Test Isolation Script
# Usage: ./run-ccip-sui-test-isolated.sh <test_pattern> <timeout_minutes>
# This script is specifically designed for SUI tests and ensures SUI RPC is available

TEST_PATTERN="${1:-}"
TIMEOUT_MINUTES="${2:-15}"
TEST_ID="ccip-$(date +%s)-$$"
if [[ -z "$TEST_PATTERN" ]]; then
    echo "Usage: $0 <test_pattern> [timeout_minutes]"
    echo "Example: $0 'Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithRateLimit' 20"
    echo "Note: This script is specifically for SUI tests and requires SUI RPC to be available"
    exit 1
fi

echo "🚀 Starting isolated CCIP SUI test: $TEST_PATTERN"
echo "📊 Test ID: $TEST_ID"

# Export test environment
export CCIP_TEST_ID="$TEST_ID"

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up test $TEST_ID..."
    
    # Kill processes by test ID
    pkill -f "$TEST_ID" 2>/dev/null || true
    
    echo "🧹 SUI containers left running for test compatibility..."

    # Health check SUI RPC connection
    if command -v curl >/dev/null 2>&1; then
        echo "🔍 Checking SUI RPC health on port 9000..."
        if curl -s -f http://127.0.0.1:9000 >/dev/null 2>&1; then
            echo "✅ SUI RPC is responding on port 9000"
        else
            echo "⚠️  SUI RPC not responding on port 9000 - this may cause test failures"
        fi
    fi
    
    echo "✅ Cleanup completed"
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

# Pre-test cleanup
echo "🧹 Pre-test cleanup..."
cleanup

# Wait for SUI RPC to be ready (required for all SUI tests)
echo "⏳ Waiting for SUI RPC to be ready..."
SUI_RPC_READY=false
for i in {1..30}; do
    if command -v curl >/dev/null 2>&1 && curl -s -f http://127.0.0.1:9000 >/dev/null 2>&1; then
        echo "✅ SUI RPC is ready on port 9000"
        SUI_RPC_READY=true
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ SUI RPC not ready after 30 seconds"
        echo "💡 This will cause SUI tests to fail"
        echo "🔧 Try: docker ps | grep sui"
        echo "🔧 Try: docker logs <sui_container_id>"
        exit 1
    fi
    sleep 1
done

if [ "$SUI_RPC_READY" = false ]; then
    echo "❌ SUI RPC check failed - cannot proceed with SUI tests"
    exit 1
fi

sleep 2

# Run the test
echo "▶️  Running test with ${TIMEOUT_MINUTES}m timeout..."
timeout "${TIMEOUT_MINUTES}m" go test ./smoke/ccip \
    -run "$TEST_PATTERN" \
    -test.parallel=1 \
    -count=1 \
    -json \
    -test.timeout=$((TIMEOUT_MINUTES - 2))m

echo "✅ SUI test completed successfully"
