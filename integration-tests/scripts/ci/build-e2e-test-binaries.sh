#!/usr/bin/env bash
# Build all integration-tests packages that contain tests into .e2e-test-bin/<dir>/<pkg>.test
# Run from the chainlink repository root (parent of integration-tests/).
set -euo pipefail

INTEGRATION_TESTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$INTEGRATION_TESTS_DIR"

rm -rf .e2e-test-bin
mkdir -p .e2e-test-bin

prefix="github.com/smartcontractkit/chainlink/integration-tests/"
mapfile -t pkgs < <(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)

for ip in "${pkgs[@]}"; do
  rel="${ip#"$prefix"}"
  if [[ "$rel" == "$ip" ]]; then
    echo "Skipping package outside integration-tests module: $ip"
    continue
  fi
  base=$(basename "$rel")
  outdir=".e2e-test-bin/$rel"
  mkdir -p "$outdir"
  echo "Building $ip -> $outdir/$base.test"
  go test -c -cover -covermode=atomic -o "$outdir/$base.test" "./$rel"
done

echo "Built ${#pkgs[@]} packages under $INTEGRATION_TESTS_DIR/.e2e-test-bin"
