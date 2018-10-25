use base64;
use std::{num, vec::Vec};
use wasmi::{self, ImportsBuilder, ModuleInstance, NopExternals};

use result::{self, RunResult, get_value};

#[derive(Debug)]
pub enum WasmError {
    ArgumentTypeNotImplementedYet,
    Base64DecoderError(base64::DecodeError),
    EmptyResultError,
    ParseIntError(num::ParseIntError),
    ParseFloatError(num::ParseFloatError),
    ResultError(result::Error),
    WasmError(wasmi::Error),
    WasmTrap(wasmi::Trap),
}

impl_from_error!(base64::DecodeError, WasmError::Base64DecoderError);
impl_from_error!(num::ParseIntError, WasmError::ParseIntError);
impl_from_error!(num::ParseFloatError, WasmError::ParseFloatError);
impl_from_error!(result::Error, WasmError::ResultError);
impl_from_error!(wasmi::Error, WasmError::WasmError);
impl_from_error!(wasmi::Trap, WasmError::WasmTrap);

pub type WasmResult = Result<serde_json::Value, WasmError>;

pub fn perform(adapter: &serde_json::Value, input: &RunResult) -> WasmResult {
    let encoded_program = get_value(&adapter, "wasm")?;
    let data = base64::decode(&encoded_program)?;
    let module = wasmi::Module::from_buffer(data)?;
    let module_ref = ModuleInstance::new(&module, &ImportsBuilder::default())?;
    let instance = module_ref.run_start(&mut NopExternals)?;

    let arguments = json_as_wasm_arguments(&input.data["value"])?;
    match instance.invoke_export("perform", &arguments.as_slice(), &mut NopExternals)? {
        Some(v) => Ok(json!({"value": wasm_as_json(&v)?})),
        _ => Err(WasmError::EmptyResultError),
    }
}

fn wasm_as_json(input: &wasmi::RuntimeValue) -> Result<serde_json::Value, WasmError> {
    match input {
        // RunResult in chainlink only supports string values
        wasmi::RuntimeValue::I32(v) => Ok(serde_json::Value::String(format!("{}", v))),
        _ => Err(WasmError::ArgumentTypeNotImplementedYet),
    }
}

fn json_as_wasm_arguments(input: &serde_json::Value) -> Result<Vec<wasmi::RuntimeValue>, WasmError> {
    match input {
        serde_json::Value::Array(vec) => vec.into_iter().map(json_to_wasm).collect(),

        // No value was specified in the input RunResult
        serde_json::Value::Null => Ok(vec![]),

        _ => Ok(vec![json_to_wasm(input)?]),
    }
}

fn json_to_wasm(input: &serde_json::Value) -> Result<wasmi::RuntimeValue, WasmError> {
    match input {
        serde_json::Value::Number(num) => {
            if input.is_f64()  {
                Ok(wasmi::RuntimeValue::F64(num.as_f64().unwrap().into()))
            } else if input.is_i64() {
                Ok(wasmi::RuntimeValue::F64((num.as_i64().unwrap() as f64).into()))
            } else {
                Err(WasmError::ArgumentTypeNotImplementedYet)
            }
        },
        _ => Err(WasmError::ArgumentTypeNotImplementedYet)
    }
}
