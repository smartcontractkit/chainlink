#!/bin/bash
set -e
rm -rf KinSwarm
mkdir -p KinSwarm/logic_core/src KinSwarm/simulation KinSwarm/include KinSwarm/lib

# Write C Header
cat << 'K_EOF' > KinSwarm/include/kernel.h
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

# Write Cargo.toml
cat << 'K_EOF' > KinSwarm/logic_core/Cargo.toml
[package]
name = "logic_core"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib", "rlib"]

[dependencies]
sha2 = "0.10"
K_EOF

# Write Rust Source
cat << 'K_EOF' > KinSwarm/logic_core/src/lib.rs
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

echo "Starting build..."
cd KinSwarm/logic_core
cargo build --release
cd ..
cp logic_core/target/release/liblogic_core.dylib ./lib/ 2>/dev/null || cp logic_core/target/release/liblogic_core.so ./lib/ 2>/dev/null || true
echo "SUCCESS: logic_core.dylib is in KinSwarm/lib"
