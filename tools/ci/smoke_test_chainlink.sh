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

IMAGE_META="$( aws ecr describe-images --repository-name=chainlink --image-ids=imageTag=$APPS_CHAINLINK_VERSION 2> /dev/null )"

if [[ $? == 0 ]]; then
  IMAGE_TAGS="$( jq '.imageDetails[0].imageTags[0]' -r < ${IMAGE_META})"
  echo "$1:$2 found"
else
  echo "$1:$2 not found"
  echo "Not running any tests as no new image for this commit is available"
  exit 0
fi

# if [[ $performance == 0 ]]; then
#   ginkgo -r -p -keepGoing --trace --randomizeAllSpecs -skipPackage=./integration/suite/performance,./integration/suite/chaos ./integration/suite/...
# else
#   ginkgo -r -p -keepGoing --trace --randomizeAllSpecs ./integration/suite/performance ./integration/suite/chaos
# fi