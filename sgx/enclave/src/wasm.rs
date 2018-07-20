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

impl_from_error!(base64::DecodeError, Error::Base64DecoderError);
impl_from_error!(num::ParseIntError, Error::ParseIntError);
impl_from_error!(traits::ParseFloatError, Error::ParseFloatError);
impl_from_error!(wasmi::Error, Error::WasmError);
impl_from_error!(wasmi::Trap, Error::WasmTrap);

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
