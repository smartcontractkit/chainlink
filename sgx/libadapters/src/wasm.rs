use libc;
use sgx_types::*;
use base64;
use std::{num, vec::Vec};
use wasmi::{self, ImportsBuilder, ModuleInstance, NopExternals, ExternVal};
use std::ffi::CStr;

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
    adapter_ptr: *const libc::c_char,
    input_ptr: *const libc::c_char,
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    let adapter_str = unsafe { CStr::from_ptr(adapter_ptr) }.to_str().expect("from_ptr failed on adapter_ptr");
    let adapter : serde_json::Value = serde_json::from_str(&adapter_str)
        .expect("serde_json::from_str failed on adapter_str");

    let encoded_program = &adapter.pointer("/wasm")
        .expect("no wasm in data")
        .as_str().expect("input not string");
    let data = base64::decode(&encoded_program).expect("base64::decode failed");
    let module = wasmi::Module::from_buffer(data).expect("wasmi::from_buffer failed");
    let module_ref = ModuleInstance::new(&module, &ImportsBuilder::default()).expect("ModuleInstance::new failed");


    let instance = module_ref.run_start(&mut NopExternals).expect("module_ref.run_start failed");

    let arguments = encode_json_as_wasm(input_ptr);
    let output = match instance.invoke_export("perform", &arguments.as_slice(), &mut NopExternals)
        .expect("instance.invoke_export failed") {
        Some(v) => json!({"result": wasm_as_json(&v)}),
        _ => panic!("empty result error"),
    };

    println!("output: {:?}", output);
}

fn encode_json_as_wasm(input: *const libc::c_char) -> Vec<wasmi::RuntimeValue> {
    vec![wasmi::RuntimeValue::I32(input as i32)]
}

fn wasm_as_json(input: &wasmi::RuntimeValue) -> serde_json::Value {
    match input {
        // RunResult in chainlink only supports string values
        wasmi::RuntimeValue::I32(v) => serde_json::Value::String(format!("{}", v)),
        _ => panic!("argument type not implemented yet"),
    }
}
