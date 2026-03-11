#!/usr/bin/env bash
# Quick Aptos dev helper: push local changes, update go.mod, start/stop CRE.
#
# Run from: core/scripts/cre/environment
#
#   ./aptos_dev.sh push       # commit+push capabilities & chainlink-aptos, then go get + tidy
#   ./aptos_dev.sh startenv   # start CRE with Aptos topology
#   ./aptos_dev.sh stopenv    # stop CRE

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHAINLINK_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

CAPABILITIES_DIR="$HOME/capabilities"
CHAINLINK_APTOS_DIR="$HOME/chainlink-aptos"

CAPABILITIES_BRANCH="${CAPABILITIES_BRANCH:-aptos-write-report}"
CHAINLINK_APTOS_BRANCH="${CHAINLINK_APTOS_BRANCH:-aptos-service}"

APTOS_PUSHED=0

push_if_dirty() {
  local dir="$1"
  local branch="$2"
  local name
  name="$(basename "$dir")"

  cd "$dir"
  if [[ -z "$(git status --porcelain)" ]]; then
    echo "[$name] No changes to commit."
    return 1
  fi

  echo "[$name] Staging and pushing changes on branch $branch..."
  git add -A
  git commit -m "wip: $(date +%Y-%m-%dT%H:%M:%S)"
  git push origin "$branch"
  echo "[$name] Pushed."
  return 0
}

cmd_push() {
  push_if_dirty "$CAPABILITIES_DIR" "$CAPABILITIES_BRANCH" || true
  push_if_dirty "$CHAINLINK_APTOS_DIR" "$CHAINLINK_APTOS_BRANCH" && APTOS_PUSHED=1 || true

  if [[ "$APTOS_PUSHED" -eq 1 ]]; then
    echo ""
    echo "chainlink-aptos changed — updating chainlink go.mod..."
    cd "$CHAINLINK_ROOT"
    go get "github.com/smartcontractkit/chainlink-aptos@${CHAINLINK_APTOS_BRANCH}"
    make gomodtidy
    echo "Done. chainlink go.mod updated."
  else
    echo ""
    echo "chainlink-aptos unchanged — skipping go get + tidy."
  fi
}

cmd_startenv() {
  cd "$SCRIPT_DIR"
  CTF_CONFIGS=./configs/workflow-gateway-don-aptos.toml \
  AWS_ECR=804282218731.dkr.ecr.us-west-2.amazonaws.com \
    go run . env start --auto-setup
}

cmd_stopenv() {
  cd "$SCRIPT_DIR"
  go run . env stop
}

usage() {
  cat <<'EOF'
Usage: aptos_dev.sh <command>

Commands:
  push       Commit+push capabilities & chainlink-aptos, then go get + tidy in chainlink
  startenv   Start CRE environment with Aptos topology
  stopenv    Stop CRE environment
EOF
}

case "${1:-}" in
  push)     cmd_push ;;
  startenv) cmd_startenv ;;
  stopenv)  cmd_stopenv ;;
  -h|--help) usage ;;
  *)
    usage
    exit 1
    ;;
esac
