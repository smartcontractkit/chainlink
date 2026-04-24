#!/usr/bin/env bash

BASE_IMAGE="$REGISTRY.dkr.ecr.us-west-2.amazonaws.com/test-base-image"
IMAGE_VERSION="latest"
SUITES="chaos migration reorg smoke load benchmark ccip-tests/load ccip-tests/smoke ccip-tests/chaos"
TAG_NAME="chainlink-tests:latest"
DOCKERFILE="test.Dockerfile"

docker build \
  --platform linux/amd64 \
  --build-arg BASE_IMAGE="$BASE_IMAGE" \
  --build-arg IMAGE_VERSION="$IMAGE_VERSION" \
  --build-arg SUITES="$SUITES" \
  -t "$TAG_NAME" \
  -f "$DOCKERFILE" ..