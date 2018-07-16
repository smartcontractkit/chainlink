use libc;
use sgx_types::*;
use std::ffi::CStr;
use std::ptr;

use ENCLAVE;

extern "C" {
    fn sgx_multiply(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        adapter: *const u8,
        adapter_len: usize,
        input: *const u8,
        input_len: usize,
    ) -> sgx_status_t;
}

fn cstr_len(string: *const libc::c_char) -> usize {
    let buffer = unsafe { CStr::from_ptr(string).to_bytes() };
    buffer.len()
}

#[no_mangle]
pub extern "C" fn multiply(
        adapter: *const libc::c_char,
        input: *const libc::c_char,
    ) -> *const libc::c_char {

    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_multiply(
            ENCLAVE.geteid(),
            &mut retval,
            adapter as *const u8,
            cstr_len(adapter),
            input as *const u8,
            cstr_len(input),
        )
    };
    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            println!("Call into Enclave multiplier failed: {}", result.as_str());
        }
    }
    ptr::null()
}
