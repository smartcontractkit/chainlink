#!/bin/bash
set -e

# 1. Environment and Pathing
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"

# 2. Go Orchestrator for Chainlink External Adapter Logic
cat << 'K_EOF' > main.go
package main

/*
#cgo LDFLAGS: -L${SRCDIR}/KinSwarm/lib -llogic_core
#cgo CFLAGS: -I${SRCDIR}/KinSwarm/include
#include "kernel.h"
*/
import "C"
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"unsafe"
)

type ChainlinkResponse struct {
	JobRunID string `json:"id"`
	Data     struct {
		AnchorHash string `json:"anchor"`
		Status     string `json:"status"`
	} `json:"data"`
}

func main() {
	// Standard Chainlink External Adapter Input Format
	inputHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	inputBytes, _ := hex.DecodeString(inputHex)

	// Execute Rust Kernel for Cryptographic Settlement
	outcome := C.execute_settlement_anchor(
		(*C.uint8_t)(unsafe.Pointer(&inputBytes[0])),
		C.uint64_t(1000000),
		C.uint32_t(777), // Oracle Job ID Mock
	)

	if bool(outcome.success) {
		resHash := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
		anchor := hex.EncodeToString(resHash)

		response := ChainlinkResponse{
			JobRunID: "K_SWARM_001",
		}
		response.Data.AnchorHash = "0x" + anchor
		response.Data.Status = "VERIFIED_OFF_CHAIN"

		output, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(output))
		
		// Write to Chainlink-ready broadcast file
		os.WriteFile("chainlink_payload.json", output, 0644)
	}
}
K_EOF

# 3. Execution and Payload Generation
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"
go mod tidy
go run main.go

echo "--- Chainlink Payload Ready ---"
cat chainlink_payload.json
