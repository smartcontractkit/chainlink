#!/usr/bin/env bash

set -euo pipefail

#############################
#                __________
#               < CRIBbit! >
#                ----------
#      _    _    /
#     (o)--(o)  /
#    /.______.\
#    \________/
#   ./        \.
#  ( .        , )
#   \ \_\\//_/ /
#    ~~  ~~  ~~
#
# Initialize your CRIB
# environment.
#############################

DEVSPACE_NAMESPACE="${1:-}"
if [[ -z "${DEVSPACE_NAMESPACE}" ]]; then
  echo "Usage: $0 <DEVSPACE_NAMESPACE>"
  exit 1
fi

# Bail if $DEVSPACE_NAMESPACE does not begin with a crib- prefix or does not have an override set.
if [[ ! "${DEVSPACE_NAMESPACE}" =~ ^crib- ]] && [[ -z "${CRIB_IGNORE_NAMESPACE_PREFIX:-}" ]]; then
  echo "Error: DEVSPACE_NAMESPACE must begin with 'crib-' prefix."
  exit 1
fi

# Path to the .env file
repo_root=$(git rev-parse --show-toplevel 2>/dev/null || echo ".")
env_file="${repo_root}/crib/.env"

# Source .env file if it exists
if [[ -f "${env_file}" ]]; then
  # shellcheck disable=SC1090
  source "${env_file}"
else
  echo "Error: .env file not found at $env_file"
  exit 1
fi

# List of required environment variables
required_vars=(
  "DEVSPACE_IMAGE"
  "HOME"
)

missing_vars=0 # Counter for missing variables

for var in "${required_vars[@]}"; do
  if [[ -z "${!var:-}" ]]; then # If variable is unset or empty
    echo "Error: Environment variable ${var} is not set."
    missing_vars=$((missing_vars + 1))
  fi
done

# Exit with an error if any variables were missing
if [[ $missing_vars -ne 0 ]]; then
  echo "Error: Total missing environment variables: $missing_vars"
  exit 1
fi

##
# Setup AWS Profile
##

path_aws_config="$HOME/.aws/config"
aws_account_id_ecr_registry=$(echo "${DEVSPACE_IMAGE}" | cut -d'.' -f1)
aws_profile_name="staging-crib"

if grep -q "$aws_profile_name" "$path_aws_config"; then
  echo "Info: Skip updating ${path_aws_config}. Profile already set: ${aws_profile_name}"
else
  # List of required environment variables
  required_aws_vars=(
    "AWS_REGION"
    # Should be the short name and not the full IAM role ARN.
    "AWS_SSO_ROLE_NAME"
    # The AWS SSO start URL, e.g. https://<org name>.awsapps.com/start
    "AWS_SSO_START_URL"
  )
  missing_aws_vars=0 # Counter for missing variables
  for var in "${required_aws_vars[@]}"; do
    if [[ -z "${!var:-}" ]]; then # If variable is unset or empty
      echo "Error: Environment variable ${var} is not set."
      missing_aws_vars=$((missing_aws_vars + 1))
    fi
  done

  # Exit with an error if any variables were missing
  if [[ $missing_aws_vars -ne 0 ]]; then
    echo "Error: Total missing environment variables: $missing_aws_vars"
    exit 1
  fi

  cat <<EOF >>"$path_aws_config"
[profile $aws_profile_name]
region=${AWS_REGION}
sso_start_url=${AWS_SSO_START_URL}
sso_region=${AWS_REGION}
sso_account_id=${aws_account_id_ecr_registry}
sso_role_name=${AWS_SSO_ROLE_NAME}
EOF
  echo "Info: ${path_aws_config} modified. Added profile: ${aws_profile_name}"
fi

echo "Info: Setting AWS Profile env var: AWS_PROFILE=${aws_profile_name}"
export AWS_PROFILE=${aws_profile_name}

if aws sts get-caller-identity >/dev/null 2>&1; then
  echo "Info: AWS credentials working."
else
  echo "Info: AWS credentials not detected. Attempting to login through SSO."
  aws sso login
fi

# Check again and fail this time if not successful
if ! aws sts get-caller-identity >/dev/null 2>&1; then
  echo "Error: AWS credentials still not detected. Exiting."
  exit 1
fi

##
# Setup EKS KUBECONFIG
##

path_kubeconfig="${KUBECONFIG:-$HOME/.kube/config}"
eks_cluster_name="${CRIB_EKS_CLUSTER_NAME:-main-stage-cluster}"
eks_alias_name="${CRIB_EKS_ALIAS_NAME:-main-stage-cluster-crib}"

if [[ ! -f "${path_kubeconfig}" ]] || ! grep -q "name: ${eks_alias_name}" "${path_kubeconfig}"; then
  echo "Info: KUBECONFIG file (${path_kubeconfig}) not found or alias (${eks_alias_name}) not found. Attempting to update kubeconfig."
  aws eks update-kubeconfig \
    --name "${eks_cluster_name}" \
    --alias "${eks_alias_name}" \
    --region "${AWS_REGION}"
else
  echo "Info: Alias '${eks_alias_name}' already exists in kubeconfig. No update needed."
  echo "Info: Setting kubernetes context to: ${eks_alias_name}"
  kubectl config use-context "${eks_alias_name}"
fi

##
# Check Docker Daemon
##

if docker info >/dev/null 2>&1; then
  echo "Info: Docker daemon is running, authorizing registry"
else
  echo "Error: Docker daemon is not running. Exiting."
  exit 1
fi

##
# AWS ECR Login
##

# Function to extract the host URI of the ECR registry from OCI URI
extract_ecr_host_uri() {
  local ecr_uri="$1"
  # Regex to capture the ECR host URI
  if [[ $ecr_uri =~ oci:\/\/([0-9]+\.dkr\.ecr\.[a-zA-Z0-9-]+\.amazonaws\.com) ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    echo "No valid ECR host URI found in the URI."
    echo "Have you set CHAINLINK_CLUSTER_HELM_CHART_URI env var?"
    exit 1
  fi
}

# Set env var CRIB_SKIP_ECR_LOGIN=true to skip ECR login.
if [[ -n "${CRIB_SKIP_ECR_LOGIN:-}" ]]; then
  echo "Info: Skipping ECR login."
else
  echo "Info: Logging docker into AWS ECR registry."
  aws ecr get-login-password \
    --region "${AWS_REGION}" |
    docker login --username AWS \
      --password-stdin "${aws_account_id_ecr_registry}.dkr.ecr.${AWS_REGION}.amazonaws.com"

  echo "Info: Logging helm into AWS ECR registry."
  helm_registry_uri=$(extract_ecr_host_uri "${CHAINLINK_CLUSTER_HELM_CHART_URI}")
  aws ecr get-login-password --region "${AWS_REGION}" |
    helm registry login "$helm_registry_uri" --username AWS --password-stdin
fi

##
# Setup DevSpace
##

devspace use namespace "${DEVSPACE_NAMESPACE}"
