#!/usr/bin/env bash
# Tests for resolve-ccip-release-baseline.sh.
#
# The derivation + selection logic is exercised with a stubbed `docker` so no
# network access is required: the stub treats any tag listed in EXISTING_TAGS as
# published and everything else as absent.
#
# With LIVE_ECR=1, an additional set of cases runs against the real public ECR
# (public.ecr.aws/chainlink/ccip) to validate the probe path end-to-end.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RESOLVER="${SCRIPT_DIR}/resolve-ccip-release-baseline.sh"

TESTS_RUN=0
TESTS_FAILED=0

RUN_STATUS=0
RUN_STDOUT=""
RUN_STDERR=""

# run_resolver <CHAINLINK_VERSION> <existing-tags-space-separated>
# Uses a stubbed docker that publishes exactly the listed tags.
run_resolver() {
  local ver="$1"
  local existing="$2"
  local tmpbin
  tmpbin="$(mktemp -d)"
  cat >"${tmpbin}/docker" <<STUB
#!/usr/bin/env bash
img="\${@: -1}"
tag="\${img##*:}"
case " \${EXISTING_TAGS:-} " in
  *" \${tag} "*) exit 0;;
  *) exit 1;;
esac
STUB
  chmod +x "${tmpbin}/docker"

  local stdout_file stderr_file
  stdout_file="$(mktemp)"
  stderr_file="$(mktemp)"

  set +e
  env -i PATH="${tmpbin}:${PATH}" \
    CHAINLINK_VERSION="${ver}" \
    EXISTING_TAGS="${existing}" \
    PROBE_ATTEMPTS="1" \
    PROBE_RETRY_DELAY="0" \
    bash "${RESOLVER}" >"${stdout_file}" 2>"${stderr_file}"
  RUN_STATUS=$?
  set -e

  RUN_STDOUT="$(<"${stdout_file}")"
  RUN_STDERR="$(<"${stderr_file}")"
  rm -f "${stdout_file}" "${stderr_file}"
  rm -rf "${tmpbin}"
}

# run_resolver_live <CHAINLINK_VERSION> - real docker, real public ECR.
run_resolver_live() {
  local ver="$1"
  local stdout_file stderr_file
  stdout_file="$(mktemp)"
  stderr_file="$(mktemp)"

  set +e
  env -i PATH="${PATH}" \
    CHAINLINK_VERSION="${ver}" \
    bash "${RESOLVER}" >"${stdout_file}" 2>"${stderr_file}"
  RUN_STATUS=$?
  set -e

  RUN_STDOUT="$(<"${stdout_file}")"
  RUN_STDERR="$(<"${stderr_file}")"
  rm -f "${stdout_file}" "${stderr_file}"
}

get_val() { # <key>
  echo "${RUN_STDOUT}" | grep -E "^${1}=" | head -1 | cut -d= -f2-
}

