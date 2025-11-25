#!/bin/bash
set -euo pipefail

# CCIP SUI Test Isolation Script
# Usage: ./run-ccip-sui-test-isolated.sh <test_pattern> <timeout_minutes>
# This script is specifically designed for SUI tests

TEST_PATTERN="${1:-}"
TIMEOUT_MINUTES="${2:-15}"
TEST_ID="ccip-$(date +%s)-$$"
if [[ -z "$TEST_PATTERN" ]]; then
    echo "Usage: $0 <test_pattern> [timeout_minutes]"
    echo "Example: $0 'Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithRateLimit' 20"
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
    
    echo "✅ Cleanup completed"
}

# Set trap for cleanup
trap cleanup EXIT INT TERM

# Pre-test cleanup
echo "🧹 Pre-test cleanup..."
cleanup
sleep 1

# Run the test
echo "▶️  Running test with ${TIMEOUT_MINUTES}m timeout..."
timeout "${TIMEOUT_MINUTES}m" go test ./smoke/ccip \
    -run "$TEST_PATTERN" \
    -test.parallel=1 \
    -count=1 \
    -json \
    -test.timeout=$((TIMEOUT_MINUTES - 2))m

echo "✅ SUI test completed successfully"
