use result::{self, get_value, RunResult};
use bigdecimal::BigDecimal;
use std::str::FromStr;
use std::string::String;

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
    let multiplicand_str = get_value(&input.data, "result")?;

    let multiplicand = BigDecimal::from_str(&multiplicand_str)?;
    let multiplier = BigDecimal::from_str(&multiplier_str)?;

    let result = multiplicand * multiplier;

    Ok(json!({"result": format_decimal(&result)}))
}

// format_decimal returns the result without any trailing 0s
fn format_decimal(value: &BigDecimal) -> String {
    let output = format!("{}", value);
    if output.contains(".") {
        return output.trim_end_matches('0').trim_end_matches('.').into()
    }
    output
}
