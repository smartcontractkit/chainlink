#!/bin/bash
set -euo pipefail

# Ensure flakeguard is on PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Ensure required environment variables are set
: "${MATRIX_TYPE_CMD:?Environment variable MATRIX_TYPE_CMD is not set}"
: "${GITHUB_REPOSITORY:?Environment variable GITHUB_REPOSITORY is not set}"
: "${GITHUB_RUN_ID:?Environment variable GITHUB_RUN_ID is not set}"
: "${FLAKEGUARD_SPLUNK_ENDPOINT:?Environment variable FLAKEGUARD_SPLUNK_ENDPOINT is not set}"
: "${FLAKEGUARD_SPLUNK_HEC:?Environment variable FLAKEGUARD_SPLUNK_HEC is not set}"
: "${GITHUB_EVENT_NAME:?Environment variable GITHUB_EVENT_NAME is not set}"

######################################
# Retrieve Failed Logs Artifact Links
######################################

# MAIN failed logs artifact retrieval
MAIN_FAILED_LOGS_PATH="${MATRIX_TYPE_CMD}_flakeguard_report/main/failed-test-report-with-logs.json"
if [ -f "$MAIN_FAILED_LOGS_PATH" ]; then
  MAIN_FAILED_LOGS_ARTIFACT_NAME="${MATRIX_TYPE_CMD}_flakeguard_report_main_failed_tests.json"
  echo "Main failed-logs report file found at '$MAIN_FAILED_LOGS_PATH'."
else
  MAIN_FAILED_LOGS_ARTIFACT_NAME=""
  echo "No main failed-logs report file found at '$MAIN_FAILED_LOGS_PATH'."
fi

if [ -n "$MAIN_FAILED_LOGS_ARTIFACT_NAME" ]; then
  echo "Attempting to retrieve main failed-logs artifact link..."
  MAIN_ARTIFACT_LINK="$(flakeguard get-gh-artifact \
    --github-repository "$GITHUB_REPOSITORY" \
    --github-run-id "$GITHUB_RUN_ID" \
    --failed-tests-artifact-name "$MAIN_FAILED_LOGS_ARTIFACT_NAME" \
  )"
else
  MAIN_ARTIFACT_LINK=""
fi
echo "Main Artifact Link: $MAIN_ARTIFACT_LINK"

# RERUN failed logs artifact retrieval
RERUN_FAILED_LOGS_PATH="${MATRIX_TYPE_CMD}_flakeguard_report/rerun/failed-test-report-with-logs.json"
if [ -f "$RERUN_FAILED_LOGS_PATH" ]; then
  RERUN_FAILED_LOGS_ARTIFACT_NAME="${MATRIX_TYPE_CMD}_flakeguard_report_rerun_failed_tests.json"
  echo "Rerun failed-logs report file found at '$RERUN_FAILED_LOGS_PATH'."
else
  RERUN_FAILED_LOGS_ARTIFACT_NAME=""
  echo "No rerun failed-logs report file found at '$RERUN_FAILED_LOGS_PATH'."
fi

if [ -n "$RERUN_FAILED_LOGS_ARTIFACT_NAME" ]; then
  echo "Attempting to retrieve rerun failed-logs artifact link..."
  RERUN_ARTIFACT_LINK="$(flakeguard get-gh-artifact \
    --github-repository "$GITHUB_REPOSITORY" \
    --github-run-id "$GITHUB_RUN_ID" \
    --failed-tests-artifact-name "$RERUN_FAILED_LOGS_ARTIFACT_NAME" \
  )"
else
  RERUN_ARTIFACT_LINK=""
fi
echo "Rerun Artifact Link: $RERUN_ARTIFACT_LINK"

######################################
# Send Flakeguard Test Reports to Splunk
######################################

# Send main test report to Splunk if it exists
MAIN_FILE="${MATRIX_TYPE_CMD}_flakeguard_report/main/all-test-report.json"
if [ -f "$MAIN_FILE" ]; then
  echo "Main test report found at '$MAIN_FILE'. Sending main report to Splunk..."
  flakeguard send-to-splunk \
    --report-path "$MAIN_FILE" \
    --splunk-url "$FLAKEGUARD_SPLUNK_ENDPOINT" \
    --splunk-token "$FLAKEGUARD_SPLUNK_HEC" \
    --splunk-event "$GITHUB_EVENT_NAME" \
    --failed-logs-url "$MAIN_ARTIFACT_LINK"
  
  EXIT_CODE=$?
  if [ $EXIT_CODE -ne 0 ]; then
    echo "ERROR: Flakeguard encountered an error while sending the main report to Splunk"
    exit $EXIT_CODE
  fi
else
  echo "No main report found at '$MAIN_FILE'. Skipping Splunk send for main report."
fi

# Send rerun test report to Splunk if it exists
RERUN_FILE="${MATRIX_TYPE_CMD}_flakeguard_report/rerun/all-test-report.json"
if [ -f "$RERUN_FILE" ]; then
  echo "Rerun test report found at '$RERUN_FILE'. Sending rerun report to Splunk..."
  flakeguard send-to-splunk \
    --report-path "$RERUN_FILE" \
    --splunk-url "$FLAKEGUARD_SPLUNK_ENDPOINT" \
    --splunk-token "$FLAKEGUARD_SPLUNK_HEC" \
    --splunk-event "$GITHUB_EVENT_NAME" \
    --failed-logs-url "$RERUN_ARTIFACT_LINK"
  
  EXIT_CODE=$?
  if [ $EXIT_CODE -ne 0 ]; then
    echo "ERROR: Flakeguard encountered an error while sending the rerun report to Splunk"
    exit $EXIT_CODE
  fi
else
  echo "No rerun report found at '$RERUN_FILE'. Skipping Splunk send for rerun report."
fi