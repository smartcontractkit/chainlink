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

test_full_image_tag() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_FULL_IMAGE=public.ecr.aws/chainlink/chainlink:2.0.0"
  assert_eq "${RUN_STATUS}" "0" "full image with tag exits 0"
  assert_eq "${RUN_STDOUT}" "public.ecr.aws/chainlink/chainlink:2.0.0" "full image is returned to stdout"
  assert_eq "${RUN_STDERR}" "" "full image success does not write stderr"
}

test_full_image_digest() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_FULL_IMAGE=public.ecr.aws/chainlink/chainlink@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  assert_eq "${RUN_STATUS}" "0" "full image with digest exits 0"
  assert_eq "${RUN_STDOUT}" "public.ecr.aws/chainlink/chainlink@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" "digest image is returned to stdout"
}

test_repo_and_tag() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_IMAGE_REPO=chainlink-integration-tests" \
    "CHAINLINK_IMAGE_TAG=v2.1.0" \
    "AWS_ACCOUNT_NUMBER=123456789012" \
    "AWS_REGION=us-west-2"
  assert_eq "${RUN_STATUS}" "0" "repo and tag exits 0"
  assert_eq "${RUN_STDOUT}" "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0" "repo and tag resolve to ECR image"
  assert_eq "${RUN_STDERR}" "" "repo/tag success does not write stderr"
}

test_default_repo_and_version() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_VERSION=abc123" \
    "AWS_ACCOUNT_NUMBER=123456789012" \
    "AWS_REGION=us-west-2"
  assert_eq "${RUN_STATUS}" "0" "version fallback exits 0"
  assert_eq "${RUN_STDOUT}" "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink:abc123" "default repo is chainlink"
}

test_tag_wins_over_version() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_IMAGE_TAG=explicit-tag" \
    "CHAINLINK_VERSION=abc123" \
    "AWS_ACCOUNT_NUMBER=123456789012" \
    "AWS_REGION=us-west-2"
  assert_eq "${RUN_STATUS}" "0" "tag override exits 0"
  assert_eq "${RUN_STDOUT}" "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink:explicit-tag" "tag wins over version"
}

test_mutual_exclusion_full_and_tag() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_FULL_IMAGE=public.ecr.aws/chainlink/chainlink:2.0.0" \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "full image and tag exits 1"
  assert_contains "${RUN_STDERR}" "mutually exclusive" "mutual exclusion error is reported"
}

test_invalid_full_image() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_FULL_IMAGE=chainlink:2.0.0"
  assert_eq "${RUN_STATUS}" "1" "invalid full image exits 1"
  assert_contains "${RUN_STDERR}" "Invalid 'CHAINLINK_FULL_IMAGE' format" "invalid full image format is reported"
}

test_missing_tag_and_version() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "AWS_ACCOUNT_NUMBER=123456789012" \
    "AWS_REGION=us-west-2"
  assert_eq "${RUN_STATUS}" "1" "missing tag and version exits 1"
  assert_contains "${RUN_STDERR}" "Either provide 'CHAINLINK_FULL_IMAGE'" "missing tag/version error is reported"
}

test_missing_aws_envs() {
  TESTS_RUN=$((TESTS_RUN + 1))
  run_script \
    "CHAINLINK_IMAGE_TAG=v2.1.0"
  assert_eq "${RUN_STATUS}" "1" "missing AWS env vars exits 1"
  assert_contains "${RUN_STDERR}" "'AWS_ACCOUNT_NUMBER' and 'AWS_REGION'" "missing AWS env vars error is reported"
}

main() {
  test_full_image_tag
  test_full_image_digest
  test_repo_and_tag
  test_default_repo_and_version
  test_tag_wins_over_version
  test_mutual_exclusion_full_and_tag
  test_invalid_full_image
  test_missing_tag_and_version
  test_missing_aws_envs

  if [[ "${TESTS_FAILED}" -ne 0 ]]; then
    echo
    echo "Tests failed: ${TESTS_FAILED}/${TESTS_RUN}"
    exit 1
  fi

  echo "All tests passed: ${TESTS_RUN}"
}

main
