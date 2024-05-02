#!/usr/bin/env bash

set -euo pipefail

# Update helm dependencies and builds charts for sub-charts.

# Function to add Helm repository if it does not exist

repos_added=0

add_helm_repo() {
    local repo_name="$1"
    local repo_url="$2"
    if ! helm repo list | grep -q "$repo_url"; then
        echo "Adding missing Helm repository: $repo_name"
        helm repo add "$repo_name" "$repo_url"
        repos_added=1
    else
        echo "Repository $repo_name already exists."
    fi
}

# Add required repositories
add_helm_repo grafana https://grafana.github.io/helm-charts
add_helm_repo mock-server https://www.mock-server.com
add_helm_repo opentelemetry-collector https://open-telemetry.github.io/opentelemetry-helm-charts
add_helm_repo tempo https://grafana.github.io/helm-charts

# Update repositories to make sure we have the latest versions of charts
if [[ "${repos_added}" -eq 1 ]]; then
    helm repo update
fi

charts_path="../charts"
local_charts=(chainlink-cluster)

for chart in "${local_charts[@]}"; do
  echo "Building chart for $chart from $charts_path/$chart/Chart.lock"
  helm dependency build "$charts_path/$chart"
done
