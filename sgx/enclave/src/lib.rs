#![crate_name = "enclave"]
#![crate_type = "staticlib"]

#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]
#![feature(alloc)]

extern crate base64;
extern crate num;
extern crate serde;
#[macro_use] extern crate serde_derive;
extern crate serde_json;
#[cfg(not(target_env = "sgx"))]
#[macro_use] extern crate sgx_tstd as std;
extern crate sgx_types;
#[macro_use] extern crate utils;
extern crate wasmi;

mod wasm;

use sgx_types::*;
use std::string::String;
use std::string::ToString;
use utils::{copy_string_to_cstr_ptr, string_from_cstr_with_len};

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
    OutputCStrError(utils::OutputCStrError),
    UnexpectedOutputError,
}

impl_from_error!(std::string::FromUtf8Error, WasmError::FromUtf8Error);
impl_from_error!(wasm::Error, WasmError::ExecError);
impl_from_error!(utils::OutputCStrError, WasmError::OutputCStrError);

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

#[no_mangle]
pub extern "C" fn sgx_multiply(
    adapter_str_ptr: *const u8,
    adapter_str_len: usize,
    input_str_ptr: *const u8,
    input_str_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> sgx_status_t {
    match multiply(
        adapter_str_ptr,
        adapter_str_len,
        input_str_ptr,
        input_str_len,
        result_ptr,
        result_capacity,
        result_len,
    ) {
        Ok(_) => sgx_status_t::SGX_SUCCESS,
        _ => sgx_status_t::SGX_ERROR_UNEXPECTED,
    }
}

fn multiply(
    adapter_str_ptr: *const u8,
    adapter_str_len: usize,
    input_str_ptr: *const u8,
    input_str_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> Result<(), sgx_status_t> {
    let adapter_str = string_from_cstr_with_len(adapter_str_ptr, adapter_str_len).unwrap();
    let adapter: serde_json::Value = match serde_json::from_str(&adapter_str) {
        Ok(result) => result,
        Err(_err) => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    };
    let input_str = string_from_cstr_with_len(input_str_ptr, input_str_len).unwrap();
    let mut input = parse_run_result_json(&input_str)?;
    let multiplier = get_json_string(&adapter, "times".to_string())?;
    let multiplicand = get_json_string(&input.data, "value".to_string())?;
    let result = parse_and_multiply(&multiplicand, &multiplier)?;

    input.status = Some("completed".to_string());
    input.add("value".to_string(), serde_json::Value::String(result));
    let rr_json = match serde_json::to_string(&input) {
        Ok(v) => v,
        _ => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    };

    match copy_string_to_cstr_ptr(&rr_json, result_ptr, result_capacity, result_len) {
        Err(_) => Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
        _ => Ok(())
    }
}

fn get_json_string(object: &serde_json::Value, key: String) -> Result<String, sgx_status_t> {
    match &object[key] {
        serde_json::Value::String(v) => Ok(v.clone()),
        serde_json::Value::Number(v) => Ok(format!("{}", v)),
        _ => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    }
}

fn parse_and_multiply(
    multiplicand_str: &str,
    multiplier_str: &str,
) -> Result<String, sgx_status_t> {
    let multiplicand = match i128::from_str_radix(multiplicand_str, 10) {
        Ok(result) => result,
        _ => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    };
    let multiplier = match i128::from_str_radix(multiplier_str, 10) {
        Ok(result) => result,
        _ => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    };

    Ok(format!("{:?}", multiplicand * multiplier))
}

fn parse_run_result_json(input: &str) -> Result<RunResult, sgx_status_t> {
    match serde_json::from_str(input) {
        Ok(rr) => Ok(rr),
        _ => return Err(sgx_status_t::SGX_ERROR_INVALID_PARAMETER),
    }
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

impl RunResult {
    fn add(&mut self, key: String, value: serde_json::Value) {
        self.data[key] = value;
    }
}