assert_eq() {
  local got="$1"
  local want="$2"
  local msg="$3"
  TESTS_RUN=$((TESTS_RUN + 1))
  if [[ "${got}" != "${want}" ]]; then
    echo "FAIL: ${msg}"
    echo "  expected: ${want}"
    echo "  got:      ${got}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

assert_match() {
  local got="$1"
  local pattern="$2"
  local msg="$3"
  TESTS_RUN=$((TESTS_RUN + 1))
  if ! [[ "${got}" =~ ${pattern} ]]; then
    echo "FAIL: ${msg}"
    echo "  expected to match: ${pattern}"
    echo "  got:               ${got}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

# --- stubbed tests ---

test_rc_n_uses_same_version_rc0() {
  run_resolver "v2.63.1-rc.4" "2.63.1-ccip-rc.0"
  assert_eq "$(get_val ccip_pr_tag)" "2.63.1-ccip-rc.4" "rcN pr tag derived"
  assert_eq "$(get_val baseline_image_tag)" "2.63.1-ccip-rc.0" "rcN baseline is same-version rc0"
  assert_eq "$(get_val skip)" "false" "rcN not skipped when rc0 published"
}

test_rc_n_falls_back_to_minor_base_rc0() {
  # same-version rc0 NOT published -> fall back to this minor's .0 rc0
  run_resolver "v2.63.1-rc.4" "2.63.0-ccip-rc.0"
  assert_eq "$(get_val baseline_image_tag)" "2.63.0-ccip-rc.0" "rcN falls back to minor base rc0"
  assert_eq "$(get_val skip)" "false" "rcN fallback not skipped"
}

test_rc0_patch_uses_minor_base_rc0() {
  run_resolver "v2.63.1-rc.0" "2.63.0-ccip-rc.0"
  assert_eq "$(get_val ccip_pr_tag)" "2.63.1-ccip-rc.0" "rc0 patch pr tag derived"
  assert_eq "$(get_val baseline_image_tag)" "2.63.0-ccip-rc.0" "rc0 patch baseline is minor base rc0"
  assert_eq "$(get_val skip)" "false" "rc0 patch not skipped"
}

test_rc0_minor_uses_previous_minor_rc0() {
  run_resolver "v2.63.0-rc.0" "2.62.0-ccip-rc.0"
  assert_eq "$(get_val ccip_pr_tag)" "2.63.0-ccip-rc.0" "rc0 minor pr tag derived"
  assert_eq "$(get_val baseline_image_tag)" "2.62.0-ccip-rc.0" "rc0 minor baseline is previous minor rc0"
  assert_eq "$(get_val skip)" "false" "rc0 minor not skipped"
}

test_walks_past_missing_minors() {
  # 2.62 and 2.61 rc0 absent -> walk to 2.60.0-ccip-rc.0
  run_resolver "v2.63.0-rc.0" "2.60.0-ccip-rc.0"
  assert_eq "$(get_val baseline_image_tag)" "2.60.0-ccip-rc.0" "walks past missing minor rc0s"
  assert_eq "$(get_val skip)" "false" "walk fallback not skipped"
}

test_skips_when_nothing_published() {
  run_resolver "v2.63.0-rc.0" ""
  assert_eq "$(get_val baseline_image_tag)" "" "no baseline tag when nothing published"
  assert_eq "$(get_val skip)" "true" "skips when nothing published"
  # ::warning:: is emitted on stdout so GitHub Actions recognizes the annotation.
  assert_match "${RUN_STDOUT}" "No published baseline" "skip emits a warning"
}

test_stable_derives_no_ccip_suffix() {
  run_resolver "v2.63.0" "2.62.0-ccip-rc.0"
  assert_eq "$(get_val ccip_pr_tag)" "2.63.0" "stable pr tag has no -ccip suffix"
  assert_eq "$(get_val baseline_image_tag)" "2.62.0-ccip-rc.0" "stable baseline is previous minor rc0"
  assert_eq "$(get_val skip)" "false" "stable not skipped"
}

test_beta_derives_ccip_beta_tag() {
  run_resolver "v2.63.0-beta.0" "2.62.0-ccip-rc.0"
  assert_eq "$(get_val ccip_pr_tag)" "2.63.0-ccip-beta.0" "beta pr tag derived"
  assert_eq "$(get_val baseline_image_tag)" "2.62.0-ccip-rc.0" "beta baseline is previous minor rc0"
}

test_rejects_non_v_tag() {
  run_resolver "2.63.0-rc.0" ""
  assert_eq "${RUN_STATUS}" "1" "non-v tag exits 1"
  assert_match "${RUN_STDERR}" "must begin with 'v'" "non-v tag error reported"
}

# --- live tests (only with LIVE_ECR=1) ---

test_live_recent_rc_resolves() {
  [[ "${LIVE_ECR:-0}" == "1" ]] || return 0
  local resolved_any=false
  for v in "v2.63.0-rc.0" "v2.62.0-rc.0" "v2.61.0-rc.0"; do
    run_resolver_live "${v}"
    assert_eq "${RUN_STATUS}" "0" "live resolver exits 0 for ${v}"
    assert_match "$(get_val ccip_pr_tag)" "^[0-9].*-ccip-rc\.0$" "live pr tag well-formed for ${v}"
    if [[ "$(get_val skip)" == "false" ]]; then
      assert_match "$(get_val baseline_image_tag)" "^[0-9].*-ccip-rc\.0$" "live baseline well-formed for ${v}"
      resolved_any=true
    fi
  done
  TESTS_RUN=$((TESTS_RUN + 1))
  if [[ "${resolved_any}" != "true" ]]; then
    echo "FAIL: expected at least one recent rc to resolve a live baseline"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

main() {
  test_rc_n_uses_same_version_rc0
  test_rc_n_falls_back_to_minor_base_rc0
  test_rc0_patch_uses_minor_base_rc0
  test_rc0_minor_uses_previous_minor_rc0
  test_walks_past_missing_minors
  test_skips_when_nothing_published
  test_stable_derives_no_ccip_suffix
  test_beta_derives_ccip_beta_tag
  test_rejects_non_v_tag
  test_live_recent_rc_resolves

  if [[ "${TESTS_FAILED}" -ne 0 ]]; then
    echo
    echo "Tests failed: ${TESTS_FAILED}/${TESTS_RUN}"
    exit 1
  fi

  echo "All tests passed: ${TESTS_RUN}"
}

main
