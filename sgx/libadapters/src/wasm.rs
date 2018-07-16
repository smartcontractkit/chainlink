use libc;
use sgx_types::*;
use std::ffi::CStr;
use std::ptr;

use ENCLAVE;

extern "C" {
    fn sgx_wasm(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        wasmt: *const u8,
        wasmt_len: usize,
    ) -> sgx_status_t;
}

fn cstr_len(string: *const libc::c_char) -> usize {
    let buffer = unsafe { CStr::from_ptr(string).to_bytes() };
    buffer.len()
}

#[no_mangle]
pub extern "C" fn wasm(wasmt: *const libc::c_char) -> *const libc::c_char {
    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_wasm(
            ENCLAVE.geteid(),
            &mut retval,
            wasmt as *const u8,
            cstr_len(wasmt),
        )
    };
    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            println!("Call into Enclave wasm failed: {}", result.as_str());
        }
    }
    ptr::null()
}
