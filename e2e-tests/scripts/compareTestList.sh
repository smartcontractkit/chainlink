#!/usr/bin/env bash

# accepts a path to a test file to compare the test list against

set -e

# get this scripts directory
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

FILENAME=$1

TESTLIST=$(cat ${FILENAME} | grep "func Test.*\(t \*testing.T\)" | grep -o 'Test[A-Za-z0-9_]*')

# convert the test list from above into json in the form {"tests":[{"name":"TestName"}]}
TESTLISTJSON=$(echo $TESTLIST | jq -R -s -c '{tests: split(" ") | map({"name":.})}')

# Get list of test names from JSON file
JSONFILE="${FILENAME}_test_list.json"
JSONTESTLIST=$(jq -r '.tests[].name' ${JSONFILE})

# Convert lists to arrays
TESTLIST_ARRAY=($(echo "$TESTLIST"))
JSONTESTLIST_ARRAY=($(echo "$JSONTESTLIST"))

ERRORS_FOUND=false

# Compare TESTLIST_ARRAY against JSONTESTLIST_ARRAY
for test in "${TESTLIST_ARRAY[@]}"; do
  if [[ ! " ${JSONTESTLIST_ARRAY[@]} " =~ " ${test} " ]]; then
    echo "$test exists only in ${FILENAME}."
    ERRORS_FOUND=true
  fi
done

# Compare JSONTESTLIST_ARRAY against TESTLIST_ARRAY
for test in "${JSONTESTLIST_ARRAY[@]}"; do
  if [[ ! " ${TESTLIST_ARRAY[@]} " =~ " ${test} " ]]; then
    echo "$test exists only in ${JSONFILE}."
    ERRORS_FOUND=true
  fi
done

if [ "$ERRORS_FOUND" = true ] ; then
  echo "Test lists do not match. Please update ${JSONFILE} with the updated tests to run in CI."
  exit 1
fi
