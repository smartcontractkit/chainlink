use errno::{set_errno, Errno};
use libc;
use sgx_types::*;
use utils::cstr_len;

use use_enclave;

extern "C" {
    fn sgx_wasm(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        adapter: *const u8,
        adapter_len: usize,
        input: *const u8,
        input_len: usize,
        result_ptr: *mut u8,
        result_capacity: usize,
        result_len: *mut usize,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn wasm(
    adapter: *const libc::c_char,
    input: *const libc::c_char,
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    use_enclave(|enclave_id| {
        let mut retval = sgx_status_t::SGX_SUCCESS;
        let result = unsafe {
            sgx_wasm(
                enclave_id,
                &mut retval,
                adapter as *const u8,
                cstr_len(adapter),
                input as *const u8,
                cstr_len(input),
                result_ptr as *mut u8,
                result_capacity,
                result_len as *mut usize,
            )
        };

        if result != sgx_status_t::SGX_SUCCESS {
                set_errno(Errno(result as i32));
                return;
        }

        if retval != sgx_status_t::SGX_SUCCESS {
                set_errno(Errno(retval as i32));
                return;
        }

        set_errno(Errno(0));
    });
}
