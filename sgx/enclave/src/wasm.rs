use base64;
use num::traits;
use num::Num;
use sgxwasm;
use std::num;
use wasmi::{self, ImportsBuilder, ModuleInstance, NopExternals};

#[derive(Debug)]
pub enum Error {
    Base64DecoderError(base64::DecodeError),
    EmptyResultError,
    ParseIntError(num::ParseIntError),
    ParseFloatError(traits::ParseFloatError),
    WasmError(sgxwasm::Error),
    WasmTrap(wasmi::Trap),
}

impl From<base64::DecodeError> for Error {
    fn from(e: base64::DecodeError) -> Self {
        Error::Base64DecoderError(e)
    }
}

impl From<sgxwasm::Error> for Error {
    fn from(e: sgxwasm::Error) -> Self {
        Error::WasmError(e)
    }
}

impl From<wasmi::Trap> for Error {
    fn from(e: wasmi::Trap) -> Self {
        Error::WasmTrap(e)
    }
}

impl From<num::ParseIntError> for Error {
    fn from(e: num::ParseIntError) -> Self {
        Error::ParseIntError(e)
    }
}

impl From<traits::ParseFloatError> for Error {
    fn from(e: traits::ParseFloatError) -> Self {
        Error::ParseFloatError(e)
    }
}

pub fn exec(encoded_program: &str, arguments: &str) -> Result<wasmi::RuntimeValue, Error> {
    println!("exec: {:?}", encoded_program);
    let data = base64::decode(encoded_program)?;
    println!("data: {:?}", data);
    // FIXME: Compiler can't find the From trait for the below?
    let module = wasmi::Module::from_buffer(data) //?;
        .expect("module error");
    println!("module loaded");
    let module_ref = ModuleInstance::new(&module, &ImportsBuilder::default()) //?;
        .expect("module_ref error");
    println!("module_ref loaded");
    let instance = module_ref.run_start(&mut NopExternals) //?;
        .expect("run_start error");

    println!("instance loaded");
    // FIXME: assumes a single value argument that is a float
    let value = <f64 as Num>::from_str_radix(arguments, 10)?;

    println!("value loaded {:?}", value);
    let arguments = [wasmi::RuntimeValue::I64(value as i64)];

    println!("loading permission: {:?}", instance);
    match instance
        .invoke_export("perform", &arguments, &mut NopExternals)
        .expect("invoke_export")
    {
        Some(v) => Ok(v),
        _ => Err(Error::EmptyResultError),
    }
}
