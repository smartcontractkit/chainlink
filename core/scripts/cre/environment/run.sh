#!/usr/bin/env bash
#
# run.sh — fast local-iteration loop for the CRE node.
#
#   1. starts env + observability (builds node + capabilities on-the-fly)
#   2. compiles + deploys the workflow
#   3. watches results until Ctrl-C (stops env upon exit or Ctrl-C)
#
# This is the iteration loop only. One-time setup is NOT done here — bootstrap
# these once (see memory: chainlink-host-crossbuild-arm64) before first use:
#   - ~/clcross/{bin/{cc,cxx},dyn/libstdc++.so*,lib/{libstdc++.a,libsupc++.a}}
#   - lcre CLI built at $LCRE
#   - managed images present (go run . env setup)
#   - any go.mod replaces you want (e.g. local chainlink-common) — run.sh does
#     NOT add/force them; it builds whatever your go.mod currently points at.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

CLCROSS="${CLCROSS:-$HOME/clcross}"
DEV_IMAGE="${DEV_IMAGE:-cre-node:dev}"
SRC_TOPO="${SRC_TOPO:-configs/workflow-gateway-capabilities-don.toml}"

CAPS_REPO="${CAPS_REPO:-$HOME/go/src/github.com/smartcontractkit/capabilities}"

LCRE="${LCRE:-$HOME/go/src/github.com/smartcontractkit/cre-cli/cre}"
WORKFLOW_PROJECT="${WORKFLOW_PROJECT:-./lcreWorkflow/proj}"
WORKFLOW_DIR="${WORKFLOW_DIR:-my-workflow}"
WORKFLOW_NAME="${WORKFLOW_NAME:-registryreaderdemo}"
WORKFLOW_CONFIG="${WORKFLOW_CONFIG:-$WORKFLOW_PROJECT/$WORKFLOW_DIR/config.staging.json}"
WATCH_NODE="${WATCH_NODE:-workflow-node0}"

# --- 0. lightweight prereq checks (fail fast, no setup) ---------------------
for f in "$CLCROSS/bin/cc" "$CLCROSS/bin/cxx" "$CLCROSS/dyn/libstdc++.so" \
         "$CLCROSS/lib/libstdc++.a" "$CLCROSS/lib/libsupc++.a" "$LCRE"; do
  [ -e "$f" ] || { echo "ERROR: missing prereq '$f' — run one-time setup first (see run.sh header)"; exit 1; }
done
command -v brotli >/dev/null || { echo "ERROR: brotli not found (brew install brotli)"; exit 1; }
[ -d "$CAPS_REPO" ] || { echo "ERROR: capabilities repo not found: $CAPS_REPO"; exit 1; }
[ -d "$WORKFLOW_PROJECT/$WORKFLOW_DIR" ] || { echo "ERROR: workflow dir not found: $WORKFLOW_PROJECT/$WORKFLOW_DIR"; exit 1; }

# --- Cleanup handler on Ctrl+C / exit ----------------------------------------
cleanup() {
  trap - INT TERM EXIT
  echo ""
  echo "==> Stopping local env..."
  go run . env stop || true
  
  # force-remove leftovers env stop can miss when the state file is lost
  orphans="$(
    { docker ps -a --format '{{.ID}} {{.Names}}' \
        | grep -aE 'anvil-|chip-router|chip-ingress|jd-|jd-db|workflow-node|capabilities-node|bootstrap-gateway|-ns-postgresql|compose-|postgres_exporter|cadvisor|promtail' \
        | grep -av 'postgres-chip-config' | awk '{print $1}'
      docker ps -aq --filter label=framework=ctf
    } | sort -u | grep -a . || true
  )"
  [ -n "$orphans" ] && { echo "    removing $(echo "$orphans" | wc -l | tr -d ' ') orphan container(s)"; docker rm -f $orphans >/dev/null 2>&1 || true; }
}

trap cleanup INT TERM EXIT

# --- 1. start the env with local builds and observability -------------------
echo "==> [1/3] Starting local CRE (observability) with local source builds"
export CRE_LOCAL_CC="$CLCROSS/bin/cc"
export CRE_LOCAL_CXX="$CLCROSS/bin/cxx"
export CRE_LOCAL_LIBDIR="$CLCROSS/dyn"

CTF_CONFIGS="$SRC_TOPO" go run . env start --with-observability \
  --local-node \
  --local-capabilities all \
  --capabilities-path "$CAPS_REPO" \
  --local-build-platform linux/arm64 \
  --local-node-image "$DEV_IMAGE"

# --- 2. compile + deploy the workflow ---------------------------------------
echo "==> [2/3] Compiling + deploying workflow '$WORKFLOW_NAME'"
( cd "$WORKFLOW_PROJECT" && "$LCRE" workflow build "$WORKFLOW_DIR" )
WF_ARTIFACT="$(mktemp -d)/$WORKFLOW_NAME.br.b64"
brotli -c -q 5 "$WORKFLOW_PROJECT/$WORKFLOW_DIR/binary.wasm" | base64 | tr -d '\n' > "$WF_ARTIFACT"
go run . env workflow deploy -w "$WF_ARTIFACT" -c "$WORKFLOW_CONFIG" -n "$WORKFLOW_NAME"

# --- 3. watch workflow results until killed ---------------------------------
echo ""
echo "==> [3/3] Watching '$WORKFLOW_NAME' results on $WATCH_NODE (Ctrl-C to stop)"
echo "    (Grafana: http://localhost:3000  |  each cron tick prints below)"
docker logs -f --since 10s "$WATCH_NODE" 2>&1 \
  | grep --line-buffered -aE "Value returned from workflow:|Error message returned from workflow" \
  | while IFS= read -r line; do
      ts="$(date '+%H:%M:%S')"
      if [[ "$line" == *"Value returned from workflow:"* ]]; then
        echo "[$ts] ✅ SUCCESS — ${line#*Value returned from workflow: }"
      else
        echo "[$ts] ❌ FAILURE — ${line#*Error message returned from workflow }"
      fi
    done
