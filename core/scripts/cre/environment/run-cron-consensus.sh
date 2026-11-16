#!/usr/bin/env bash
#
# run-cron-consensus.sh — run cron and consensus from the local capabilities repo in a local CRE.
#
#   1. starts the env, building the node and every capability from local source, including
#      both the legacy and the v2 binaries of cron and consensus
#   2. compiles + deploys a workflow that uses cron and consensus (./cronConsensusWorkflow)
#   3. watches every workflow node's executions until Ctrl-C (stops the env on exit)
#
# Which nodes run legacy and which run v2 is chosen beforehand with ./capability-mode.sh
# (legacy, v2, or both: two of each).
#
# Set KEEP_ENV=1 to leave containers running on exit, for reading logs after a
# failure. Stop them yourself afterwards with `go run . env stop`.
#
# One-time setup is NOT done here: the $CLCROSS cross-compile toolchain and the
# managed images (go run . env setup) must already be in place. The node and the
# capabilities build with whatever their go.mod files currently replace.
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

CLCROSS="${CLCROSS:-$HOME/clcross}"
DEV_IMAGE="${DEV_IMAGE:-cre-node:dev}"
TOPO="${TOPO:-configs/workflow-gateway-capabilities-don-cron-consensus.toml}"
CAPS_REPO="${CAPS_REPO:-$HOME/go/src/github.com/smartcontractkit/capabilities}"
WORKFLOW_DIR="${WORKFLOW_DIR:-./cronConsensusWorkflow}"
WORKFLOW_NAME="${WORKFLOW_NAME:-cronconsensusdemo}"
WORKFLOW_NODES="${WORKFLOW_NODES:-workflow-node0 workflow-node1 workflow-node2 workflow-node3}"

# --- 0. lightweight prereq checks (fail fast, no setup) ---------------------
for f in "$CLCROSS/bin/cc" "$CLCROSS/bin/cxx" "$CLCROSS/dyn/libstdc++.so" \
         "$CLCROSS/lib/libstdc++.a" "$CLCROSS/lib/libsupc++.a"; do
  [ -e "$f" ] || { echo "ERROR: missing prereq '$f' — run one-time setup first (see header)"; exit 1; }
done
[ -d "$CAPS_REPO" ] || { echo "ERROR: capabilities repo not found: $CAPS_REPO"; exit 1; }
[ -f "$TOPO" ] || { echo "ERROR: topology not found: $TOPO"; exit 1; }
[ -f "$WORKFLOW_DIR/main.go" ] || { echo "ERROR: workflow not found: $WORKFLOW_DIR/main.go"; exit 1; }

# --- Cleanup handler on Ctrl+C / exit ----------------------------------------
cleanup() {
  trap - INT TERM EXIT
  # shellcheck disable=SC2046 # one argument per job
  kill $(jobs -p) 2>/dev/null || true
  echo ""
  if [ -n "${KEEP_ENV:-}" ]; then
    echo "==> KEEP_ENV set; leaving env up for inspection."
    echo "    node logs:  docker logs workflow-node0"
    echo "    stop later: (cd $SCRIPT_DIR && go run . env stop)"
    return
  fi
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
echo "==> [1/3] Starting local CRE with local source builds"
awk '/# capability-mode$/ { printf "    workflow-node%d: %s\n", n++, (index($0, "[2000-01-01") ? "v2" : "legacy") }' "$TOPO"
export CRE_LOCAL_CC="$CLCROSS/bin/cc"
export CRE_LOCAL_CXX="$CLCROSS/bin/cxx"
export CRE_LOCAL_LIBDIR="$CLCROSS/dyn"

CTF_CONFIGS="$TOPO" go run . env start --with-observability \
  --local-node \
  --local-capabilities all \
  --capabilities-path "$CAPS_REPO" \
  --local-build-platform linux/arm64 \
  --local-node-image "$DEV_IMAGE"

# --- 2. compile + deploy the workflow ---------------------------------------
echo "==> [2/3] Compiling + deploying workflow '$WORKFLOW_NAME'"
go run . env workflow deploy --compile -w "$WORKFLOW_DIR/main.go" -c "$WORKFLOW_DIR/config.json" -n "$WORKFLOW_NAME"

# --- 3. watch workflow executions until killed ------------------------------
echo ""
echo "==> [3/3] Watching '$WORKFLOW_NAME' on $WORKFLOW_NODES (Ctrl-C to stop)"
echo "    (Grafana: http://localhost:3000  |  each cron tick, every 30s, prints below)"
for node in $WORKFLOW_NODES; do
  docker logs -f --since 10s "$node" 2>&1 \
    | grep --line-buffered -aE "Workflow execution finished successfully|Workflow execution failed|Consensus reached" \
    | while IFS= read -r line; do
        ts="$(date '+%H:%M:%S')"
        case "$line" in
          *"execution failed"*) echo "[$ts] $node ❌ ${line}" ;;
          *"Consensus reached"*) echo "[$ts] $node 🤝 ${line}" ;;
          *) echo "[$ts] $node ✅ ${line}" ;;
        esac
      done &
done
wait
