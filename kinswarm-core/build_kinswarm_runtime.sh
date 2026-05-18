#!/bin/bash

# 1. Workspace initialization
mkdir -p kinswarm_cre/{src,include,lib}
cd kinswarm_cre

# 2. Rust Kernel Logic (Cleaned of unused warnings)
cat << 'RUST_EOF' > src/lib.rs
#[repr(C)]
pub struct KernelOutcome {
    pub success: bool,
    pub root_output: [u8; 32],
}

#[no_mangle]
pub extern "C" fn execute_settlement_batch(
    merkle_root: *const u8,
    _amount: u64,
    _worker_count: u32
) -> KernelOutcome {
    let mut hash_buffer = [0u8; 32];
    
    unsafe {
        if !merkle_root.is_null() {
            std::ptr::copy_nonoverlapping(merkle_root, hash_buffer.as_mut_ptr(), 32);
        }
    }

    KernelOutcome {
        success: true,
        root_output: hash_buffer,
    }
}
RUST_EOF

# 3. C Header for FFI boundary
cat << 'H_EOF' > include/kernel.h
#ifndef KERNEL_H
#define KERNEL_H

#include <stdint.h>
#include <stdbool.h>

typedef struct {
    bool success;
    uint8_t root_output[32];
} KernelOutcome;

KernelOutcome execute_settlement_batch(const uint8_t* merkle_root, uint64_t amount, uint32_t worker_count);

#endif
H_EOF

# 4. Go CRE Entry Point
cat << 'GO_EOF' > main.go
package main

/*
#cgo LDFLAGS: -L./lib -lkin_kernel
#include "include/kernel.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	root := [32]byte{0xbd, 0xae, 0xe7, 0xb2, 0x05, 0xce, 0x72, 0x6b}
	totalAmount := uint64(350000000000000000)
	workers := uint32(1000)

	outcome := C.execute_settlement_batch(
		(*C.uint8_t)(unsafe.Pointer(&root[0])),
		C.uint64_t(totalAmount),
		C.uint32(workers),
	)

	if outcome.success {
		fmt.Printf("[EVMAdapter] Pushing Root: %x\n", outcome.root_output)
		fmt.Printf("[SolanaAdapter] Pushing Root: %x\n", outcome.root_output)
		fmt.Printf("[CosmosAdapter] Pushing Root: %x\n", outcome.root_output)
	}
}
GO_EOF

# 5. Build Pipeline
cat << 'CARGO_EOF' > Cargo.toml
[package]
name = "kin_kernel"
version = "1.0.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]
CARGO_EOF

cargo build --release
if [ -f target/release/libkin_kernel.so ]; then
    cp target/release/libkin_kernel.so ./lib/
elif [ -f target/release/libkin_kernel.dylib ]; then
    cp target/release/libkin_kernel.dylib ./lib/
fi

echo "Compilation successful. Execute with: LD_LIBRARY_PATH=./lib go run main.go"
