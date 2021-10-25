#!/bin/bash

# Inputs:
#   performance: 0 or 1
#     Indicates whether this is a performance test run or not.
#     - 0: No performance tests
#     - 1: Performance tests
performance=$1

echo "Using image $APPS_CHAINLINK_IMAGE at version $APPS_CHAINLINK_VERSION"

if [[ $performance == 0 ]]; then
  echo "Running smoke tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress -nodes=15 -skipPackage=./integration-tests/performance,./integration-tests/chaos ./integration-tests/...
else
  echo "Running performance and chaos tests"
  ginkgo -r -keepGoing --trace --randomizeAllSpecs --randomizeSuites --progress ./integration-tests/performance ./integration-tests/chaos
fi