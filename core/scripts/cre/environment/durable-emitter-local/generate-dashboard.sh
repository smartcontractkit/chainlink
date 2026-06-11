#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT="${SCRIPT_DIR}/grafana/dashboards/durable_emitter.json"

if [[ -z "${OBSERVABILITY_ROOT:-}" ]]; then
  echo "OBSERVABILITY_ROOT is required" >&2
  echo "Set it to your clone of github.com/smartcontractkit/chainlink-observability" >&2
  echo "Example: export OBSERVABILITY_ROOT=\$HOME/projects/chainlink-observability" >&2
  exit 1
fi

if [[ ! -d "${OBSERVABILITY_ROOT}/resources/durable_emitter" ]]; then
  echo "chainlink-observability not found at ${OBSERVABILITY_ROOT}" >&2
  echo "Clone https://github.com/smartcontractkit/chainlink-observability and set OBSERVABILITY_ROOT to that path" >&2
  exit 1
fi

mkdir -p "$(dirname "${OUTPUT}")"

GOMODCACHE="${GOMODCACHE:-${TMPDIR:-/tmp}/chainlink-observability-gomodcache}"
GOCACHE="${GOCACHE:-${TMPDIR:-/tmp}/chainlink-observability-gocache}"
mkdir -p "${GOMODCACHE}" "${GOCACHE}"

(
  cd "${OBSERVABILITY_ROOT}"
  GOMODCACHE="${GOMODCACHE}" GOCACHE="${GOCACHE}" go run ./cmd/generate-durable-emitter-dashboard/main.go \
    --local-load-test \
    --format grafana \
    --dashboard-uid durable-emitter-load-test \
    --output "${OUTPUT}"
)

echo "Generated ${OUTPUT}"
