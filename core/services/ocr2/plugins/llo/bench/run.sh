#!/usr/bin/env bash
#
# Runs the v30-vs-v31 LLO plugin benchmark matrix and (if benchstat is present)
# prints a summarized table with per-benchmark variance.
#
# Usage:
#   ./run.sh [bench-regex] [count] [benchtime]
#
# Examples:
#   ./run.sh                              # full matrix, count=6, benchtime=1s
#   ./run.sh BenchmarkFullRound 10 2s     # just the full-round bench, more samples
#   ./run.sh 'FullRound/ch=10' 6 200x     # a single workload, fixed iterations
#
# Read the /v30 and /v31 rows for the same workload side by side; the v31 rows
# additionally carry precursor_B, kvread_B, kvwrite_B and kvkeys per op.
set -euo pipefail

BENCH="${1:-Benchmark}"
COUNT="${2:-6}"
BENCHTIME="${3:-1s}"

cd "$(dirname "$0")"

OUT="bench_results.txt"

echo "running: -bench '${BENCH}' -count=${COUNT} -benchtime=${BENCHTIME}"
go test . \
  -run '^$' \
  -bench "${BENCH}" \
  -benchmem \
  -count="${COUNT}" \
  -benchtime="${BENCHTIME}" \
  -timeout=60m | tee "${OUT}"

echo
if command -v benchstat >/dev/null 2>&1; then
  echo "=== benchstat (${OUT}) ==="
  benchstat "${OUT}"
else
  echo "benchstat not found; install with: go install golang.org/x/perf/cmd/benchstat@latest"
fi
