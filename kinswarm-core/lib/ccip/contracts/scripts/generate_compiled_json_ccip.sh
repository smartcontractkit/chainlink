#!/usr/bin/env bash

set -e

ROOT="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; cd ../../ && pwd -P )"

contractName="$1"
SOLC_VERSION="$2"
destPath="$3"
path="$ROOT"/contracts/solc/v$SOLC_VERSION/"$contractName"

if [ -z "$contractName" ]; then
  echo "Error: contractName is not set."
  exit 1
fi

if [ -z "$contractName" ];  then
  echo "Error: solc version is not set."
  exit 1
fi

if [ -z "$destPath" ]; then
  echo "Error: destination path is not set."
  exit 1
fi

if [ ! -e "$destPath" ]; then
  echo "Error: $destPath does not exist."
  exit 1
fi


metadata=$(cat "$path/${contractName}_meta.json")
fileName=$(echo "$metadata" | jq -r '.settings.compilationTarget | keys[0]')
combined=$(cat "$path/combined.json")

abi=$(cat "$path/${contractName}.abi")
bytecode_object="0x$(cat "$path/${contractName}.bin" | tr -d '\n')"
bytecode_sourceMap=$(echo "$combined" | jq -r ".contracts[\"$fileName:$contractName\"].srcmap")

deployedBytecode_object="0x$(cat "$path/${contractName}.bin-runtime" | tr -d '\n')"
deployedBytecode_sourceMap=$(echo "$combined" | jq -r ".contracts[\"$fileName:$contractName\"].\"srcmap-runtime\"")

methodIdentifiers=$(echo "$combined" | jq -r ".contracts[\"$fileName:$contractName\"].hashes")
rawMetadata=$(echo "$metadata" | jq -c '.')

result=$(jq -n \
    --argjson abi "$abi" \
    --arg object "$bytecode_object" \
    --arg sourceMap "$bytecode_sourceMap" \
    --arg deployedObject "$deployedBytecode_object" \
    --arg deployedSourceMap "$deployedBytecode_sourceMap" \
    --argjson methodIdentifiers "$methodIdentifiers" \
    --arg rawMetadata "$rawMetadata" \
    --argjson metadata "$metadata" \
    '{
        abi: $abi,
        bytecode: {
            object: $object,
            sourceMap: $sourceMap
        },
        deployedBytecode: {
            object: $deployedObject,
            sourceMap: $deployedSourceMap
        },
        methodIdentifiers: $methodIdentifiers,
        rawMetadata: $rawMetadata,
        metadata: $metadata
    }'
)

echo "$result" > "$destPath"
echo "Generated artifacts at $(realpath $destPath)"
