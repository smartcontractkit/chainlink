#![crate_name = "enclave"]
#![crate_type = "staticlib"]
#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

extern crate base64;
extern crate bigdecimal;
extern crate num;
extern crate serde;
#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate serde_json;
extern crate sgx_rand;
#[cfg(not(target_env = "sgx"))]
#[macro_use]
extern crate sgx_tstd as std;
extern crate sgx_tse as tse;
extern crate sgx_types;
extern crate sgx_tcrypto;
#[macro_use]
extern crate utils;
extern crate wasmi;

mod attestation;
mod multiply;
mod result;
mod wasm;

use result::RunResult;
use sgx_types::*;
use utils::{copy_string_to_cstr_ptr, string_from_cstr_with_len};
use std::string::ToString;

#[derive(Debug)]
enum ShimError {
    OutputCStrError(utils::OutputCStrError),
    FromUtf8Error(std::string::FromUtf8Error),
    JsonError(serde_json::Error),
}

impl_from_error!(utils::OutputCStrError, ShimError::OutputCStrError);
impl_from_error!(std::string::FromUtf8Error, ShimError::FromUtf8Error);
impl_from_error!(serde_json::Error, ShimError::JsonError);

#[no_mangle]
pub extern "C" fn sgx_wasm(
    adapter_str_ptr: *const u8,
    adapter_str_len: usize,
    input_str_ptr: *const u8,
    input_str_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> sgx_status_t {
    match wasm_shim(
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

fn wasm_shim(
    adapter_str_ptr: *const u8,
    adapter_str_len: usize,
    input_str_ptr: *const u8,
    input_str_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> Result<(), ShimError> {
    let adapter_str = string_from_cstr_with_len(adapter_str_ptr, adapter_str_len)?;
    let adapter = serde_json::from_str(&adapter_str)?;
    let input_str = string_from_cstr_with_len(input_str_ptr, input_str_len)?;
    let input: RunResult = serde_json::from_str(&input_str)?;

    let result = match wasm::perform(&adapter, &input) {
        Ok(value) => result::new(&input)
            .with_data(&value)
            .with_status("completed"),
        Err(err) => result::new(&input).with_error(&format!("{:?}", err)),
    };

    let rr_json = serde_json::to_string(&result)?;
    copy_string_to_cstr_ptr(&rr_json, result_ptr, result_capacity, result_len)?;
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
    match multiply_shim(
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

fn multiply_shim(
    adapter_str_ptr: *const u8,
    adapter_str_len: usize,
    input_str_ptr: *const u8,
    input_str_len: usize,
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> Result<(), ShimError> {
    let adapter_str = string_from_cstr_with_len(adapter_str_ptr, adapter_str_len)?;
    let adapter = serde_json::from_str(&adapter_str)?;
    let input_str = string_from_cstr_with_len(input_str_ptr, input_str_len)?;
    let input: RunResult = serde_json::from_str(&input_str)?;

    let result = match multiply::perform(&adapter, &input) {
        Ok(value) => result::new(&input)
            .with_data(&value)
            .with_status("completed"),
        Err(err) => result::new(&input).with_error(&format!("{:?}", err)),
    };

    let rr_json = serde_json::to_string(&result)?;
    copy_string_to_cstr_ptr(&rr_json, result_ptr, result_capacity, result_len)?;
    Ok(())
}

#[no_mangle]
pub extern "C" fn sgx_report(
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> sgx_status_t {
    match report_shim(
        result_ptr,
        result_capacity,
        result_len,
    ) {
        Ok(_) => sgx_status_t::SGX_SUCCESS,
        _ => sgx_status_t::SGX_ERROR_UNEXPECTED,
    }
}

fn report_shim(
    result_ptr: *mut u8,
    result_capacity: usize,
    result_len: *mut usize,
) -> Result<(), ShimError> {
    let output = match attestation::report() {
        Ok(report) => json!({
            "report": {
                "body": {
                    "report_data": report.body.report_data.d.to_vec(),
                    "mr_enclave": report.body.mr_enclave.m.to_vec(),
                },
                "key_id": report.key_id.id,
                "mac": report.mac,
            }
        }).to_string(),
        Err(err) => format!("error: {:?}", err),
    };

    copy_string_to_cstr_ptr(&output, result_ptr, result_capacity, result_len)?;
    Ok(())
}
