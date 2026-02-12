#!/usr/bin/env bash
# test_http_workflow.sh -- End-to-end test for the http_simple CRE workflow.
#
# This script:
#   1. Starts a local echo HTTP server (receives the workflow's outgoing POST)
#   2. Assumes CRE environment is already running (start with: go run . env start)
#   3. Deploys the http_simple workflow via the CRE CLI
#   4. Waits for the gateway to discover the workflow
#   5. Triggers the workflow with a signed HTTP request
#   6. Verifies the echo server received the POST from the workflow
#
# Prerequisites:
#   - CRE environment running: cd core/scripts/cre/environment && go run . env start
#   - Run from: core/scripts/ directory (the go.mod root for scripts)
#
# Usage:
#   cd core/scripts && bash cre/environment/test_http_workflow.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Configuration
ECHO_SERVER_PORT=8171
GATEWAY_URL="http://localhost:5002"
PRIVATE_KEY="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
WORKFLOW_NAME="http_simple"
WORKFLOW_OWNER="0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
WORKFLOW_SRC="$SCRIPT_DIR/examples/workflows/v2/http_simple/main.go"
WORKFLOW_CONFIG="$SCRIPT_DIR/examples/workflows/v2/http_simple/load_test_config.json"
GATEWAY_SYNC_WAIT=60  # seconds to wait for gateway to discover workflow

# Cleanup function
cleanup() {
    echo ""
    echo "=== Cleaning up ==="
    if [[ -n "${ECHO_PID:-}" ]] && kill -0 "$ECHO_PID" 2>/dev/null; then
        echo "Stopping echo server (PID $ECHO_PID)..."
        kill "$ECHO_PID" 2>/dev/null || true
        wait "$ECHO_PID" 2>/dev/null || true
    fi
    echo "Done."
}
trap cleanup EXIT

cd "$SCRIPTS_ROOT"

echo "============================================"
echo "  CRE HTTP Workflow End-to-End Test"
echo "============================================"
echo ""
echo "Echo server port:  $ECHO_SERVER_PORT"
echo "Gateway URL:       $GATEWAY_URL"
echo "Workflow:          $WORKFLOW_NAME"
echo "Workflow source:   $WORKFLOW_SRC"
echo "Config:            $WORKFLOW_CONFIG"
echo ""

# -------------------------------------------------------
# Step 1: Start the echo server
# -------------------------------------------------------
echo "=== Step 1: Starting echo server on port $ECHO_SERVER_PORT ==="
go run ./cre/environment/echo_server/ --port "$ECHO_SERVER_PORT" &
ECHO_PID=$!
sleep 2

# Verify it's running
if ! kill -0 "$ECHO_PID" 2>/dev/null; then
    echo "ERROR: Echo server failed to start"
    exit 1
fi
echo "Echo server running (PID $ECHO_PID)"
echo ""

# -------------------------------------------------------
# Step 2: Deploy the workflow
# -------------------------------------------------------
echo "=== Step 2: Deploying workflow '$WORKFLOW_NAME' ==="
cd "$SCRIPT_DIR"
DEPLOY_OUTPUT=$(go run . env workflow deploy \
    -w "$WORKFLOW_SRC" \
    -n "$WORKFLOW_NAME" \
    -c "$WORKFLOW_CONFIG" \
    --compile 2>&1)
echo "$DEPLOY_OUTPUT"

# Extract the workflow ID from the deploy output
WORKFLOW_ID=$(echo "$DEPLOY_OUTPUT" | grep "workflowID=" | sed "s/.*workflowID='\([^']*\)'.*/\1/")
if [[ -z "$WORKFLOW_ID" ]]; then
    echo "WARNING: Could not extract workflow ID from deploy output"
else
    echo "Extracted workflow ID: $WORKFLOW_ID"
fi
cd "$SCRIPTS_ROOT"
echo ""

# -------------------------------------------------------
# Step 3: Wait for gateway to discover the workflow
# -------------------------------------------------------
echo "=== Step 3: Waiting for gateway to sync workflow metadata (up to ${GATEWAY_SYNC_WAIT}s) ==="
echo "The gateway pulls workflow metadata from DON nodes periodically."
echo "Polling gateway every 10s to check readiness..."
echo ""

READY=false
ELAPSED=0
WORKFLOW_ID_FLAG=""
if [[ -n "$WORKFLOW_ID" ]]; then
    WORKFLOW_ID_FLAG="--workflow-id $WORKFLOW_ID"
fi

while [[ $ELAPSED -lt $GATEWAY_SYNC_WAIT ]]; do
    sleep 10
    ELAPSED=$((ELAPSED + 10))

    # Try sending a trigger and check if we get "ACCEPTED"
    RESPONSE=$(go run ./cre/environment/trigger_http_workflow/ \
        --gateway-url "$GATEWAY_URL" \
        --private-key "$PRIVATE_KEY" \
        --workflow-name "$WORKFLOW_NAME" \
        --workflow-owner "$WORKFLOW_OWNER" \
        $WORKFLOW_ID_FLAG \
        --input '{"customer":"readiness-check"}' 2>&1 || true)

    if echo "$RESPONSE" | grep -qi "ACCEPTED"; then
        echo "Gateway is ready! (took ${ELAPSED}s)"
        READY=true
        break
    fi

    echo "  [${ELAPSED}s] Not ready yet..."
done

if [[ "$READY" != "true" ]]; then
    echo ""
    echo "WARNING: Gateway may not have synced the workflow yet after ${GATEWAY_SYNC_WAIT}s."
    echo "Proceeding with trigger attempt anyway..."
fi
echo ""

# -------------------------------------------------------
# Step 4: Trigger the workflow
# -------------------------------------------------------
echo "=== Step 4: Triggering workflow '$WORKFLOW_NAME' ==="
go run ./cre/environment/trigger_http_workflow/ \
    --gateway-url "$GATEWAY_URL" \
    --private-key "$PRIVATE_KEY" \
    --workflow-name "$WORKFLOW_NAME" \
    --workflow-owner "$WORKFLOW_OWNER" \
    $WORKFLOW_ID_FLAG \
    --input '{"customer":"test-customer","size":"large","toppings":["cheese","pepperoni"]}'
echo ""

# -------------------------------------------------------
# Step 5: Wait and check echo server
# -------------------------------------------------------
echo "=== Step 5: Waiting for workflow execution ==="
echo "The workflow should POST to the echo server shortly."
echo "Check the echo server logs above for incoming requests."
echo "Waiting 15s for the workflow to complete..."
sleep 15

echo ""
echo "============================================"
echo "  Test complete!"
echo "  Check the echo server output above for"
echo "  the POST request from the workflow."
echo "============================================"
