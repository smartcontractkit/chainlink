#!/bin/bash
set -euo pipefail

# CCIP Test Isolation Script
# Usage: ./run-ccip-test-isolated.sh <test_pattern> <timeout_minutes>

TEST_PATTERN="${1:-}"
TIMEOUT_MINUTES="${2:-15}"
TEST_ID="ccip-$(date +%s)-$$"
PORT_BASE=$(shuf -i 30000-40000 -n 1)

if [[ -z "$TEST_PATTERN" ]]; then
    echo "Usage: $0 <test_pattern> [timeout_minutes]"
    echo "Example: $0 'Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithRateLimit' 20"
    exit 1
fi

echo "🚀 Starting isolated CCIP test: $TEST_PATTERN"
echo "📊 Test ID: $TEST_ID"
echo "🔌 Port base: $PORT_BASE"

# Export test environment
export CCIP_TEST_ID="$TEST_ID"
export CCIP_TEST_PORT_BASE="$PORT_BASE"

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up test $TEST_ID..."
    
    # Kill processes by test ID
    pkill -f "$TEST_ID" 2>/dev/null || true
    
    # Kill processes in our port range
    for port in $(seq $PORT_BASE $((PORT_BASE + 200))); do
        lsof -ti:$port 2>/dev/null | xargs -r kill -9 2>/dev/null || true
    done
    
    echo "✅ Cleanup completed"
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

# Pre-test cleanup
echo "🧹 Pre-test cleanup..."
cleanup
sleep 2

# Run the test
echo "▶️  Running test with ${TIMEOUT_MINUTES}m timeout..."
timeout "${TIMEOUT_MINUTES}m" go test ./smoke/ccip \
    -run "$TEST_PATTERN" \
    -test.parallel=1 \
    -count=1 \
    -json \
    -test.timeout=$((TIMEOUT_MINUTES - 2))m

echo "✅ Test completed successfully"
