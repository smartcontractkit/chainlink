#!/bin/bash

# Inputs:
#   performance: 0 or 1
#     Indicates whether this is a performance test run or not.
#     - 0: No performance tests
#     - 1: Performance tests
performance=$1

export APPS_CHAINLINK_IMAGE=public.ecr.aws/z0b1w9r9/chainlink
export APPS_CHAINLINK_VERSION=latest.99415739365570301bb898f4b9c531e75b85a22d

echo "IMAGE: $APPS_CHAINLINK_IMAGE | TAG: $APPS_CHAINLINK_VERSION"

# Need to poll and see if an image for this commit gets pushed to ECR
# Obviously non-ideal, but having the build-and-push take place in a circle-ci step has introduced a host of issues
# trying to get this run correctly.
end_time=$((SECONDS+6))
while [[ $SECONDS -lt $end_time ]]
do 
  IMAGE_META=$( aws ecr describe-images --repository-name=chainlink --image-ids=imageTag=$APPS_CHAINLINK_VERSION )
  if [[ $? == 0 ]]; then
    IMAGE_TAGS="$( jq '.imageDetails[0].imageTags[0]' -r < ${IMAGE_META})"
    echo "$APPS_CHAINLINK_IMAGE:$APPS_CHAINLINK_VERSION found"
  fi
  sleep 3
done

if [[ $SECONDS -ge $end_time ]]; then
  echo "$APPS_CHAINLINK_IMAGE:$APPS_CHAINLINK_VERSION not found after $SECONDS seconds, not running any tests"
  exit 0
fi


if [[ $performance == 0 ]]; then
  echo "Running smoke tests"
  ginkgo -r -p -keepGoing --trace --randomizeAllSpecs -skipPackage=./integration/suite/performance,./integration/suite/chaos ./integration/suite/...
else
  echo "Running performance and chaos tests"
  ginkgo -r -p -keepGoing --trace --randomizeAllSpecs ./integration/suite/performance ./integration/suite/chaos
fi