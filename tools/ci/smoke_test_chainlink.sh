#!/bin/bash

set -ex

#
# Effectively replicates the decision logic found in push_chainlink
# This is done to get this CI runs' tag that was pushed to ECR so we can pull it down and test it.
# The build-and-push logic is done in circle-ci, thus necessitating this strategy.


# Our registry, hosted on Public ECR.
export APPS_CHAINLINK_IMAGE=public.ecr.aws/z0b1w9r9/chainlink

# version tag takes precedence.
if [ -n "${version_tag}" ]; then
  # Only if we don't have an explorer tag
  if [[ "${tag}" =~ ^explorer-v([a-zA-Z0-9.]+) ]]; then
    echo "No image published for this PR, no tests run"
    exit 0
  else
    test "latest.$sha"
  fi
elif [ -n "$branch_tag" ]; then
  # Only if we're not on explorer branch
  if [[ "${branch}" =~ ^release(s)?\/explorer-(.+)$ ]]; then
    echo "No image published for this PR, no tests run"
    exit 0
  else
    test "$branch_tag.$sha"
  fi
else
    echo "No image published for this PR, no tests run"
    exit 0
fi

test() {
  export APPS_CHAINLINK_VERSION=$1 
  echo "Testing with the image $APPS_CHAINLINK_IMAGE:$APPS_CHAINLINK_VERSION"
}