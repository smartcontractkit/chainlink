#![crate_name = "enclave"]
#![crate_type = "staticlib"]

#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

#[macro_use]
extern crate lazy_static;
extern crate num;
extern crate serde;
extern crate serde_json;
extern crate sgx_types;
#[macro_use]
extern crate sgx_tstd as std;
extern crate sgxwasm;
extern crate wasmi;

mod wasm;

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

#[no_mangle]
pub extern "C" fn sgx_wasm(wasmt_ptr: *const u8, wasmt_len: usize) -> sgx_status_t {
    let wasmt_slice = unsafe { slice::from_raw_parts(wasmt_ptr, wasmt_len) };
    let wasmt = String::from_utf8(wasmt_slice.to_vec()).unwrap();

    let data = &[
        0, 97, 115, 109, // \0ASM - magic
        1, 0, 0, 0       //  0x01 - version
    ];

    match wasmi::Module::from_buffer(data) {
        Ok(r) => {
            //println!("module: {:?}", r);
            return sgx_status_t::SGX_SUCCESS
        },
        Err(err) => {
            println!("Error executing wasm: {:?}", err);
            return sgx_status_t::SGX_ERROR_WASM_INTERPRETER_ERROR;
        },
    };
}
