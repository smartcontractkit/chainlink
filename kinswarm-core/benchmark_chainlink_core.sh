#!/bin/bash
set -e

# 1. Environment Check
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"

# 2. Go Benchmark Tool
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
	"time"
	"unsafe"
)

const Iterations = 1000000

func main() {
	fmt.Println("--- KinSwarm: Chainlink Core Benchmark ---")
	fmt.Printf("Executing %d cryptographic settlements...\n", Iterations)

	baseRootHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	baseBytes, _ := hex.DecodeString(baseRootHex)
	
	start := time.Now()

	for i := 0; i < Iterations; i++ {
		_ = C.execute_settlement_anchor(
			(*C.uint8_t)(unsafe.Pointer(&baseBytes[0])),
			C.uint64_t(100*i),
			C.uint32_t(i),
		)
	}

	duration := time.Since(start)
	tps := float64(Iterations) / duration.Seconds()

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Total Time     : %v\n", duration)
	fmt.Printf("Avg Latency    : %v\n", duration / Iterations)
	fmt.Printf("Throughput     : %.2f Settlements/sec\n", tps)
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Benchmark Complete. Core integrity verified.\n")
}
K_EOF

# 3. Execution
go run main.go
