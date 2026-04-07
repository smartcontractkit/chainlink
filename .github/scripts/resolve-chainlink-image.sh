#!/usr/bin/env bash
# Resolves the Chainlink Docker image to use based on environment variables.
# Priority:
#   1. If CHAINLINK_FULL_IMAGE is set, use it directly (must include registry and tag/digest).
#   2. Otherwise, construct the image from:
#        - Repository: CHAINLINK_IMAGE_REPO (default: "chainlink")
#        - Tag: CHAINLINK_IMAGE_TAG or fallback to CHAINLINK_VERSION
#        - Registry: AWS ECR using AWS_ACCOUNT_NUMBER and AWS_REGION
# Enforces mutual exclusivity between FULL_IMAGE and repo/tag inputs, and validates required fields.
set -euo pipefail

error() {
  echo "Error: $1" >&2
  exit 1
}

full_image="${CHAINLINK_FULL_IMAGE:-}"
explicit_repo="${CHAINLINK_IMAGE_REPO:-}"
tag="${CHAINLINK_IMAGE_TAG:-}"
version="${CHAINLINK_VERSION:-}"
aws_account="${AWS_ACCOUNT_NUMBER:-}"
aws_region="${AWS_REGION:-}"

resolved_repo="${explicit_repo:-chainlink}"

if [[ -n "${full_image}" ]]; then
  if [[ -n "${tag}" || -n "${explicit_repo}" ]]; then
    error "'CHAINLINK_FULL_IMAGE' is mutually exclusive with 'CHAINLINK_IMAGE_TAG' and 'CHAINLINK_IMAGE_REPO'."
  fi

  if [[ ! "${full_image}" =~ .+/.+(:[^@]+|@sha256:[a-f0-9]{64})$ ]]; then
    error "Invalid 'CHAINLINK_FULL_IMAGE' format: '${full_image}'. Expected '<registry>/<repo>:<tag>' or '<registry>/<repo>@sha256:<digest>'."
  fi

  printf '%s\n' "${full_image}"
  exit 0
fi

resolved_tag="${tag:-$version}"
if [[ -z "${resolved_tag}" ]]; then
  error "Either provide 'CHAINLINK_FULL_IMAGE' or provide a tag via 'CHAINLINK_IMAGE_TAG' (or 'CHAINLINK_VERSION')."
fi

if [[ -z "${aws_account}" || -z "${aws_region}" ]]; then
  error "When 'CHAINLINK_FULL_IMAGE' is not provided, both 'AWS_ACCOUNT_NUMBER' and 'AWS_REGION' environment variables must be set and non-empty."
fi

printf '%s\n' "${aws_account}.dkr.ecr.${aws_region}.amazonaws.com/${resolved_repo}:${resolved_tag}"
