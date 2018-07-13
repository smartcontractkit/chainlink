use libc;
use sgx_types::*;
use std::ffi::CStr;
use std::ptr;

use ENCLAVE;

extern "C" {
    fn sgx_multiply(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        multiplicand: *const u8,
        multiplicand_len: usize,
        multipliplier: *const u8,
        multipliplier_len: usize,
    ) -> sgx_status_t;
}

fn cstr_len(string: *const libc::c_char) -> usize {
    let buffer = unsafe { CStr::from_ptr(string).to_bytes() };
    buffer.len()
}

#[no_mangle]
pub extern "C" fn multiply(
    multiplicand: *const libc::c_char,
    multiplier: *const libc::c_char
    ) -> *const libc::c_char {

    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_multiply(
            ENCLAVE.geteid(),
            &mut retval,
            multiplicand as *const u8,
            cstr_len(multiplicand),
            multiplier as *const u8,
            cstr_len(multiplier),
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
