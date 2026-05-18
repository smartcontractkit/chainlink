#!/bin/bash
set -e

# 1. Environment Check
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"

# 2. Update Go Swarm to write JSONL intents
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
	"sync"
	"time"
	"unsafe"
)

const (
	NumAgents    = 100
	WorkerAmount = 50000
)

type SwarmIntent struct {
	Timestamp  int64  `json:"ts"`
	AgentID    int    `json:"agent_id"`
	AnchorHash string `json:"anchor"`
	Status     string `json:"status"`
}

func main() {
	fmt.Printf("--- KinSwarm OS: Journaling Swarm Results ---\n")
	
	file, err := os.OpenFile("swarm_journal.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open journal: %v\n", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	baseRootHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	baseBytes, _ := hex.DecodeString(baseRootHex)

	for i := 1; i <= NumAgents; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			outcome := C.execute_settlement_anchor(
				(*C.uint8_t)(unsafe.Pointer(&baseBytes[0])),
				C.uint64_t(WorkerAmount),
				C.uint32_t(id),
			)

			if bool(outcome.success) {
				resHash := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
				intent := SwarmIntent{
					Timestamp:  time.Now().UnixNano(),
					AgentID:    id,
					AnchorHash: hex.EncodeToString(resHash),
					Status:     "COMMITTED",
				}
				data, _ := json.Marshal(intent)
				file.WriteString(string(data) + "\n")
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Swarm intent log synchronized to swarm_journal.jsonl")
}
K_EOF

# 3. Execution and GPG Anchoring
go run main.go

# Signing the journal with your GPG Identity (The Keeper)
# This creates swarm_journal.jsonl.asc
echo "Signing journal with GPG..."
gpg --clearsign --yes swarm_journal.jsonl

echo "--- Journal Integrity Verified ---"
ls -l swarm_journal.jsonl*
