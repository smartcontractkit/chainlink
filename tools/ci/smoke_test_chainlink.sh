#!/bin/bash

# Inputs:
#   performance: 0 or 1
#     Indicates whether this is a performance test run or not.
#     - 0: No performance tests
#     - 1: Performance tests
performance=$1

echo "IMAGE: $APPS_CHAINLINK_IMAGE | TAG: $APPS_CHAINLINK_VERSION"

if [[ $performance == 0 ]]; then
  echo "Running smoke tests"
  ginkgo -r -p -keepGoing --trace --randomizeAllSpecs -skipPackage=./integration/suite/performance,./integration/suite/chaos ./integration/suite/...
else
  echo "Running performance and chaos tests"
  ginkgo -r -p -keepGoing --trace --randomizeAllSpecs ./integration/suite/performance ./integration/suite/chaos
fi