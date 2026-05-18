#!/bin/bash
set -e

# 1. Ensure Library Naming and Path Consistency
# The linker expects 'liblogic_core.dylib' for '-llogic_core'
LIB_DIR="$(pwd)/KinSwarm/lib"
INC_DIR="$(pwd)/KinSwarm/include"

if [ ! -f "$LIB_DIR/liblogic_core.dylib" ]; then
    echo "Error: liblogic_core.dylib not found in $LIB_DIR"
    exit 1
fi

# 2. Write Go Source with Absolute Path Linking
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
	"os"
	"unsafe"
)

func main() {
	fmt.Println("--- KinSwarm Sovereign Orchestrator ---")

	// Deterministic 32-byte input
	inputHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	inputBytes, err := hex.DecodeString(inputHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hex Decode Error: %v\n", err)
		os.Exit(1)
	}

	amount := uint64(1000000)
	workerCount := uint32(42)

	// FFI Boundary Cross
	outcome := C.execute_settlement_anchor(
		(*C.uint8_t)(unsafe.Pointer(&inputBytes[0])),
		C.uint64_t(amount),
		C.uint32_t(workerCount),
	)

	if bool(outcome.success) {
		resultHash := C.GoBytes(unsafe.Pointer(&outcome.root_output[0]), 32)
		fmt.Printf("Status: SUCCESS\n")
		fmt.Printf("Anchor Hash: %s\n", hex.EncodeToString(resultHash))
	} else {
		fmt.Printf("Status: FAILURE\n")
		os.Exit(1)
	}
}
K_EOF

# 3. Setup Go Module and Tidy
if [ ! -f go.mod ]; then
    go mod init kinswarm/orchestrator
fi
go mod tidy

# 4. Compile and Run with Runtime Library Path
# We export DYLD_LIBRARY_PATH specifically for the 'go run' execution
export DYLD_LIBRARY_PATH="$LIB_DIR:$DYLD_LIBRARY_PATH"

echo "Building and executing Go harness..."
go run main.go
