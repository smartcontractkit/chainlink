#!/bin/bash

config_loaded=0

read_input() {
    local prompt=$1
    local allow_empty=$2
    local input_type=$3
    local input

    while true; do
        read -p "$prompt" input

        if [[ "$input_type" == "number" ]]; then
            if [[ $input =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
                break
            else
                echo "Please enter a valid number."
                continue
            fi
        fi

        if [[ -n "$input" || "$allow_empty" == "true" ]]; then
            break
        else
            echo "This is a required field. Please enter a value."
        fi
    done

    echo "$input"
}

select_network_and_set_urls() {
    echo "Select a network to run the test on:"
    select network in SIMULATED SIMULATED_1 SIMULATED_2 SIMULATED_BESU_NONDEV_1 SIMULATED_BESU_NONDEV_2 SIMULATED_NONDEV ETHEREUM_MAINNET GOERLI SEPOLIA KLAYTN_MAINNET KLAYTN_BAOBAB METIS_ANDROMEDA METIS_STARDUST ARBITRUM_MAINNET ARBITRUM_GOERLI ARBITRUM_SEPOLIA OPTIMISM_MAINNET OPTIMISM_GOERLI OPTIMISM_SEPOLIA BASE_GOERLI BASE_SEPOLIA CELO_ALFAJORES CELO_MAINNET RSK POLYGON_MUMBAI POLYGON_MAINNET AVALANCHE_FUJI AVALANCHE_MAINNET QUORUM SCROLL_SEPOLIA SCROLL_MAINNET BASE_MAINNET BSC_TESTNET BSC_MAINNET LINEA_GOERLI LINEA_MAINNET POLYGON_ZKEVM_GOERLI POLYGON_ZKEVM_MAINNET FANTOM_TESTNET FANTOM_MAINNET WEMIX_TESTNET WEMIX_MAINNET KROMA_SEPOLIA KROMA_MAINNET ZK_SYNC_GOERLI ZK_SYNC_MAINNET; do
        if [[ -n "$network" ]]; then
            SELECTED_NETWORKS="$network"
            break
        else
            echo "Invalid selection. Please try again."
        fi
    done

    http_url=$(read_input "Enter HTTP RPC url for $SELECTED_NETWORKS: " false)
    ws_url=$(read_input "Enter WS RPC url for $SELECTED_NETWORKS: " false)
}

# Check for an existing config
load_previous_config() {
    local config_file="/tmp/${product}_config.env"
    if [ -f "$config_file" ]; then
        read -p "Previous configuration for $product found. Load it? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Loading previous configuration..."
            source "$config_file"
            config_loaded=1
            return 0
        fi
    fi
    return 1
}

# Save config
save_current_config() {
    local config_file="/tmp/${product}_config.env"
    cat << EOF > "$config_file"
export SELECTED_NETWORKS="$SELECTED_NETWORKS"
export ${SELECTED_NETWORKS}_HTTP_URLS="$http_url"
export ${SELECTED_NETWORKS}_URLS="$ws_url"
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

# Run tests
run_automation() {
    echo "Running command with provided environment variables..."
    go test -v -test.run "$SUITE" ./automation_test.go
}

# Catch errors
trap 'echo "An error occurred. Exiting..."; exit 1;' ERR

PS3="Please select a product: "
options=("Automation" "Exit")
select product in "${options[@]}"; do
    case $product in
        "Automation")
            if ! load_previous_config; then
                user=$(read_input "Please enter a user that the test will run under: " false)
                log_level=$(read_input "Select test log level (INFO, DEBUG, ERROR): " false)
                select_network_and_set_urls
                private_key=$(read_input "Enter the private key the test should use: " false)
                docker_image_repo=$(read_input "Enter docker image repository for core: " false)
                funding_per_node=$(read_number "Enter funding per node: " false, "number")
                docker_image_version=$(read_input "Enter docker image version for core: " false)
                additional_args=$(read_input "Enter the test suite: " false)
                save_current_config
            fi
            run_automation
            break
            ;;
        "Exit")
            echo "Exiting the wizard."
            exit 0
            ;;
        *) echo "Invalid option $REPLY";;
    esac
done
