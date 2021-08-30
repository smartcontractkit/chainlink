#!/bin/bash

set -ex

#
# Effectively replicates the decision logic found in push_chainlink
# This is done to get this CI runs' tag that was pushed to ECR so we can pull it down and test it.
# The build-and-push logic is done in circle-ci, thus necessitating this strategy.


# Our registry, hosted on Public ECR.
export APPS_CHAINLINK_IMAGE=public.ecr.aws/z0b1w9r9/chainlink

branch="$1"
tag="$2"
sha="$3"

branch_tag=$(tools/ci/branch2tag ${branch})     # ie: develop, latest, candidate-*, etc.
version_tag=$(tools/ci/gittag2dockertag ${tag}) # aka GIT_TAG. v0.9.1 -> 0.9.1

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
  ginkgo -r -p -keepGoing --trace --randomizeAllSpecs --progress ./integration/suite/...
}