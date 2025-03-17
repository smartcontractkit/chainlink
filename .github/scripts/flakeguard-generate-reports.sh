#!/bin/bash
set -euo pipefail

# Fix flakeguard binary path
export PATH="$PATH:$(go env GOPATH)/bin"

# Ensure required environment variables are set
: "${MATRIX_TYPE_CMD:?Environment variable MATRIX_TYPE_CMD is not set}"
: "${GITHUB_WORKSPACE:?Environment variable GITHUB_WORKSPACE is not set}"
: "${GITHUB_SHA:?Environment variable GITHUB_SHA is not set}"
: "${GITHUB_WORKFLOW:?Environment variable GITHUB_WORKFLOW is not set}"
: "${GITHUB_SERVER_URL:?Environment variable GITHUB_SERVER_URL is not set}"
: "${GITHUB_REPOSITORY:?Environment variable GITHUB_REPOSITORY is not set}"
: "${GITHUB_RUN_ID:?Environment variable GITHUB_RUN_ID is not set}"

# Determine branch name (preferring GITHUB_HEAD_REF over GITHUB_REF_NAME)
BRANCH_NAME="${GITHUB_HEAD_REF:-$GITHUB_REF_NAME}"

# Set the project path conditionally based on MATRIX_TYPE_CMD
if [ "$MATRIX_TYPE_CMD" = "go_core_ccip_deployment_tests" ]; then
  PROJECT_PATH="deployment"
else
  PROJECT_PATH="."
fi

# Generate a UUID for the main report
MAIN_REPORT_ID=$(uuidgen)
echo "Using Main Report ID: $MAIN_REPORT_ID"

# Generate the main test report
flakeguard generate-test-report \
  --test-results-dir "${MATRIX_TYPE_CMD}_flakeguard_results/main/test_results.json" \
  --output-path "${MATRIX_TYPE_CMD}_flakeguard_report/main" \
  --codeowners-path "$GITHUB_WORKSPACE/.github/CODEOWNERS" \
  --repo-path "$GITHUB_WORKSPACE" \
  --branch-name "$BRANCH_NAME" \
  --head-sha "$GITHUB_SHA" \
  --github-workflow-name "$GITHUB_WORKFLOW" \
  --github-workflow-run-url "${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}" \
  --project-path "$PROJECT_PATH" \
  --report-id "$MAIN_REPORT_ID"

EXIT_CODE=$?
if [ $EXIT_CODE -ne 0 ]; then
  echo "ERROR: Flakeguard encountered an error while generating the main test report"
  exit $EXIT_CODE
fi

# Check for rerun test results and generate a rerun report if available.
RERUN_FILE="${MATRIX_TYPE_CMD}_flakeguard_results/rerun/test_results.json"
if [ -f "$RERUN_FILE" ]; then
  echo "Rerun test results found. Generating rerun test report..."
  flakeguard generate-test-report \
    --test-results-dir "$RERUN_FILE" \
    --output-path "${MATRIX_TYPE_CMD}_flakeguard_report/rerun" \
    --codeowners-path "$GITHUB_WORKSPACE/.github/CODEOWNERS" \
    --repo-path "$GITHUB_WORKSPACE" \
    --branch-name "$BRANCH_NAME" \
    --head-sha "$GITHUB_SHA" \
    --github-workflow-name "$GITHUB_WORKFLOW" \
    --github-workflow-run-url "${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}" \
    --project-path "$PROJECT_PATH" \
    --gen-report-id \
    --rerun-of-report-id "$MAIN_REPORT_ID"

  EXIT_CODE=$?
  if [ $EXIT_CODE -ne 0 ]; then
    echo "ERROR: Flakeguard encountered an error while generating the rerun test report"
    exit $EXIT_CODE
  fi
fi