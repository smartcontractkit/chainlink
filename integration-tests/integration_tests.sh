#!/bin/bash

# Inputs:
#   test_type: 0, 1, or 2
#     Indicates whether this is a standard integration, performance, or chaos test run or not.
#     - 0: Smoke tests
#     - 1: Performance tests
#     - 2: Chaos tests
test_type=$1

echo "Using image $APPS_CHAINLINK_IMAGE at version $APPS_CHAINLINK_VERSION"

if [[ $test_type == 0 ]]; then
  echo "Running smoke tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -nodes=15 -tags integration -skipPackage=./integration-tests/performance,./integration-tests/chaos ./integration-tests/...
elif [[ $test_type == 1 ]]; then
  echo "Running performance tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -nodes=5 -tags performance ./integration-tests/performance
elif [[ $test_type == 2 ]]; then
  echo "Running chaos tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -nodes=5 -tags chaos ./integration-tests/chaos
else
  echo "Invalid input '$test_type'"
  echo "
  Inputs:
   test_type: 0, 1, or 2
     Indicates whether this is a standard integration, performance, or chaos test run or not.
     - 0: Smoke tests
     - 1: Performance tests
     - 2: Chaos tests"
fi