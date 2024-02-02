#!/usr/bin/env bash

cat << EOF >> ~/.aws/config
[profile staging-crib]
region=us-west-2
sso_start_url=https://smartcontract.awsapps.com/start
sso_region=us-west-2
sso_account_id=323150190480
sso_role_name=CRIB-ECR-Power
EOF

aws sso login --profile staging-crib
export AWS_PROFILE=staging-crib
aws eks update-kubeconfig --name main-stage-cluster --alias main-stage-cluster-crib --profile staging-crib
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 323150190480.dkr.ecr.us-west-2.amazonaws.com
export DEVSPACE_IMAGE="323150190480.dkr.ecr.us-west-2.amazonaws.com/chainlink-devspace"
devspace use namespace $1