#![crate_name = "enclave"]
#![crate_type = "staticlib"]

#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]
#![feature(alloc)]

extern crate base64;
extern crate num;
extern crate sgx_types;
#[cfg(not(target_env = "sgx"))]
#[macro_use] extern crate sgx_tstd as std;
extern crate wasmi;
extern crate serde;
#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate serde_json;

#[macro_use] mod util;
mod wasm;

use sgx_types::*;
use util::{copy_string_to_cstr_ptr, string_from_cstr_with_len};
use std::string::String;

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

pub extern "C" fn sgx_multiply(adapter_str_ptr: *const u8, adapter_str_len: usize, input_str_ptr: *const u8, input_str_len: usize) -> sgx_status_t {
    let adapter_str = string_from_cstr_with_len(adapter_str_ptr, adapter_str_len).unwrap();
    let adapter: serde_json::Value = match serde_json::from_str(&adapter_str) {
        Ok(result) => result,
        Err(_err) => return sgx_status_t::SGX_ERROR_INVALID_PARAMETER,
    };

    let input_str = string_from_cstr_with_len(input_str_ptr, input_str_len).unwrap();
    let input = match parse_run_result_json(&input_str) {
        Ok(result) => result,
        Err(_err) => return sgx_status_t::SGX_ERROR_INVALID_PARAMETER,
    };

    let multiplier = match &adapter["times"] {
        serde_json::Value::String(v) => v.clone(),
        serde_json::Value::Number(v) => format!("{}", v),
        _ => return sgx_status_t::SGX_ERROR_INVALID_PARAMETER,
    };


    let multiplicand = match &input.data["value"] {
        serde_json::Value::String(v) => v.clone(),
        serde_json::Value::Number(v) => format!("{}", v),
        _ => return sgx_status_t::SGX_ERROR_INVALID_PARAMETER,
    };

    println!("Performing MULTIPLY from within enclave with {:?} * {:?} = {:?}",
             multiplicand, multiplier, multiply(&multiplicand, &multiplier));

    sgx_status_t::SGX_SUCCESS
}

fn multiply(multiplicand_str: &str, multiplier_str: &str) -> Result<String, std::num::ParseIntError> {
    let multiplicand = match i128::from_str_radix(multiplicand_str, 10) {
        Ok(result) => result,
        Err(err) => return Err(err),
    };
    let multiplier = match i128::from_str_radix(multiplier_str, 10) {
        Ok(result) => result,
        Err(err) => return Err(err),
    };

    Ok(format!("{:?}", multiplicand * multiplier))
}

fn parse_run_result_json(input: &str) -> Result<RunResult, serde_json::Error> {
    let result: RunResult = serde_json::from_str(input)?;
    Ok(result)
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq)]
#[serde(rename_all = "camelCase")]
struct RunResult {
    job_run_id: String,
    data: serde_json::Value,
    status: Option<String>,
    error_message: Option<String>,
    amount: Option<u64>,
}
