#![crate_name = "enclave"]
#![crate_type = "staticlib"]
#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

extern crate base64;
extern crate num;
extern crate serde;
extern crate serde_json;
extern crate sgx_types;
#[cfg(not(target_env = "sgx"))]
#[macro_use] extern crate sgx_tstd as std;
extern crate wasmi;

#[macro_use] mod util;
mod wasm;

use sgx_types::*;
use util::{copy_string_to_cstr_ptr, string_from_cstr_with_len};

#[no_mangle]
pub extern "C" fn sgx_http_get(url_ptr: *const u8, url_len: usize) -> sgx_status_t {
    let url = string_from_cstr_with_len(url_ptr, url_len).unwrap();

    println!("Performing HTTP GET from within enclave with {:?}", url);
    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn sgx_http_post(
    url_ptr: *const u8,
    url_len: usize,
    body_ptr: *const u8,
    body_len: usize,
) -> sgx_status_t {
    let url = string_from_cstr_with_len(url_ptr, url_len).unwrap();
    let body = string_from_cstr_with_len(body_ptr, body_len).unwrap();

    println!(
        "Performing HTTP POST from within enclave with {:?}: {:?}",
        url, body
    );
    sgx_status_t::SGX_SUCCESS
}

#[no_mangle]
pub extern "C" fn sgx_wasm(
    wasmt_ptr: *const u8,
    wasmt_len: usize,
    arguments_ptr: *const u8,
    arguments_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> sgx_status_t {
    match wasm(
        wasmt_ptr,
        wasmt_len,
        arguments_ptr,
        arguments_len,
        result_ptr,
        result_capacity,
        result_len,
    ) {
        Ok(_) => sgx_status_t::SGX_SUCCESS,
        _ => sgx_status_t::SGX_ERROR_UNEXPECTED,
    }
}

enum WasmError {
    FromUtf8Error(std::string::FromUtf8Error),
    ExecError(wasm::Error),
    OutputCStrError(util::OutputCStrError),
    UnexpectedOutputError,
}

impl_from_error!(std::string::FromUtf8Error, WasmError::FromUtf8Error);
impl_from_error!(wasm::Error, WasmError::ExecError);
impl_from_error!(util::OutputCStrError, WasmError::OutputCStrError);

fn wasm(
    wasmt_ptr: *const u8,
    wasmt_len: usize,
    arguments_ptr: *const u8,
    arguments_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> Result<(), WasmError> {

    let wasmt = string_from_cstr_with_len(wasmt_ptr, wasmt_len)?;
    println!("wasmt: {:?}", wasmt);

    let arguments = string_from_cstr_with_len(arguments_ptr, arguments_len)?;
    println!("arguments: {:?}", arguments);

    let output = wasm::exec(&wasmt, &arguments)?;
    println!("output: {:?}", output);

    let value = match output {
        wasmi::RuntimeValue::I32(v) => format!("{}", v),
        _ => return Err(WasmError::UnexpectedOutputError),
    };

    copy_string_to_cstr_ptr(&value, result_ptr, result_capacity, result_len)?;
    Ok(())
}

pub extern "C" fn sgx_multiply(multiplicand_ptr: *const u8, multiplicand_len: usize, multiplier_ptr: *const u8, multiplier_len: usize) -> sgx_status_t {
    let multiplicand = string_from_cstr_with_len(multiplicand_ptr, multiplicand_len).unwrap();
    let multiplier = string_from_cstr_with_len(multiplier_ptr, multiplier_len).unwrap();

    println!("Performing MULTIPLY from within enclave with {:?} {:?}", multiplicand, multiplier);
    sgx_status_t::SGX_SUCCESS
}
