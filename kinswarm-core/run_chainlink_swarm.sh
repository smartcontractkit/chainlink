#!/bin/bash
set -e

# 1. Paths
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"

# 2. Go Logic for Hex-Formatted Chainlink Output
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
	"sync"
	"unsafe"
)

type ChainlinkOutput struct {
	JobRunID string `json:"id"`
	Data     map[string]string `json:"data"`
}

func main() {
	var wg sync.WaitGroup
	baseRootHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	baseBytes, _ := hex.DecodeString(baseRootHex)

	// Simulate a batch of 5 high-priority worker settlements
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			outcome := C.execute_settlement_anchor(
				(*C.uint8_t)(unsafe.Pointer(&baseBytes[0])),
				C.uint64_t(1000), // Amount
				C.uint32_t(id),   // Worker ID
			)

			if bool(outcome.success) {
				resHash := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
				
				// Format as 0x-prefixed hex for Chainlink/EVM
				formattedHex := "0x" + hex.EncodeToString(resHash)
				
				resp := ChainlinkOutput{
					JobRunID: fmt.Sprintf("WORKER_%d", id),
					Data: map[string]string{
						"anchor": formattedHex,
						"status": "COMMITTED",
					},
				}
				
				output, _ := json.Marshal(resp)
				
				// Standardized Output Log
				fmt.Printf("[ChainlinkAdapter] Sent: %s\n", formattedHex)
				_ = output // Ready for HTTP response delivery
			}
		}(i)
	}
	wg.Wait()
}
K_EOF

# 3. Execution
go run main.go
