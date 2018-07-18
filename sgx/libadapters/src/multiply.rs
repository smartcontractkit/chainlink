use libc;
use sgx_types::*;
use util::cstr_len;

use ENCLAVE;

extern "C" {
    fn sgx_multiply(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        adapter: *const u8,
        adapter_len: usize,
        input: *const u8,
        input_len: usize,
        output: *mut u8,
        output_len: usize,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn multiply(
    adapter: *const libc::c_char,
    input: *const libc::c_char,
    output: *mut libc::c_char,
) {
    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_multiply(
            ENCLAVE.geteid(),
            &mut retval,
            adapter as *const u8,
            cstr_len(adapter),
            input as *const u8,
            cstr_len(input),
            output as *mut u8,
            cstr_len(output),
        )
    };

    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            println!("Call into Enclave multiplier failed: {}", result.as_str());
        }
    }
}
