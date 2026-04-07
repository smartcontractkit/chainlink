#!/usr/bin/env bash
# Build CRE test packages into .cre-test-bin/<path>/cre.test
# Run from chainlink repo root, or invoke via: (cd system-tests/tests && bash scripts/ci/build-cre-test-binaries.sh)
# Set BUILD_CRE_REGRESSION=true to also compile ./regression/cre.
set -euo pipefail

TESTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$TESTS_DIR"

rm -rf .cre-test-bin
mkdir -p .cre-test-bin/smoke/cre

echo "Building smoke/cre -> .cre-test-bin/smoke/cre/cre.test"
go test -c -cover -covermode=atomic -o .cre-test-bin/smoke/cre/cre.test ./smoke/cre

if [[ "${BUILD_CRE_REGRESSION:-}" == "true" ]]; then
  mkdir -p .cre-test-bin/regression/cre
  echo "Building regression/cre -> .cre-test-bin/regression/cre/cre.test"
  go test -c -cover -covermode=atomic -o .cre-test-bin/regression/cre/cre.test ./regression/cre
fi

echo "Done under $TESTS_DIR/.cre-test-bin"
