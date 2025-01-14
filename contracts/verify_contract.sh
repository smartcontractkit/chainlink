#!/bin/bash

# Default values, you will change these
ETHERSCAN_API_KEY="H6D2GJ5CY33THT4SDWWRWVYQ7QC9H6M45Q"
VERIFIER_URL="https://api-sepolia.arbiscan.io/api"
COMPILER_VERSION="0.8.19+commit.7dd6d404" #Do not change
NUM_OF_OPTIMIZATIONS=1000000 #Do not change
DEPLOYMENT_INFO_FILE="deployment-info.json" #Do not change
VERIFIER="etherscan" 
CHAIN_ID=421614

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
    cmd+=" --skip-is-verified-check"
    cmd+=" --verifier $VERIFIER"
    cmd+=" --chain-id $CHAIN_ID"


    if [ -n "$constructor_args" ]; then
        cmd+=" --constructor-args $constructor_args"
    fi

    cmd+=" --num-of-optimizations $NUM_OF_OPTIMIZATIONS"
    cmd+=" --watch"
    # cmd+=" --show-standard-json-input"
    echo "Verifying $contract_name..."
    echo "$cmd"

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

# Verify FeeManager if present
if [ "$(get_json_value '.contracts.FeeManager')" != "" ]; then
    address=$(get_json_value '.contracts.FeeManager.address')
    linkToken=$(get_json_value '.contracts.FeeManager.params.linkToken')
    nativeToken=$(get_json_value '.contracts.FeeManager.params.nativeToken')
    verifierProxy=$(get_json_value '.contracts.FeeManager.params.verifierProxy')
    rewardManager=$(get_json_value '.contracts.FeeManager.params.rewardManager')

    if [ -n "$address" ] && [ -n "$linkToken" ] && [ -n "$nativeToken" ] && [ -n "$verifierProxy" ] && [ -n "$rewardManager" ]; then
        constructor_args=$(cast abi-encode "constructor(address,address,address,address)" "$linkToken" "$nativeToken" "$verifierProxy" "$rewardManager")
        echo "verify_contract \"FeeManager\" \"$address\" \"$constructor_args\" \"src/v0.8/llo-feeds/v0.5.0/FeeManager.sol\""
        verify_contract "FeeManager" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.5.0/FeeManager.sol"
    else
        echo "Error: Missing parameters for FeeManager"
        exit 1
    fi
else
    echo "FeeManager not found in deployment info. Skipping."
fi

# Verify RewardManager if present
if [ "$(get_json_value '.contracts.RewardManager')" != "" ]; then
    address=$(get_json_value '.contracts.RewardManager.address')
    linkToken=$(get_json_value '.contracts.RewardManager.params.linkToken')

    if [ -n "$address" ] && [ -n "$linkToken" ]; then
        constructor_args=$(cast abi-encode "constructor(address)" "$linkToken")
        verify_contract "RewardManager" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.5.0/RewardManager.sol"
    else
        echo "Error: Missing parameters for RewardManager"
        exit 1
    fi
else
    echo "RewardManager not found in deployment info. Skipping."
fi

# Verify Verifier if present
if [ "$(get_json_value '.contracts.Verifier')" != "" ]; then
    address=$(get_json_value '.contracts.Verifier.address')
    verifierProxy=$(get_json_value '.contracts.Verifier.params.verifierProxy')

    if [ -n "$address" ] && [ -n "$verifierProxy" ]; then
        constructor_args=$(cast abi-encode "constructor(address)" "$verifierProxy")
        verify_contract "Verifier" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.5.0/Verifier.sol"
    else
        echo "Error: Missing parameters for Verifier"
        exit 1
    fi
else
    echo "Verifier not found in deployment info. Skipping."
fi

# Verify VerifierProxy if present
if [ "$(get_json_value '.contracts.VerifierProxy')" != "" ]; then
    address=$(get_json_value '.contracts.VerifierProxy.address')
    accessController=$(get_json_value '.contracts.VerifierProxy.params.accessController')

    if [ -n "$address" ] && [ -n "$accessController" ]; then
        constructor_args=$(cast abi-encode "constructor(address)" "$accessController")
        verify_contract "VerifierProxy" "$address" "$constructor_args" "src/v0.8/llo-feeds/v0.5.0/VerifierProxy.sol"
    else
        echo "Error: Missing parameters for VerifierProxy"
        exit 1
    fi
else
    echo "VerifierProxy not found in deployment info. Skipping."
fi

echo "All present contracts verified successfully."
