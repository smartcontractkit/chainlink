#!/usr/bin/env bash
# Resolves the Chainlink Docker image based on ECR type and environment variables.
# Required:
#   - ECR_TYPE: one of "sdlc" or "public" (case-insensitive)
#   - CHAINLINK_IMAGE_REPO
#   - CHAINLINK_IMAGE_TAG
# For ECR_TYPE=sdlc, also required:
#   - AWS_ACCOUNT_NUMBER
#   - AWS_REGION
# Inputs are normalized to lowercase before validation and output.
set -euo pipefail

error() {
  echo "Error: $1" >&2
  exit 1
}

# Trim leading/trailing whitespace.
trim() {
  local s="$1"
  s="${s#"${s%%[![:space:]]*}"}"
  s="${s%"${s##*[![:space:]]}"}"
  printf '%s' "${s}"
}

repository="${CHAINLINK_IMAGE_REPO:-}"
image_tag="${CHAINLINK_IMAGE_TAG:-}"
ecr_type="${ECR_TYPE:-}"
aws_account="${AWS_ACCOUNT_NUMBER:-}"
aws_region="${AWS_REGION:-}"

# Trim whitespace before normalization/validation.
repository="$(trim "${repository}")"
image_tag="$(trim "${image_tag}")"
ecr_type="$(trim "${ecr_type}")"
aws_account="$(trim "${aws_account}")"
aws_region="$(trim "${aws_region}")"

# Normalize all input fields to lowercase for case-insensitive handling.
repository="$(printf '%s' "${repository}" | tr '[:upper:]' '[:lower:]')"
image_tag="$(printf '%s' "${image_tag}" | tr '[:upper:]' '[:lower:]')"
ecr_type="$(printf '%s' "${ecr_type}" | tr '[:upper:]' '[:lower:]')"
aws_account="$(printf '%s' "${aws_account}" | tr '[:upper:]' '[:lower:]')"
aws_region="$(printf '%s' "${aws_region}" | tr '[:upper:]' '[:lower:]')"

if [[ -z "${ecr_type}" ]]; then
  error "'ECR_TYPE' must be set and non-empty. Allowed values: 'sdlc' or 'public'."
fi

if [[ "${ecr_type}" != "sdlc" && "${ecr_type}" != "public" ]]; then
  error "Invalid 'ECR_TYPE': '${ecr_type}'. Allowed values: 'sdlc' or 'public'."
fi

if [[ -z "${repository}" ]]; then
  error "'CHAINLINK_IMAGE_REPO' must be set and non-empty."
fi

if [[ -z "${image_tag}" ]]; then
  error "'CHAINLINK_IMAGE_TAG' must be set and non-empty."
fi

if [[ "${ecr_type}" == "public" ]]; then
  printf '%s\n' "gallery.ecr.aws/${repository}:${image_tag}"
  exit 0
fi

# ECR_TYPE=sdlc
if [[ -z "${aws_account}" || -z "${aws_region}" ]]; then
  error "For 'ECR_TYPE=sdlc', both 'AWS_ACCOUNT_NUMBER' and 'AWS_REGION' must be set and non-empty."
fi

printf '%s\n' "${aws_account}.dkr.ecr.${aws_region}.amazonaws.com/${repository}:${image_tag}"
