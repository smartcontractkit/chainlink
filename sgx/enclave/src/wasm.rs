use base64;
use num::Num;
use num::traits;
use std::num;
use wasmi::{self, ImportsBuilder, ModuleInstance, NopExternals};

#[derive(Debug)]
pub enum Error {
    Base64DecoderError(base64::DecodeError),
    EmptyResultError,
    ParseIntError(num::ParseIntError),
    ParseFloatError(traits::ParseFloatError),
    WasmError(wasmi::Error),
    WasmTrap(wasmi::Trap),
}

impl From<base64::DecodeError> for Error {
    fn from(e: base64::DecodeError) -> Self {
        Error::Base64DecoderError(e)
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

impl From<wasmi::Error> for Error {
    fn from(e: wasmi::Error) -> Self {
        Error::WasmError(e)
    }
}

impl From<wasmi::Trap> for Error {
    fn from(e: wasmi::Trap) -> Self {
        Error::WasmTrap(e)
    }
}

pub fn exec(encoded_program: &str, arguments: &str) -> Result<wasmi::RuntimeValue, Error> {
    let data = base64::decode(encoded_program)?;
    let module = wasmi::Module::from_buffer(data)?;
    let module_ref = ModuleInstance::new(&module, &ImportsBuilder::default())?;
    let instance = module_ref.run_start(&mut NopExternals)?;

    // FIXME: assumes a single value argument that is a float
    let value = <f64 as Num>::from_str_radix(arguments, 10)?;

    let arguments = &[wasmi::RuntimeValue::I64(value as i64)];
    match instance.invoke_export("perform", arguments, &mut NopExternals)? {
        Some(v) => Ok(v),
        _ => Err(Error::EmptyResultError),
    }
}
