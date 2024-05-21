#!/usr/bin/env bash

# requires a path to a json file with all the tests it should run
# requires a node label to be passed in, for example "ubuntu-latest"

set -e

# get this script's directory
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

JSONFILE=$1
NODE_LABEL=$2

COUNTER=1

# Build a JSON object in the format expected by our evm-version-compatibility-tests workflow matrix
matrix_output() {
  local counter=$1
  local job_name=$2
  local test_name=$3
  local node_label=$4
  local eth_client=$5
  local docker_image=$6
  local product=$7
  local counter_out=$(printf "%02d\n" $counter)
  echo -n "{\"name\": \"${job_name}-${counter_out}\", \"os\": \"${node_label}\", \"product\": \"${product}\", \"eth_client\": \"${eth_client}\", \"docker_image\": \"${docker_image}\", \"run\": \"-run '^${test_name}$'\"}"
}

# Read the JSON file and loop through 'tests' and 'run'
jq -c '.tests[]' ${JSONFILE} | while read -r test; do
  testName=$(echo ${test} | jq -r '.name')
  label=$(echo ${test} | jq -r '.label // empty')
  effective_node_label=${label:-$NODE_LABEL}
  eth_client=$(echo ${test} | jq -r '.eth_client')
  docker_image=$(echo ${test} | jq -r '.docker_image')
  product=$(echo ${test} | jq -r '.product')
  subTests=$(echo ${test} | jq -r '.run[]?.name // empty')
  output=""

  if [ $COUNTER -ne 1 ]; then
      echo -n ","
  fi

  # Loop through subtests, if any, and print in the desired format
  if [ -n "$subTests" ]; then
    subTestString=""
    subTestCounter=1
    for subTest in $subTests; do
      if [ $subTestCounter -ne 1 ]; then
        subTestString+="|"
      fi
      subTestString+="${testName}\/${subTest}"
      ((subTestCounter++))
    done
    testName="${subTestString}"
  fi
  matrix_output $COUNTER "emv-node-version-compatibility-test" "${testName}" ${effective_node_label} "${eth_client}" "${docker_image}" "${product}"
  ((COUNTER++))
done > "./tmpout.json"
OUTPUT=$(cat ./tmpout.json)
echo "[${OUTPUT}]"
rm ./tmpout.json