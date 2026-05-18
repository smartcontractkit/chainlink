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
