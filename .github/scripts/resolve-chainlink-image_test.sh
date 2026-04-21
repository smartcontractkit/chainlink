#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_SCRIPT="${SCRIPT_DIR}/resolve-chainlink-image.sh"

TESTS_RUN=0
TESTS_FAILED=0

RUN_STATUS=0
RUN_STDOUT=""
RUN_STDERR=""

run_script() {
  local -a env_vars=("$@")
  local stdout_file
  local stderr_file
  stdout_file="$(mktemp)"
  stderr_file="$(mktemp)"

  set +e
  env -i PATH="${PATH}" "${env_vars[@]}" bash "${TARGET_SCRIPT}" >"${stdout_file}" 2>"${stderr_file}"
  RUN_STATUS=$?
  set -e

  RUN_STDOUT="$(<"${stdout_file}")"
  RUN_STDERR="$(<"${stderr_file}")"
  rm -f "${stdout_file}" "${stderr_file}"
}

assert_eq() {
  local got="$1"
  local want="$2"
  local msg="$3"
  if [[ "${got}" != "${want}" ]]; then
    echo "FAIL: ${msg}"
    echo "  expected: ${want}"
    echo "  got:      ${got}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  local msg="$3"
  if [[ "${haystack}" != *"${needle}"* ]]; then
    echo "FAIL: ${msg}"
    echo "  expected substring: ${needle}"
    echo "  got:                ${haystack}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
}

test_public_ecr() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=public" \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "0" "public ecr exits 0"
  assert_eq "${RUN_STDOUT}" "public.ecr.aws/chainlink:v2.1.0" "public image is returned to stdout"
  assert_eq "${RUN_STDERR}" "" "public ecr success does not write stderr"
}

test_private_ecr() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=sdlc" \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink-integration-tests" \
    "CHAINLINK_IMAGE_TAG=v2.1.0" \
    "AWS_ACCOUNT_NUMBER=123456789012" \
    "AWS_REGION=us-west-2"
  assert_eq "${RUN_STATUS}" "0" "repo and tag exits 0"
  assert_eq "${RUN_STDOUT}" "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0" "repo and tag resolve to ECR image"
  assert_eq "${RUN_STDERR}" "" "repo/tag success does not write stderr"
}

test_inputs_are_case_insensitive() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=PuBLic" \
    "CHAINLINK_IMAGE_REPO_PATH=ChAinLink" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "0" "case-insensitive inputs exit 0"
  assert_eq "${RUN_STDOUT}" "public.ecr.aws/chainlink:v2.1.0" "inputs are lowercased"
}

test_inputs_are_case_insensitive_but_not_tag() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=PuBLic" \
    "CHAINLINK_IMAGE_REPO_PATH=ChAinLink" \
    "CHAINLINK_IMAGE_TAG=V2.1.0"
  assert_eq "${RUN_STATUS}" "0" "case-insensitive inputs exit 0"
  assert_eq "${RUN_STDOUT}" "public.ecr.aws/chainlink:V2.1.0" "tag is not lowercased"
}

test_whitespace_trimming() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=  SDLC  " \
    "CHAINLINK_IMAGE_REPO_PATH=  chainlink-integration-tests  " \
    "CHAINLINK_IMAGE_TAG=  v2.1.0  " \
    "AWS_ACCOUNT_NUMBER=  123456789012  " \
    "AWS_REGION=  US-WEST-2  "
  assert_eq "${RUN_STATUS}" "0" "whitespace-padded inputs exit 0"
  assert_eq "${RUN_STDOUT}" "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0" "inputs are trimmed and lowercased"
}

test_missing_ecr_type() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "missing ecr type exits 1"
  assert_contains "${RUN_STDERR}" "'ECR_TYPE' must be set" "missing ecr type error is reported"
}

test_invalid_ecr_type() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=other" \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "invalid ecr type exits 1"
  assert_contains "${RUN_STDERR}" "Invalid 'ECR_TYPE'" "invalid ecr type error is reported"
}

test_missing_repo() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=public" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "missing repo exits 1"
  assert_contains "${RUN_STDERR}" "'CHAINLINK_IMAGE_REPO_PATH' must be set" "missing repo error is reported"
}

test_missing_tag() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=public" \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink"
  assert_eq "${RUN_STATUS}" "1" "missing tag exits 1"
  assert_contains "${RUN_STDERR}" "'CHAINLINK_IMAGE_TAG' must be set" "missing tag error is reported"
}

test_missing_aws_envs() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "ECR_TYPE=sdlc" \
    "CHAINLINK_IMAGE_REPO_PATH=chainlink" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "missing AWS env vars exits 1"
  assert_contains "${RUN_STDERR}" "For 'ECR_TYPE=sdlc'" "missing AWS env vars error is reported"
}

main() {
  test_public_ecr
  test_private_ecr
  test_inputs_are_case_insensitive
  test_inputs_are_case_insensitive_but_not_tag
  test_whitespace_trimming
  test_missing_ecr_type
  test_invalid_ecr_type
  test_missing_repo
  test_missing_tag
  test_missing_aws_envs

  if [[ "${TESTS_FAILED}" -ne 0 ]]; then
    echo
    echo "Tests failed: ${TESTS_FAILED}/${TESTS_RUN}"
    exit 1
  fi

  echo "All tests passed: ${TESTS_RUN}"
}

main
