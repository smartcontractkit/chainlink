#!/usr/bin/env bash
#
# run.sh — fast local-iteration loop for the CRE node.
#
#   1. stops the local CRE env (if running)
#   2. cross-builds the chainlink node on the HOST for linux/arm64, using the
#      local `replace => ../chainlink-common` in go.mod (no in-Docker compile)
#   3. cross-builds the needed CRE capability plugins (from the local capabilities
#      repo) and bakes them + the node into a minimal self-contained image
#      (ubuntu:24.04 + runtime libs) — no dependency on any prebuilt chainlink image
#   4. generates a topology pointed at that image
#   5. starts the env with observability
#   6. compiles and deploys the workflow
#
# Capability plugins live at /usr/local/bin/<binary_name> in the container
# (DefaultCapabilitiesDir); the env's standard-capability jobs exec them there.
# don-time and vault are in-process (no binary). EVM relayer is in-process (CL_EVM_CMD='').
#
# Prereqs (one-time, see memory: chainlink-host-crossbuild-arm64):
#   ~/clcross/bin/{cc,cxx}                  zig wrappers (-target aarch64-linux-gnu.2.39 -nostdlib++)
#   ~/clcross/dyn/libstdc++.so*             shared GNU libstdc++ for arm64
#   ~/clcross/lib/{libstdc++.a,libsupc++.a} static backfill
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

# Local capabilities repo + which plugins to bake in (dir : output binary_name : extra go flags).
# binary_name must match capability_defaults.toml. don-time/vault are in-process (omitted).
CAPS_REPO="${CAPS_REPO:-$HOME/go/src/github.com/smartcontractkit/capabilities}"
CAP_SPECS=(
  "cron:cron:-tags timetzdata"
  "consensus:consensus:"
  "http_action:http_action:"
  "http_trigger:http_trigger:"
  "chain_capabilities/evm:evm:"
)

# Workflow to compile + deploy after the env is up.
LCRE="${LCRE:-$HOME/go/src/github.com/smartcontractkit/cre-cli/cre}"
WORKFLOW_PROJECT="${WORKFLOW_PROJECT:-$HOME/lcreWorkflow/proj}"
WORKFLOW_DIR="${WORKFLOW_DIR:-my-workflow}"
WORKFLOW_NAME="${WORKFLOW_NAME:-registryreaderdemo}"
WORKFLOW_CONFIG="${WORKFLOW_CONFIG:-$WORKFLOW_PROJECT/$WORKFLOW_DIR/config.staging.json}"

# --- 0. prereq checks -------------------------------------------------------
for f in "$CLCROSS/bin/cc" "$CLCROSS/bin/cxx" "$CLCROSS/dyn/libstdc++.so" \
         "$CLCROSS/lib/libstdc++.a" "$CLCROSS/lib/libsupc++.a"; do
  [ -e "$f" ] || { echo "ERROR: missing prereq $f (see memory: chainlink-host-crossbuild-arm64)"; exit 1; }
done
command -v brotli >/dev/null || { echo "ERROR: brotli not found (brew install brotli)"; exit 1; }
[ -x "$LCRE" ] || { echo "ERROR: lcre not found/executable at $LCRE"; exit 1; }
[ -d "$WORKFLOW_PROJECT/$WORKFLOW_DIR" ] || { echo "ERROR: workflow dir not found: $WORKFLOW_PROJECT/$WORKFLOW_DIR"; exit 1; }
[ -d "$CAPS_REPO" ] || { echo "ERROR: capabilities repo not found: $CAPS_REPO"; exit 1; }
if ! grep -q 'replace github.com/smartcontractkit/chainlink-common => ../chainlink-common' "$REPO_ROOT/go.mod"; then
  echo "WARN: local chainlink-common replace not found in go.mod — building the pinned (published) version."
fi

# --- 1. stop env if running -------------------------------------------------
echo "==> [1/6] Stopping local env (if running)"
go run . env stop || true

# --- 2. host cross-build the node (linux/arm64) -----------------------------
echo "==> [2/6] Building chainlink node on host for linux/arm64"
( cd "$REPO_ROOT" &&
  CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    CC="$CLCROSS/bin/cc" CXX="$CLCROSS/bin/cxx" \
    LIBRARY_PATH="$CLCROSS/dyn" \
    CGO_LDFLAGS="-L$CLCROSS/dyn $CLCROSS/lib/libstdc++.a $CLCROSS/lib/libsupc++.a" \
    go build -o "$BIN" . )
echo "    built $BIN ($(du -h "$BIN" | cut -f1))"

# --- 3. build capability plugins + bake a self-contained image --------------
echo "==> [3/6] Building capability plugins + dev image $DEV_IMAGE"
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
echo "==> [4/6] Generating topology $GEN_TOPO (image=$DEV_IMAGE)"
sed -E \
  -e 's#^([[:space:]]*)docker_ctx = .*#\1image = "'"$DEV_IMAGE"'"#' \
  -e '/^[[:space:]]*docker_file = "core\/chainlink.Dockerfile"/d' \
  "$SRC_TOPO" > "$GEN_TOPO"

# --- 5. start the env with observability ------------------------------------
echo "==> [5/6] Starting local CRE (observability) with $DEV_IMAGE"
CTF_CONFIGS="$GEN_TOPO" go run . env start --with-observability

# --- 6. compile + deploy the workflow ---------------------------------------
echo "==> [6/6] Compiling + deploying workflow '$WORKFLOW_NAME'"
( cd "$WORKFLOW_PROJECT" && "$LCRE" workflow build "$WORKFLOW_DIR" )
WF_ARTIFACT="$(mktemp -d)/$WORKFLOW_NAME.br.b64"
brotli -c -q 5 "$WORKFLOW_PROJECT/$WORKFLOW_DIR/binary.wasm" | base64 | tr -d '\n' > "$WF_ARTIFACT"
go run . env workflow deploy -w "$WF_ARTIFACT" -c "$WORKFLOW_CONFIG" -n "$WORKFLOW_NAME"

echo ""
echo "==> Done: env up on $DEV_IMAGE (node + capability plugins), workflow '$WORKFLOW_NAME' deployed."
echo "    Logs: http://localhost:3000  (Explore -> Loki -> {job=\"ctf\"} |= \"User log\")"
