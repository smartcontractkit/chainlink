#!/bin/bash
set -e

# 1. Environment and Path Normalization
PROJECT_ROOT=$(pwd)
BASE_DIR="$PROJECT_ROOT/KinSwarm"
mkdir -p "$BASE_DIR/logic_core/src" "$BASE_DIR/include" "$BASE_DIR/lib"

# 2. Write Rust Infrastructure
cat << 'K_EOF' > "$BASE_DIR/logic_core/Cargo.toml"
[package]
name = "logic_core"
version = "0.1.0"
edition = "2021"

[lib]
name = "logic_core"
crate-type = ["cdylib"]

[dependencies]
sha2 = "0.10"
K_EOF

cat << 'K_EOF' > "$BASE_DIR/logic_core/src/lib.rs"
use sha2::{Sha256, Digest};

#[repr(C)]
pub struct SettlementOutcome {
    pub root_output: [u8; 32],
    pub success: bool,
}

#[no_mangle]
pub extern "C" fn execute_settlement_anchor(root_input: *const u8, amount: u64, worker_count: u32) -> SettlementOutcome {
    let mut outcome = SettlementOutcome { root_output: [0u8; 32], success: true };
    if root_input.is_null() { outcome.success = false; return outcome; }
    let input_slice = unsafe { std::slice::from_raw_parts(root_input, 32) };
    let mut hasher = Sha256::new();
    hasher.update(input_slice);
    hasher.update(amount.to_le_bytes());
    hasher.update(worker_count.to_le_bytes());
    let result = hasher.finalize();
    outcome.root_output.copy_from_slice(&result);
    outcome
}
K_EOF

# 3. Write C Header
cat << 'K_EOF' > "$BASE_DIR/include/kernel.h"
#ifndef KERNEL_H
#define KERNEL_H
#include <stdint.h>
#include <stdbool.h>
typedef struct {
    uint8_t root_output[32];
    bool success;
} SettlementOutcome;
SettlementOutcome execute_settlement_anchor(const uint8_t* root_input, uint64_t amount, uint32_t worker_count);
#endif
K_EOF

# 4. Build and Extract Library
echo "Building Logic Core..."
cd "$BASE_DIR/logic_core"
cargo build --release

# Programmatically find the target directory
TARGET_DIR=$(cargo metadata --format-version 1 | python3 -c "import sys, json; print(json.load(sys.stdin)['target_directory'])")
DYLIB_PATH="$TARGET_DIR/release/liblogic_core.dylib"

if [ -f "$DYLIB_PATH" ]; then
    cp "$DYLIB_PATH" "$BASE_DIR/lib/"
    echo "Library anchored: $BASE_DIR/lib/liblogic_core.dylib"
else
    echo "FAILED: Could not locate dylib at $DYLIB_PATH"
    exit 1
fi

# 5. Write Go Orchestrator
cd "$PROJECT_ROOT"
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
	inputHex := "4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b"
	inputBytes, _ := hex.DecodeString(inputHex)

	outcome := C.execute_settlement_anchor(
		(*C.uint8_t)(unsafe.Pointer(&inputBytes[0])),
		C.uint64_t(1000000),
		C.uint32_t(42),
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

# 6. Initialize and Run
if [ ! -f go.mod ]; then go mod init kinswarm/orchestrator; fi
go mod tidy
export DYLD_LIBRARY_PATH="$BASE_DIR/lib:$DYLD_LIBRARY_PATH"

echo "Executing Go harness..."
go run main.go
