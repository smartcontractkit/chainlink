#!/usr/bin/env bash
set -euo pipefail

##
# Tests for ./ontriggerlint.sh
##

_compare_result() {
    local TEST_NAME=$1
    local EXPECTED_RESULT=$2
    local ACTUAL_RESULT=$3

    if [[ "${EXPECTED_RESULT:-}" != "${ACTUAL_RESULT:-}" ]]; then
        echo "Fail: ${TEST_NAME} expected ${EXPECTED_RESULT} got ${ACTUAL_RESULT}"
        exit 1
    else
         echo "Pass: ${TEST_NAME} expected ${EXPECTED_RESULT} got ${ACTUAL_RESULT}"
    fi
}

test_schedule() {
    local TEST_NAME="Trigger on schedule"
    local EXPECTED_RESULT="on_trigger=true"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=schedule ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"
}

test_pull_request() {
    local TEST_NAME="Trigger on pull_request for non release branches"
    local EXPECTED_RESULT="on_trigger=true"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=pull_request GITHUB_BASE_REF=develop ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"

    local TEST_NAME="No trigger on pull_request for release branches"
    local EXPECTED_RESULT="on_trigger=false"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=pull_request GITHUB_BASE_REF=release/1.2.3 ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"
}

test_push() {
    local TEST_NAME="Trigger on push to the staging branch"
    local EXPECTED_RESULT="on_trigger=true"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=push GITHUB_REF=staging ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"

    local TEST_NAME="Trigger on push to the trying branch"
    local EXPECTED_RESULT="on_trigger=true"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=push GITHUB_REF=staging ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"

    local TEST_NAME="Trigger on push to the rollup branch"
    local EXPECTED_RESULT="on_trigger=true"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=push GITHUB_REF=rollup ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"

    local TEST_NAME="No trigger on push to the develop branch"
    local EXPECTED_RESULT="on_trigger=false"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=push GITHUB_REF=develop ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"
}

test_misc() {
    local TEST_NAME="No trigger on invalid event name"
    local EXPECTED_RESULT="on_trigger=false"
    ACTUAL_RESULT=$(GITHUB_EVENT_NAME=invalid_event_name ./ontriggerlint.sh)
    _compare_result "${TEST_NAME}" "${EXPECTED_RESULT}" "${ACTUAL_RESULT}"    
}

test_schedule
test_pull_request
test_push
test_misc
