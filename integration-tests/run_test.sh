#!/bin/bash

# Initialize variables
total_tests=20
failures=0
passes=0

# Run tests in a loop
for ((i=1; i<=$total_tests; i++)); do
    echo "Running test $i..."
    go test -v -run $1 $2

    # Check the exit status of the test
    if [ $? -eq 0 ]; then
        ((passes++))
    else
        ((failures++))
    fi
done

# Display results
echo -e "\nTest Results:"
echo "Total Tests: $total_tests"
echo "Passes: $passes"
echo "Failures: $failures"

# Exit with success if there are no failures, otherwise, exit with failure
if [ $failures -eq 0 ]; then
    exit 0
else
    exit 1
fi

