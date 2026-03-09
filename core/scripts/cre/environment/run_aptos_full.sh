#!/usr/bin/env bash
# Start CRE with Aptos topology using remote plugin refs and the normal Chainlink Dockerfile.
# Plugin versions are resolved via plugins.private.yaml/plugins.public.yaml and installed at image-build time.
#
# Prerequisites: Docker, Foundry (anvil). Remote branches must be pushed.
# Run from: core/scripts/cre/environment.
#
# Optional env vars (defaults shown).
#   CHAINLINK_APTOS_BRANCH=aptos-service
#   CAPABILITIES_BRANCH=feature/aptos-service-tmp-2
#   # Optional fallback if aptos-service head is temporarily broken:
#   # CAPABILITIES_BRANCH=fcb512c64aa9
#   CHAINLINK_APTOS_IMAGE=chainlink:aptos-remote

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Chainlink repo root (core/scripts/cre/environment -> 4 levels up)
CHAINLINK_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
BRANCH_APTOS="${CHAINLINK_APTOS_BRANCH:-feature/aptos-cre-tmp-2}"
BRANCH_CAPABILITIES="${CAPABILITIES_BRANCH:-feature/aptos-service-tmp-3}"
BRANCH_DATA_STREAMS="${CHAINLINK_DATA_STREAMS_BRANCH:-master}"
BRANCH_SOLANA="${CHAINLINK_SOLANA_BRANCH:-develop}"
IMAGE_TAG="${CHAINLINK_APTOS_IMAGE:-chainlink:aptos-remote}"

