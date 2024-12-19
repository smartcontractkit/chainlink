#!/bin/bash

# Default values
ETHERSCAN_API_KEY="xxx"
VERIFIER_URL="https://api-sepolia.blastscan.io/api"
COMPILER_VERSION="0.8.19+commit.7dd6d404"
NUM_OF_OPTIMIZATIONS=1000000
DEPLOYMENT_INFO_FILE="deployment-info.json"

# Function to parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --etherscan-api-key)
                ETHERSCAN_API_KEY="$2"
                shift 2
                ;;
            --verifier-url)
                VERIFIER_URL="$2"
                shift 2
                ;;
            --compiler-version)
                COMPILER_VERSION="$2"
                shift 2
                ;;
            --num-of-optimizations)
                NUM_OF_OPTIMIZATIONS="$2"
                shift 2
                ;;
            --deployment-info)
                DEPLOYMENT_INFO_FILE="$2"
                shift 2
                ;;
            *)
                echo "Unknown option: $1"
                exit 1
                ;;
        esac
    done
}

# Function to verify a contract
verify_contract() {
    local contract_name="$1"
    local address="$2"
    local constructor_args="$3"
    local file_path="$4"

    if [ -z "$address" ] || [ "$address" == "null" ]; then
        echo "Error: Invalid address for $contract_name"
        exit 1
    fi

    local cmd="forge verify-contract"
    cmd+=" --etherscan-api-key $ETHERSCAN_API_KEY"
    cmd+=" --verifier-url $VERIFIER_URL"
    cmd+=" $address"
    cmd+=" $file_path:$contract_name"
    cmd+=" --compiler-version $COMPILER_VERSION"
    
    if [ -n "$constructor_args" ]; then
        cmd+=" --constructor-args $constructor_args"
    fi
    
    cmd+=" --num-of-optimizations $NUM_OF_OPTIMIZATIONS"

    echo "Verifying $contract_name..."
    
    # Execute the command
    if ! eval "$cmd"; then
        echo "Verification failed for $contract_name"
        exit 1
    fi

    echo "Verification successful for $contract_name"
}

# Function to safely get JSON value
get_json_value() {
    local value=$(echo "$DEPLOYMENT_INFO" | jq -r "$1")
    if [ "$value" == "null" ] || [ -z "$value" ]; then
        echo ""
    else
        echo "$value"
    fi
}

# Main execution
parse_args "$@"

# Check if the deployment info file exists
if [ ! -f "$DEPLOYMENT_INFO_FILE" ]; then
    echo "Error: Deployment info file '$DEPLOYMENT_INFO_FILE' not found."
    exit 1
fi

# Read deployment info
DEPLOYMENT_INFO=$(cat "$DEPLOYMENT_INFO_FILE")

# Verify DestinationFeeManager if present
if [ "$(get_json_value '.contracts.DestinationFeeManager')" != "" ]; then
    address=$(get_json_value '.contracts.DestinationFeeManager.address')
    linkToken=$(get_json_value '.contracts.DestinationFeeManager.params.linkToken')
    nativeToken=$(get_json_value '.contracts.DestinationFeeManager.params.nativeToken')
    verifier=$(get_json_value '.contracts.DestinationFeeManager.params.verifier')
    rewardManager=$(get_json_value '.contracts.DestinationFeeManager.params.rewardManager')

    if [ -n "$address" ] && [ -n "$linkToken" ] && [ -n "$nativeToken" ] && [ -n "$verifier" ] && [ -n "$rewardManager" ]; then
        constructor_args=$(cast abi-encode "constructor(address,address,address,address)" "$linkToken" "$nativeToken" "$verifier" "$rewardManager")
        verify_contract "DestinationFeeManager" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.4.0/DestinationFeeManager.sol"
    else
        echo "Error: Missing parameters for DestinationFeeManager"
        exit 1
    fi
else
    echo "DestinationFeeManager not found in deployment info. Skipping."
fi

# Verify DestinationRewardManager if present
if [ "$(get_json_value '.contracts.DestinationRewardManager')" != "" ]; then
    address=$(get_json_value '.contracts.DestinationRewardManager.address')
    linkToken=$(get_json_value '.contracts.DestinationRewardManager.params.linkToken')

    if [ -n "$address" ] && [ -n "$linkToken" ]; then
        constructor_args=$(cast abi-encode "constructor(address)" "$linkToken")
        verify_contract "DestinationRewardManager" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.4.0/DestinationRewardManager.sol"
    else
        echo "Error: Missing parameters for DestinationRewardManager"
        exit 1
    fi
else
    echo "DestinationRewardManager not found in deployment info. Skipping."
fi

# Verify DestinationVerifier if present
if [ "$(get_json_value '.contracts.DestinationVerifier')" != "" ]; then
    address=$(get_json_value '.contracts.DestinationVerifier.address')
    destinationVerifierProxy=$(get_json_value '.contracts.DestinationVerifier.params.destinationVerifierProxy')

    if [ -n "$address" ] && [ -n "$destinationVerifierProxy" ]; then
        constructor_args=$(cast abi-encode "constructor(address)" "$destinationVerifierProxy")
        verify_contract "DestinationVerifier" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.4.0/DestinationVerifier.sol"
    else
        echo "Error: Missing parameters for DestinationVerifier"
        exit 1
    fi
else
    echo "DestinationVerifier not found in deployment info. Skipping."
fi

# Verify DestinationVerifierProxy if present
if [ "$(get_json_value '.contracts.DestinationVerifierProxy')" != "" ]; then
    address=$(get_json_value '.contracts.DestinationVerifierProxy.address')

    if [ -n "$address" ]; then
        verify_contract "DestinationVerifierProxy" "$address" "" "src/v0.8/llo-feeds/v0.4.0/DestinationVerifierProxy.sol"
    else
        echo "Error: Missing address for DestinationVerifierProxy"
        exit 1
    fi
else
    echo "DestinationVerifierProxy not found in deployment info. Skipping."
fi

echo "All present contracts verified successfully."
