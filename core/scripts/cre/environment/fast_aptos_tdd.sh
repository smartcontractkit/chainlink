#!/usr/bin/env bash
# Fast Aptos TDD loop for local CRE development.
#
# Run from anywhere:
#   ./core/scripts/cre/environment/fast_aptos_tdd.sh unit
#   ./core/scripts/cre/environment/fast_aptos_tdd.sh read
#   ./core/scripts/cre/environment/fast_aptos_tdd.sh all
#
# Notes:
# - E2E targets require local CRE to already be running (state/local_cre.toml present).
# - Keep this script updated as new Aptos tests are added.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHAINLINK_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
CAPABILITIES_ROOT="${CAPABILITIES_ROOT:-$(cd "$CHAINLINK_ROOT/.." && pwd)/capabilities}"
SYSTEM_TESTS_DIR="$CHAINLINK_ROOT/system-tests/tests"
CRE_STATE_FILE="$SCRIPT_DIR/state/local_cre.toml"

GOTEST_TIMEOUT="${GOTEST_TIMEOUT:-15m}"
GOTEST_COUNT="${GOTEST_COUNT:-1}"

usage() {
  cat <<'EOF'
Usage:
  fast_aptos_tdd.sh [target ...]

Targets:
  unit        Run fast unit/integration package tests (default)
  read        Run Aptos read-only CRE test
  suite       Run full Aptos CRE suite test
  all         Run: unit + read + suite
  list        Print targets

Environment variables:
  CAPABILITIES_ROOT  Override capabilities repo path
  GOTEST_TIMEOUT     Go test timeout for CRE tests (default: 15m)
  GOTEST_COUNT       Go test -count value (default: 1)
EOF
}

print_targets() {
  echo "Available targets: unit read suite all list"
}

run_go_test() {
  local dir="$1"
  shift
  echo
  echo "==> (cd $dir && go test $*)"
  (
    cd "$dir"
    go test "$@"
  )
}

require_cre_env() {
  if [[ -f "$CRE_STATE_FILE" ]]; then
    return 0
  fi

  echo "ERROR: CRE local state not found: $CRE_STATE_FILE"
  echo "Start local Aptos CRE first, e.g.:"
  echo "  cd $SCRIPT_DIR && ./run_aptos_full.sh"
  exit 1
}

run_unit() {
  if [[ ! -d "$CAPABILITIES_ROOT/chain_capabilities/aptos" ]]; then
    echo "ERROR: capabilities repo not found at $CAPABILITIES_ROOT"
    echo "Set CAPABILITIES_ROOT or clone capabilities adjacent to chainlink."
    exit 1
  fi

  run_go_test \
    "$CAPABILITIES_ROOT/chain_capabilities/aptos" \
    ./actions \
    -count="$GOTEST_COUNT"

  run_go_test \
    "$CHAINLINK_ROOT" \
    ./core/capabilities/remote/executable/request \
    -count="$GOTEST_COUNT"

  run_go_test \
    "$CHAINLINK_ROOT" \
    ./core/services/standardcapabilities \
    -run Test_getCapabilityID \
    -count="$GOTEST_COUNT"
}

run_cre_test() {
  local test_name="$1"
  require_cre_env
  run_go_test \
    "$SYSTEM_TESTS_DIR" \
    -timeout "$GOTEST_TIMEOUT" \
    -run "^${test_name}$" \
    ./smoke/cre/ \
    -count="$GOTEST_COUNT" \
    -v
}

run_target() {
  local target="$1"
  case "$target" in
    unit)
      run_unit
      ;;
    read)
      run_cre_test "Test_CRE_V2_Aptos_Read"
      ;;
    suite)
      run_cre_test "Test_CRE_V2_Aptos_Suite"
      ;;
    all)
      run_unit
      run_cre_test "Test_CRE_V2_Aptos_Read"
      run_cre_test "Test_CRE_V2_Aptos_Suite"
      ;;
    list)
      print_targets
      ;;
    *)
      echo "Unknown target: $target"
      echo
      usage
      exit 1
      ;;
  esac
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ $# -eq 0 ]]; then
  run_target unit
  exit 0
fi

for target in "$@"; do
  run_target "$target"
done
