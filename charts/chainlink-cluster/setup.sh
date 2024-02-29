#!/usr/bin/env bash

set -e
user_home="$HOME"
file_path="$user_home/.aws/config"
image=""
registry_id=$(echo "$DEVSPACE_IMAGE" | cut -d'.' -f1)

if grep -q "staging-crib" "$file_path"; then
  echo "Staging AWS config is already applied, role is 'staging-crib'"
else
  cat <<EOF >> "$file_path"
[profile staging-crib]
region=us-west-2
sso_start_url=https://smartcontract.awsapps.com/start
sso_region=us-west-2
sso_account_id=${registry_id}
sso_role_name=CRIB-ECR-Power
EOF
  echo "~/.aws/config modified, added 'staging-crib"
fi

# Login through SSO
aws sso login --profile staging-crib
# Update kubeconfig and switch context
export AWS_PROFILE=staging-crib
aws eks update-kubeconfig --name main-stage-cluster --alias main-stage-cluster-crib --profile staging-crib

# Check if the Docker daemon is running
if docker info > /dev/null 2>&1; then
  echo "Docker daemon is running, authorizing registry"
else
  echo "Docker daemon is not running, exiting"
  exit 1
fi
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${registry_id}.dkr.ecr.us-west-2.amazonaws.com
devspace use namespace $1
