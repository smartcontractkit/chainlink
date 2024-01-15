#!/bin/bash
decoded_toml=$(echo "$BASE64_CONFIG_OVERRIDE" | base64 -d)

CHAINLINK_VERSION=$(echo "$decoded_toml" | awk -F'=' '/^[[:space:]]*version[[:space:]]*=/ {gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2); print $2}' 2>/dev/null)
CHAINLINK_IMAGE=$(echo "$decoded_toml" | awk -F'=' '/^[[:space:]]*image[[:space:]]*=/ {gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2); print $2}' 2>/dev/null)
NETWORKS=$(echo "$decoded_toml" | awk -F'=' '/^[[:space:]]*selected_networks[[:space:]]*=/ {gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2); print $2}' 2>/dev/null)
ETH2_EL_CLIENT=$(echo "$decoded_toml" | awk -F'=' '/^[[:space:]]*execution_layer[[:space:]]*=/ {gsub(/^[[:space:]]+|[[:space:]]+$/, "", $2); print $2}' 2>/dev/null)

::add-mask:: "$CHAINLINK_IMAGE"
echo "VERSION: $CHAINLINK_VERSION"
echo "NETWORKS: $NETWORKS"
echo "ETH2_EL_CLIENT: $ETH2_EL_CLIENT"