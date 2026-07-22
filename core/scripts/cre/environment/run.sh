#!/usr/bin/env bash
#
# run.sh — fast local-iteration loop for the CRE node.
#
#   1. stops the env  2. host-builds the node (linux/arm64, local chainlink-common)
#   3. builds capability plugins + a self-contained image
#   4. generates a topology using that image  5. starts env + observability
#   6. compiles + deploys the workflow  7. watches results until Ctrl-C
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
BIN="$CLCROSS/chainlink_dev"
SRC_TOPO="${SRC_TOPO:-configs/workflow-gateway-capabilities-don.toml}"
GEN_TOPO="${GEN_TOPO:-configs/run-local-dev.toml}"

CAPS_REPO="${CAPS_REPO:-$HOME/go/src/github.com/smartcontractkit/capabilities}"
CAP_SPECS=(
  "cron:cron:-tags timetzdata"
  "consensus:consensus:"
  "http_action:http_action:"
  "http_trigger:http_trigger:"
  "chain_capabilities/evm:evm:"
  "main:main:"
)

LCRE="${LCRE:-$HOME/go/src/github.com/smartcontractkit/cre-cli/cre}"
WORKFLOW_PROJECT="${WORKFLOW_PROJECT:-$HOME/lcreWorkflow/proj}"
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

# --- 1. stop env if running -------------------------------------------------
echo "==> [1/7] Stopping local env (if running)"
go run . env stop || true
# force-remove leftovers env stop can miss when the state file is lost (orphaned
# anvil/observability containers aren't all ctf-labeled) so the next start doesn't
# hit a container-name conflict. Keeps the unrelated postgres-chip-config.
orphans="$(
  { docker ps -a --format '{{.ID}} {{.Names}}' \
      | grep -aE 'anvil-|chip-router|chip-ingress|jd-|jd-db|workflow-node|capabilities-node|bootstrap-gateway|-ns-postgresql|compose-|postgres_exporter|cadvisor|promtail' \
      | grep -av 'postgres-chip-config' | awk '{print $1}'
    docker ps -aq --filter label=framework=ctf
  } | sort -u | grep -a . || true
)"
[ -n "$orphans" ] && { echo "    removing $(echo "$orphans" | wc -l | tr -d ' ') orphan container(s)"; docker rm -f $orphans >/dev/null 2>&1 || true; }

# --- 2. host cross-build the node (linux/arm64) -----------------------------
echo "==> [2/7] Building chainlink node on host for linux/arm64"
( cd "$REPO_ROOT" &&
  CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    CC="$CLCROSS/bin/cc" CXX="$CLCROSS/bin/cxx" \
    LIBRARY_PATH="$CLCROSS/dyn" \
    CGO_LDFLAGS="-L$CLCROSS/dyn $CLCROSS/lib/libstdc++.a $CLCROSS/lib/libsupc++.a" \
    go build -o "$BIN" . )
echo "    built $BIN ($(du -h "$BIN" | cut -f1))"

# --- 3. build capability plugins + bake a self-contained image --------------
echo "==> [3/7] Building capability plugins + dev image $DEV_IMAGE"
CTX="$(mktemp -d)"
cp "$BIN" "$CTX/chainlink"
mkdir -p "$CTX/caps"
for spec in "${CAP_SPECS[@]}"; do
  IFS=: read -r capdir capname capflags <<< "$spec"
  echo "    - capability '$capname' (from $capdir)"
  ( cd "$CAPS_REPO/$capdir" &&
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $capflags -o "$CTX/caps/$capname" . )
done
cat > "$CTX/Dockerfile" <<'DOCKERFILE'
FROM ubuntu:24.04
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl libstdc++6 \
 && rm -rf /var/lib/apt/lists/*
RUN useradd --uid 14933 --create-home chainlink
COPY chainlink /usr/local/bin/chainlink
COPY caps/ /usr/local/bin/
RUN chmod -R 0755 /usr/local/bin
USER chainlink
WORKDIR /home/chainlink
EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["local", "node"]
DOCKERFILE
docker build --platform linux/arm64 -t "$DEV_IMAGE" "$CTX"
rm -rf "$CTX"

# --- 4. generate a topology that uses the dev image (no source build) -------
echo "==> [4/7] Generating topology $GEN_TOPO (image=$DEV_IMAGE)"
sed -E \
  -e 's#^([[:space:]]*)docker_ctx = .*#\1image = "'"$DEV_IMAGE"'"#' \
  -e '/^[[:space:]]*docker_file = "core\/chainlink.Dockerfile"/d' \
  "$SRC_TOPO" > "$GEN_TOPO"

# --- 5. start the env with observability ------------------------------------
echo "==> [5/7] Starting local CRE (observability) with $DEV_IMAGE"
CTF_CONFIGS="$GEN_TOPO" go run . env start --with-observability

# --- 6. compile + deploy the workflow ---------------------------------------
echo "==> [6/7] Compiling + deploying workflow '$WORKFLOW_NAME'"
( cd "$WORKFLOW_PROJECT" && "$LCRE" workflow build "$WORKFLOW_DIR" )
WF_ARTIFACT="$(mktemp -d)/$WORKFLOW_NAME.br.b64"
brotli -c -q 5 "$WORKFLOW_PROJECT/$WORKFLOW_DIR/binary.wasm" | base64 | tr -d '\n' > "$WF_ARTIFACT"
go run . env workflow deploy -w "$WF_ARTIFACT" -c "$WORKFLOW_CONFIG" -n "$WORKFLOW_NAME"

# --- 7. watch workflow results until killed ---------------------------------
echo ""
echo "==> [7/7] Watching '$WORKFLOW_NAME' results on $WATCH_NODE (Ctrl-C to stop)"
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
