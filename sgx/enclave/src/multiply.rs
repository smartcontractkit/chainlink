use result::{self, get_value, RunResult};
use std::string::ToString;
use bigdecimal::BigDecimal;
use std::str::FromStr;

#[derive(Debug)]
pub enum MultiplyError {
    ParseBigDecimalError(bigdecimal::ParseBigDecimalError),
    ResultError(result::Error),
}

impl_from_error!(bigdecimal::ParseBigDecimalError, MultiplyError::ParseBigDecimalError);
impl_from_error!(result::Error, MultiplyError::ResultError);

pub type MultiplyResult = Result<serde_json::Value, MultiplyError>;

pub fn perform(adapter: &serde_json::Value, input: &RunResult) -> MultiplyResult {
    let multiplier_str = get_value(&adapter, "times")?;
    let multiplicand_str = get_value(&input.data, "value")?;

    let multiplicand = BigDecimal::from_str(&multiplicand_str)?;
    let multiplier = BigDecimal::from_str(&multiplier_str)?;

    let result = multiplicand * multiplier;

    Ok(json!({"value": result.to_string()}))
}
