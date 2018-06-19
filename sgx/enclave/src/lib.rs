#![crate_name = "enclave"]
#![crate_type = "staticlib"]

#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

extern crate sgx_types;
#[macro_use]
extern crate sgx_tstd as std;

use sgx_types::*;
use std::string::String;
use std::slice;

#[no_mangle]
pub extern "C" fn sgx_http_get(url_ptr: *const u8, url_len: usize) -> sgx_status_t {
    let url_slice = unsafe { slice::from_raw_parts(url_ptr, url_len) };
    let url = String::from_utf8(url_slice.to_vec()).unwrap();

    println!("Performing HTTP GET from within enclave with {:?}", url);
    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn sgx_http_post(url_ptr: *const u8, url_len: usize, body_ptr: *const u8, body_len: usize) -> sgx_status_t {
    let url_slice = unsafe { slice::from_raw_parts(url_ptr, url_len) };
    let url = String::from_utf8(url_slice.to_vec()).unwrap();

    let body_slice = unsafe { slice::from_raw_parts(body_ptr, body_len) };
    let body = String::from_utf8(body_slice.to_vec()).unwrap();

    println!("Performing HTTP POST from within enclave with {:?}: {:?}", url, body);
    sgx_status_t::SGX_SUCCESS
}
