#!/usr/bin/env bash

set -euo pipefail

###
# To be invoked by `devspace` after a successful DevSpace deploy via a hook.
###

if [[ -z "${DEVSPACE_HOOK_KUBE_NAMESPACE:-}" ]]; then
  echo "Error: DEVSPACE_HOOK_KUBE_NAMESPACE is not set. Make sure to run from devspace."
  exit 1
fi

INGRESS_NAME="${1:-}"
if [[ -z "${INGRESS_NAME}" ]]; then
    echo "Usage: $0 INGRESS_NAME"
    exit 1
fi

max_retries=10
sleep_duration_retry=10 # 10 seconds
sleep_duration_propagate=60 # 60 seconds
timeout=$((60 * 2)) # 2 minutes
elapsed=0 # Track the elapsed time

# Loop until conditions are met or we reach max retries or timeout
for ((i=1; i<=max_retries && elapsed<=timeout; i++)); do
  ingress_hostname_aws=$(kubectl get ingress "${INGRESS_NAME}" -n "${DEVSPACE_HOOK_KUBE_NAMESPACE}" \
    -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
  
  # Sometimes the events on the ingress are "<none>" instead of "successfully reconciled".
  # So we use the AWS hostname as a signal that the ingress has been created successfully.
  if echo "${ingress_hostname_aws}" | grep -q ".elb.amazonaws.com"; then
    echo "#############################################################"
    echo "# Ingress hostnames:"
    echo "#############################################################"
    devspace run ingress-hosts
    echo
    echo "Sleeping for ${sleep_duration_propagate} seconds to allow DNS records to propagate... (Use CTRL+C to safely skip this step.)"
    sleep $sleep_duration_propagate
    echo "...done. NOTE: If you have an issue with the DNS records, try to reset your local and/or VPN DNS cache."
    exit 0
  else
    echo "Attempt $i: Waiting for the ingress to be created..."
    sleep $sleep_duration_retry
    ((elapsed += sleep_duration_retry))
  fi
done

# If we reached here, it means we hit the retry limit or the timeout
echo "Error: Ingress was not successfully created within the given constraints."
exit 1
