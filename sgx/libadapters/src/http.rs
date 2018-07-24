use libc;
use sgx_types::*;
use std::ptr;
use utils::cstr_len;

use ENCLAVE;

extern "C" {
    fn sgx_http_get(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        url: *const u8,
        url_len: usize,
    ) -> sgx_status_t;
    fn sgx_http_post(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        url: *const u8,
        url_len: usize,
        body: *const u8,
        body_len: usize,
    ) -> sgx_status_t;
}

#[no_mangle]
pub extern "C" fn http_get(url: *const libc::c_char) -> *const libc::c_char {
    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_http_get(
            ENCLAVE.geteid(),
            &mut retval,
            url as *const u8,
            cstr_len(url),
        )
    };
    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            println!("Call into Enclave sgx_http_get failed: {}", result.as_str());
        }
    }
    ptr::null()
}

#[no_mangle]
pub extern "C" fn http_post(
    url: *const libc::c_char,
    body: *const libc::c_char,
) -> *const libc::c_char {
    let mut retval = sgx_status_t::SGX_SUCCESS;
    let result = unsafe {
        sgx_http_post(
            ENCLAVE.geteid(),
            &mut retval,
            url as *const u8,
            cstr_len(url),
            body as *const u8,
            cstr_len(body),
        )
    };
    match result {
        sgx_status_t::SGX_SUCCESS => {}
        _ => {
            println!(
                "Call into Enclave sgx_http_post failed: {}",
                result.as_str()
            );
        }
    }
    ptr::null()
}
