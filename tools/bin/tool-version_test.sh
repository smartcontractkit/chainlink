#!/usr/bin/env bash
#
# Tests for tools/bin/tool-version. No network, no go install — uses the
# TOOL_VERSIONS_FILE / GO_TOOLS_FILE overrides to point at fixtures.
#
# Run: bash tools/bin/tool-version_test.sh

set -u

HERE="$(cd "$(dirname "$0")" && pwd)"
TOOL_VERSION="$HERE/tool-version"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

cat >"$tmp/.tool-versions" <<'EOF'
golang        1.26.4
mockery       2.53.0
protoc        29.3
golangci-lint 2.12.2
nodejs        20.13.1
EOF

cat >"$tmp/go-tools.txt" <<'EOF'
# comment line, ignored
github.com/jmank88/gomods            0.1.7
github.com/ugorji/go/codec/codecgen  1.2.10
github.com/smartcontractkit/gencodec 42dc7da8c2874db550e91c656f98d05fca3c2f98
EOF

export TOOL_VERSIONS_FILE="$tmp/.tool-versions"
export GO_TOOLS_FILE="$tmp/go-tools.txt"

fail=0
check() { # name, expected, actual
  if [ "$2" = "$3" ]; then
    echo "ok   - $1"
  else
    echo "FAIL - $1"
    echo "         expected: [$2]"
    echo "         actual:   [$3]"
    fail=1
  fi
}

check_rc() { # name, expected_rc, actual_rc
  if [ "$2" = "$3" ]; then
    echo "ok   - $1"
  else
    echo "FAIL - $1 (expected rc $2, got $3)"
    fail=1
  fi
}

# --- get: runtimes from .tool-versions ---
check "get mockery"        "2.53.0"   "$("$TOOL_VERSION" get mockery)"
check "get protoc"         "29.3"     "$("$TOOL_VERSION" get protoc)"
check "get golangci-lint"  "2.12.2"   "$("$TOOL_VERSION" get golangci-lint)"
check "get golang"         "1.26.4"   "$("$TOOL_VERSION" get golang)"

# --- get: CLIs from go-tools.txt (by import path) ---
check "get gomods path"    "0.1.7"    "$("$TOOL_VERSION" get github.com/jmank88/gomods)"
check "get codecgen path"  "1.2.10"   "$("$TOOL_VERSION" get github.com/ugorji/go/codec/codecgen)"

# --- target: the `go install` argument (module@version) ---
check "target mockery maps name->module, prepends v" \
  "github.com/vektra/mockery/v2@v2.53.0" "$("$TOOL_VERSION" target mockery)"
check "target gomods uses path verbatim, prepends v" \
  "github.com/jmank88/gomods@v0.1.7" "$("$TOOL_VERSION" target github.com/jmank88/gomods)"
check "target codecgen" \
  "github.com/ugorji/go/codec/codecgen@v1.2.10" "$("$TOOL_VERSION" target github.com/ugorji/go/codec/codecgen)"

# --- target: non-semver pin (SHA) is NOT given a v prefix ---
check "target gencodec SHA pin has no v prefix" \
  "github.com/smartcontractkit/gencodec@42dc7da8c2874db550e91c656f98d05fca3c2f98" \
  "$("$TOOL_VERSION" target github.com/smartcontractkit/gencodec)"

# --- unknown key fails ---
"$TOOL_VERSION" get does-not-exist >/dev/null 2>&1
check_rc "get unknown key exits non-zero" "1" "$?"

# --- target on a non-go runtime (protoc) fails (no module mapping) ---
"$TOOL_VERSION" target protoc >/dev/null 2>&1
check_rc "target protoc (no go module) exits non-zero" "1" "$?"

# --- list emits all pairs from both files ---
list_out="$("$TOOL_VERSION" list)"
case "$list_out" in
  *"mockery"*"2.53.0"*) echo "ok   - list includes mockery" ;;
  *) echo "FAIL - list missing mockery"; fail=1 ;;
esac
case "$list_out" in
  *"github.com/jmank88/gomods"*"0.1.7"*) echo "ok   - list includes gomods" ;;
  *) echo "FAIL - list missing gomods"; fail=1 ;;
esac

if [ "$fail" -eq 0 ]; then
  echo "PASS"
else
  echo "SOME TESTS FAILED"
fi
exit "$fail"
