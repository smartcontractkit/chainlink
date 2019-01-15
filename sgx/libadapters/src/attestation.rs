use errno::{set_errno, Errno};
use libc;
use sgx_types::*;

use use_enclave;

extern "C" {
    fn sgx_report(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        result_ptr: *mut u8,
        result_capacity: usize,
        result_len: *mut usize,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn report(
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    use_enclave(|enclave_id| {
        let mut retval = sgx_status_t::SGX_SUCCESS;
        let result = unsafe {
            sgx_report(
                enclave_id,
                &mut retval,
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
