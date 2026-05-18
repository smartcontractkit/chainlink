#!/bin/bash
set -e

# 1. Setup Environment
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"

# 2. Write Multi-Threaded Swarm Coordinator
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
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	NumAgents    = 100
	WorkerAmount = 50000
)

type AgentResult struct {
	ID         int
	AnchorHash string
	Duration   time.Duration
}

func main() {
	fmt.Printf("--- KinSwarm OS: Swarm Coordination Layer ---\n")
	fmt.Printf("Initialising %d agents targeting Logic Core...\n\n", NumAgents)

	var wg sync.WaitGroup
	results := make(chan AgentResult, NumAgents)
	startTime := time.Now()

	// Mock state root for the swarm
	baseRootHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	baseBytes, _ := hex.DecodeString(baseRootHex)

	for i := 1; i <= NumAgents; i++ {
		wg.Add(1)
		go func(agentID int) {
			defer wg.Done()
			
			agentStart := time.Now()
			
			// Execute FFI Call
			outcome := C.execute_settlement_anchor(
				(*C.uint8_t)(unsafe.Pointer(&baseBytes[0])),
				C.uint64_t(WorkerAmount),
				C.uint32_t(agentID),
			)

			if bool(outcome.success) {
				resHash := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
				results <- AgentResult{
					ID:         agentID,
					AnchorHash: hex.EncodeToString(resHash),
					Duration:   time.Since(agentStart),
				}
			}
		}(i)
	}

	// Closer
	go func() {
		wg.Wait()
		close(results)
	}()

	// Processing aggregated results
	successCount := 0
	for res := range results {
		successCount++
		if successCount <= 5 || successCount > 95 {
			fmt.Printf("[Agent %03d] Anchor: %s...%s | Latency: %v\n", 
				res.ID, res.AnchorHash[:8], res.AnchorHash[56:], res.Duration)
		} else if successCount == 6 {
			fmt.Println("...")
		}
	}

	totalDuration := time.Since(startTime)
	fmt.Printf("\n--- Swarm Execution Complete ---\n")
	fmt.Printf("Total Agents:     %d\n", successCount)
	fmt.Printf("Total Time:       %v\n", totalDuration)
	fmt.Printf("Avg Per Agent:    %v\n", totalDuration/time.Duration(NumAgents))
}
K_EOF

# 3. Execution
go mod tidy
echo "Launching KinSwarm Swarm..."
go run main.go
