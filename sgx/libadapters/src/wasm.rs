use errno::{set_errno, Errno};
use libc;
use sgx_types::*;
use utils::cstr_len;

use ENCLAVE;

extern "C" {
    fn sgx_wasm(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        wasmt_ptr: *const u8,
        wasmt_len: usize,
        arguments_ptr: *const u8,
        arguments_len: usize,
        result_ptr: *mut u8,
        result_capacity: usize,
        result_len: *mut usize,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn wasm(
    wasm: *const libc::c_char,
    arguments: *const libc::c_char,
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_wasm(
            ENCLAVE.geteid(),
            &mut retval,
            wasm as *const u8,
            cstr_len(wasm),
            arguments as *const u8,
            cstr_len(arguments),
            result_ptr as *mut u8,
            result_capacity,
            result_len as *mut usize,
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {
            println!("Call into enclave succeded");
            set_errno(Errno(0));
        }
        _ => {
            println!("Call into Enclave wasm failed: {}", result.as_str());
            set_errno(Errno(result as i32));
            return;
        }
    }

    match retval {
        sgx_status_t::SGX_SUCCESS => {
            println!("wasm call succeeded");
            set_errno(Errno(0));
        }
        _ => {
            println!("wasm returned error: {}", retval.as_str());
            set_errno(Errno(result as i32));
            return;
        }
    }
}
