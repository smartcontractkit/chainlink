#!/bin/bash

# Get working dir of script
SCRIPT=$(readlink -f "$0")
TEST_PATH=$(dirname "$SCRIPT")

# Init
SELECTED_NETWORKS=""
http_url=""
ws_url=""
user=""
log_level="DEBUG" # Default value
private_key=""
docker_image_repo=""
funding_per_node="0.1" # Default value
docker_image_version=""
additional_args=""
product=""
list_networks=0
load_config=""
show_help=0

# Help
usage() {
    echo "Usage: $0 [--list-networks] [--load-config file] [--selected-networks SELECTED_NETWORKS] [--http-url http_url] [--ws-url ws_url] [--user user] [--log-level log_level] [--private-key private_key] [--docker-image-repo docker_image_repo] [--funding-per-node funding_per_node] [--docker-image-version docker_image_version] [--test-suite test_suite] [--product product]"
    echo
    echo "Options:"
    echo "  --list-networks         List all available networks"
    echo "  --load-config           Load configuration from a file"
    echo "  --selected-networks     Specify the selected networks"
    echo "  --http-url              Specify HTTP RPC url"
    echo "  --ws-url                Specify WS RPC url"
    echo "  --user                  Specify the user"
    echo "  --log-level             Specify the log level (INFO, DEBUG, ERROR). Default: DEBUG"
    echo "  --private-key           Specify the private key"
    echo "  --docker-image-repo     Specify docker image repository"
    echo "  --funding-per-node      Specify funding per node. Default: 0.1"
    echo "  --docker-image-version  Specify docker image version"
    echo "  --test-suite            Specify additional test suite args"
    echo "  --product               Specify the product (Automation)"
    exit 0
}

list_networks() {
    echo "Available networks:"
    networks="SIMULATED SIMULATED_1 SIMULATED_2 SIMULATED_BESU_NONDEV_1 SIMULATED_BESU_NONDEV_2 SIMULATED_NONDEV ETHEREUM_MAINNET GOERLI SEPOLIA KLAYTN_MAINNET KLAYTN_BAOBAB METIS_ANDROMEDA METIS_STARDUST ARBITRUM_MAINNET ARBITRUM_GOERLI ARBITRUM_SEPOLIA OPTIMISM_MAINNET OPTIMISM_GOERLI OPTIMISM_SEPOLIA BASE_GOERLI BASE_SEPOLIA CELO_ALFAJORES CELO_MAINNET RSK POLYGON_MUMBAI POLYGON_MAINNET AVALANCHE_FUJI AVALANCHE_MAINNET QUORUM SCROLL_SEPOLIA SCROLL_MAINNET BASE_MAINNET BSC_TESTNET BSC_MAINNET LINEA_GOERLI LINEA_MAINNET POLYGON_ZKEVM_GOERLI POLYGON_ZKEVM_MAINNET FANTOM_TESTNET FANTOM_MAINNET WEMIX_TESTNET WEMIX_MAINNET KROMA_SEPOLIA KROMA_MAINNET ZK_SYNC_GOERLI ZK_SYNC_MAINNET"
    read -r -a network_array <<< "$networks"
    columns=4
    for (( i=0; i<${#network_array[@]}; i++ )); do
        # Print network name
        printf '%-30s' "${network_array[i]}"

        # New line every 'columns' networks
        if (( (i + 1) % columns == 0 )); then
            echo
        fi
    done
    echo
    exit 0
}

# Validate config
load_and_validate_config() {
    if [ -f "$load_config" ]; then
        source "$load_config"

        # Validate mandatory variables
        if [ -z "$SELECTED_NETWORKS" ]; then
            echo "Error: SELECTED_NETWORKS variable is missing in the configuration file."
            exit 1
        fi
    else
        echo "Error: Configuration file '$load_config' not found."
        exit 1
    fi
}

# Parse command-line options
while [[ $# -gt 0 ]]; do
    case "$1" in
        --list-networks ) list_networks=1; shift ;;
        --help ) show_help=1; shift ;;
        --load-config ) load_config="$2"; shift 2 ;;
        --selected-networks ) SELECTED_NETWORKS="$2"; shift 2 ;;
        --http-url ) http_url="$2"; shift 2 ;;
        --ws-url ) ws_url="$2"; shift 2 ;;
        --user ) user="$2"; shift 2 ;;
        --log-level ) log_level="$2"; shift 2 ;;
        --private-key ) private_key="$2"; shift 2 ;;
        --docker-image-repo ) docker_image_repo="$2"; shift 2 ;;
        --funding-per-node ) funding_per_node="$2"; shift 2 ;;
        --docker-image-version ) docker_image_version="$2"; shift 2 ;;
        --test-suite ) additional_args="$2"; shift 2 ;;
        --product ) product="$2"; shift 2 ;;
        --) shift; break ;;
        *) break ;;
    esac
done

if [[ "$list_networks" -eq 1 ]]; then
    list_networks
fi

# Run tests
run_automation() {
    echo "Running command with provided environment variables..."
    go test -v -test.run "$SUITE" "${TEST_PATH}/automation_test.go"
}

# Load config
if [[ -n "$load_config" ]]; then
    load_and_validate_config
    run_automation
    exit 0
fi

if [[ "$show_help" -eq 1 ]]; then
    usage
    exit 0
fi

# Check arguments
if [[ -z "$SELECTED_NETWORKS" || -z "$http_url" || -z "$ws_url" || -z "$user" || -z "$private_key" || -z "$docker_image_repo" || -z "$docker_image_version" || -z "$additional_args" || -z "$product" ]]; then
    echo "All options are mandatory."
    usage
    exit 1
fi

# Save config
save_current_config() {
    local config_file="/tmp/${product}_config.env"
    cat << EOF > "$config_file"
export SELECTED_NETWORKS="$SELECTED_NETWORKS"
export ${SELECTED_NETWORKS}_HTTP_URLS="$http_url"
export ${SELECTED_NETWORKS}_WS_URLS="$ws_url"
export CHAINLINK_ENV_USER="$user"
export TEST_LOG_LEVEL="$log_level"
export ${SELECTED_NETWORKS}_KEYS="$private_key"
export CHAINLINK_IMAGE="$docker_image_repo"
export CHAINLINK_NODE_FUNDING="$funding_per_node"
export CHAINLINK_VERSION="$docker_image_version"
export product="$product"
export SUITE="$additional_args"
EOF
    source "$config_file"
}

trap 'echo "An error occurred. Exiting..."; exit 1;' ERR

# Main execution
save_current_config
run_automation
