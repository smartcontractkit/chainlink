use result::{self, get_value, RunResult};
use std::string::ToString;

#[derive(Debug)]
pub enum MultiplyError {
    ParseFloatError(core::num::ParseFloatError),
    ResultError(result::Error),
}

impl_from_error!(core::num::ParseFloatError, MultiplyError::ParseFloatError);
impl_from_error!(result::Error, MultiplyError::ResultError);

pub type MultiplyResult = Result<serde_json::Value, MultiplyError>;

pub fn perform(adapter: &serde_json::Value, input: &RunResult) -> MultiplyResult {
    let multiplier_str = get_value(&adapter, "times")?;
    let multiplicand_str = get_value(&input.data, "value")?;

    let multiplicand = multiplicand_str.parse::<f64>()?;
    let multiplier = multiplier_str.parse::<f64>()?;

    let result = multiplicand * multiplier;

    Ok(json!({"value": result.to_string()}))
}