resolve_module_ref() {
  local module="$1"
  local ref="$2"
  GOPRIVATE=github.com/smartcontractkit/* go list -m -f '{{.Version}}' "${module}@${ref}"
}

RESOLVED_APTOS_REF="$(resolve_module_ref github.com/smartcontractkit/chainlink-aptos "$BRANCH_APTOS")"
RESOLVED_CAP_APTOS_REF="$(resolve_module_ref github.com/smartcontractkit/capabilities/chain_capabilities/aptos "$BRANCH_CAPABILITIES")"
RESOLVED_CAP_CONSENSUS_REF="$(resolve_module_ref github.com/smartcontractkit/capabilities/consensus "$BRANCH_CAPABILITIES")"
RESOLVED_DATA_STREAMS_REF="$(resolve_module_ref github.com/smartcontractkit/chainlink-data-streams "$BRANCH_DATA_STREAMS")"
RESOLVED_SOLANA_REF="$(resolve_module_ref github.com/smartcontractkit/chainlink-solana "$BRANCH_SOLANA")"

if [ -z "${CRE_APTOS_CONTRACTS_PATH:-}" ]; then
  DEFAULT_APTOS_CONTRACTS_PATH="$(cd "$CHAINLINK_ROOT/.." && pwd)/chainlink-aptos/contracts"
  if [ -d "$DEFAULT_APTOS_CONTRACTS_PATH" ]; then
    export CRE_APTOS_CONTRACTS_PATH="$DEFAULT_APTOS_CONTRACTS_PATH"
    echo "Using local Aptos contracts path: ${CRE_APTOS_CONTRACTS_PATH}"
  fi
fi

echo "Using remote refs: chainlink-aptos@${BRANCH_APTOS}, capabilities@${BRANCH_CAPABILITIES}"
echo "Resolved chainlink-aptos module ref: ${RESOLVED_APTOS_REF}"
echo "Resolved chainlink-data-streams module ref: ${RESOLVED_DATA_STREAMS_REF}"
echo "Resolved chainlink-solana module ref: ${RESOLVED_SOLANA_REF}"
echo "Resolved Aptos capabilities module ref: ${RESOLVED_CAP_APTOS_REF}"
echo "Resolved consensus capabilities module ref: ${RESOLVED_CAP_CONSENSUS_REF}"

cd "$CHAINLINK_ROOT"

# Keep remote-branch workflow: point private Aptos + consensus capability plugins
# to CAPABILITIES_BRANCH so report-generation algo changes are picked up together.
PRIVATE_PLUGINS_FILE="$CHAINLINK_ROOT/plugins/plugins.private.yaml"
PUBLIC_PLUGINS_FILE="$CHAINLINK_ROOT/plugins/plugins.public.yaml"
for capability in aptos consensus; do
  capability_ref="$RESOLVED_CAP_APTOS_REF"
  if [ "$capability" = "consensus" ]; then
    capability_ref="$RESOLVED_CAP_CONSENSUS_REF"
  fi
  TMP_PRIVATE_PLUGINS="$(mktemp)"
  awk -v branch="$capability_ref" -v capability="$capability" '
    $0 ~ "^  " capability ":$" { in_capability = 1 }
    in_capability && /gitRef:/ {
      $0 = "      gitRef: \"" branch "\""
      in_capability = 0
    }
    { print }
  ' "$PRIVATE_PLUGINS_FILE" > "$TMP_PRIVATE_PLUGINS"
  mv "$TMP_PRIVATE_PLUGINS" "$PRIVATE_PLUGINS_FILE"
done
echo "Using remote capabilities refs: aptos@${RESOLVED_CAP_APTOS_REF}, consensus@${RESOLVED_CAP_CONSENSUS_REF}"

# Keep Aptos relayer plugin aligned with CHAINLINK_APTOS_BRANCH so new Aptos gRPC methods
# are available to capabilities (e.g. AccountTransactions).
TMP_PUBLIC_PLUGINS="$(mktemp)"
awk -v aptos_ref="$RESOLVED_APTOS_REF" '
  $0 ~ "^  aptos:$" { in_aptos = 1 }
  in_aptos && /gitRef:/ {
    $0 = "      gitRef: \"" aptos_ref "\""
    in_aptos = 0
  }
  { print }
' "$PUBLIC_PLUGINS_FILE" > "$TMP_PUBLIC_PLUGINS"
mv "$TMP_PUBLIC_PLUGINS" "$PUBLIC_PLUGINS_FILE"
echo "Using Aptos relayer plugin ref: chainlink-aptos@${RESOLVED_APTOS_REF}"
echo "Using Aptos capabilities plugin ref: capabilities/chain_capabilities/aptos@${RESOLVED_CAP_APTOS_REF}"

# Pre-pull DB images used by CRE/JD stacks. After aggressive local prune these may
# be missing and cause env start to fail with "No such image".
for image in postgres:16 postgres:12.0; do
  if ! docker image inspect "$image" >/dev/null 2>&1; then
    echo "Pulling required image: $image"
    docker pull "$image"
  fi
done

# Build node image with normal Dockerfile (context = chainlink repo only).
# CL_INSTALL_PRIVATE_PLUGINS=true required for cron (cron-trigger). Capabilities repo is private, so GIT_AUTH_TOKEN is required.
# Best-effort fallback so local runs work even when token is not exported.
if [ -z "${GIT_AUTH_TOKEN:-}" ] && command -v gh >/dev/null 2>&1; then
  GIT_AUTH_TOKEN="$(gh auth token 2>/dev/null | tr -d '\n' || true)"
  if [ -n "${GIT_AUTH_TOKEN:-}" ]; then
    export GIT_AUTH_TOKEN
    echo "Using GIT_AUTH_TOKEN from gh auth token"
  fi
fi

if [ -z "${GIT_AUTH_TOKEN:-}" ]; then
  echo "ERROR: GIT_AUTH_TOKEN must be set to build with private plugins (capabilities/cron). Export a GitHub token with access to the capabilities repo."
  exit 1
fi
echo "Building Chainlink node image (core/chainlink.Dockerfile)..."
docker build -t "$IMAGE_TAG" -f core/chainlink.Dockerfile \
  --secret id=GIT_AUTH_TOKEN,env=GIT_AUTH_TOKEN \
  --build-arg CL_IS_PROD_BUILD=false \
  --build-arg CL_INSTALL_PRIVATE_PLUGINS=true \
  --build-arg CL_CAPABILITIES_APTOS_REF="${RESOLVED_CAP_APTOS_REF}" \
  .

cd "$SCRIPT_DIR"

export CTF_CONFIGS=configs/workflow-gateway-don-aptos.toml

STATE_FILE="$SCRIPT_DIR/state/local_cre.toml"
READINESS_TIMEOUT=1200
POLL_INTERVAL=20

echo "Starting CRE environment (Aptos topology) with plugins image $IMAGE_TAG (timeout ${READINESS_TIMEOUT}s)..."
go run . env start -p "$IMAGE_TAG" &
ENV_PID=$!

echo "Waiting for state to be cleared by current run..."
for _ in $(seq 1 30); do
  if [ ! -f "$STATE_FILE" ]; then
    break
  fi
  sleep 2
done

elapsed=0
while [ $elapsed -lt $READINESS_TIMEOUT ]; do
  sleep $POLL_INTERVAL
  elapsed=$((elapsed + POLL_INTERVAL))
  if [ -f "$STATE_FILE" ]; then
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q workflow-node; then
      CAPABILITY_WAIT=${CRE_CAPABILITY_WAIT_SECONDS:-10}
      echo "Waiting ${CAPABILITY_WAIT}s for workflow nodes to sync capabilities (cron, Aptos)..."
      sleep $CAPABILITY_WAIT
      echo "Environment ready (state file and workflow nodes present)."
      echo ""
      echo "Run the Aptos suite test from system-tests/tests:"
      echo "  go test -timeout 15m -run '^Test_CRE_V2_Aptos_Suite\$' ./smoke/cre/"
      echo ""
      echo "To stop CRE when done: go run . env stop -a"
      exit 0
    fi
  fi
  echo "Waiting for environment... ${elapsed}s"
done

echo "Timeout: state file did not appear after ${READINESS_TIMEOUT}s."
echo "You can run manually: CTF_CONFIGS=configs/workflow-gateway-don-aptos.toml go run . env start -p $IMAGE_TAG"
kill $ENV_PID 2>/dev/null || true
exit 1
