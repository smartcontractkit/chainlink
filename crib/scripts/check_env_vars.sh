#!/bin/bash

# List of required environment variables
required_vars=(
  "DEVSPACE_IMAGE"
  "DEVSPACE_PROFILE"
  "DEVSPACE_INGRESS_CIDRS"
  "DEVSPACE_INGRESS_BASE_DOMAIN"
  "DEVSPACE_INGRESS_CERT_ARN"
  "DEVSPACE_K8S_POD_WAIT_TIMEOUT"
  "CHAINLINK_CLUSTER_HELM_CHART_URI"
  "NS_TTL"
)

missing_vars=0 # Counter for missing variables

# Check each variable
for var in "${required_vars[@]}"; do
  if [ -z "${!var}" ]; then # If variable is unset or empty
    echo "Error: Environment variable $var is not set."
    missing_vars=$((missing_vars + 1))
  fi
done

# Check for keystone specific profiles
if [[ "${DEVSPACE_PROFILE}" == "keystone" ]]; then
  keystone_vars=(
    "KEYSTONE_ETH_WS_URL"
    "KEYSTONE_ETH_HTTP_URL"
    "KEYSTONE_ACCOUNT_KEY"
  )

  for var in "${keystone_vars[@]}"; do
    if [ -z "${!var}" ]; then # If variable is unset or empty
      echo "Error: Environment variable $var is not set."
      missing_vars=$((missing_vars + 1))
    fi
  done
fi

# Exit with an error if any variables were missing
if [ $missing_vars -ne 0 ]; then
  echo "Total missing environment variables: $missing_vars"
  echo "To fix it, add missing variables in the \"crib/.env\" file."
  echo "you can find example of the .env config in the \"crib/.env.example\""
  exit 1
else
  echo "All required environment variables are set."
fi
