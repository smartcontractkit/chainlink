#!/bin/bash
# This is a script to run integration tests on a Go project

# Step into the directory
cd core/chains/evm/headtracker

# Export the database URL
export CL_DATABASE_URL=postgresql://127.0.0.1:5432/chainlink_test?sslmode=disable

# Run the tests with the specified flags
go test -v  ./head_saver_in_mem_test.go  -race -tags integration -count=10 | tee test.txt
