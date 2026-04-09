# Default package parallelism for `go test ./...` matches `go help build`:
#   -p n  ... The default for -p is GOMAXPROCS, normally the number of CPUs available.
# Honor the same inputs: optional CI_GO_TEST_PROFILE_PARALLEL, then $GOMAXPROCS, then go env, then OS CPU count.
go_test_package_parallelism() {
  local n
  if [[ -n "${CI_GO_TEST_PROFILE_PARALLEL:-}" ]] && [[ "${CI_GO_TEST_PROFILE_PARALLEL}" =~ ^[0-9]+$ ]] && [[ "${CI_GO_TEST_PROFILE_PARALLEL}" -gt 0 ]]; then
    echo "${CI_GO_TEST_PROFILE_PARALLEL}"
    return
  fi
  if [[ -n "${GOMAXPROCS:-}" ]] && [[ "${GOMAXPROCS}" =~ ^[0-9]+$ ]] && [[ "${GOMAXPROCS}" -gt 0 ]]; then
    echo "${GOMAXPROCS}"
    return
  fi
  n=$(go env GOMAXPROCS 2>/dev/null | tr -d ' \r\n\t' || true)
  if [[ -n "$n" ]] && [[ "$n" =~ ^[0-9]+$ ]] && [[ "$n" -gt 0 ]]; then
    echo "$n"
    return
  fi
  if command -v nproc >/dev/null 2>&1; then
    nproc
  elif command -v getconf >/dev/null 2>&1; then
    getconf _NPROCESSORS_ONLN 2>/dev/null || echo 1
  elif [[ "$(uname -s)" == "Darwin" ]]; then
    sysctl -n hw.logicalcpu 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 1
  else
    echo 1
  fi
}
