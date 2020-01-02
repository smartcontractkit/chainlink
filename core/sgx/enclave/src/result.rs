use std::string::{String, ToString};

#[derive(Debug)]
pub enum Error {
    InvalidEncoding,
}

pub fn get_value(object: &serde_json::Value, key: &str) -> Result<String, Error> {
    match &object[key] {
        serde_json::Value::String(v) => Ok(v.clone()),
        serde_json::Value::Number(v) => Ok(format!("{}", v)),
        _ => return Err(Error::InvalidEncoding),
    }
}

pub fn new(result: &RunResult) -> RunResult {
    RunResult {
        ..Default::default()
    }
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct RunResult {
    pub data: serde_json::Value,
    pub status: String,
    pub error: Option<String>,
}

impl RunResult {
    pub fn with_data(&self, data: &serde_json::Value) -> RunResult {
        RunResult {
            data: data.clone(),
            status: self.status.clone(),
            error: self.error.clone(),
        }
    }

    pub fn with_error(&self, error: &str) -> RunResult {
        RunResult {
            data: self.data.clone(),
            status: "errored".to_string(),
            error: Some(error.to_string()),
        }
    }

    pub fn with_status(&self, status: &str) -> RunResult {
        RunResult {
            data: self.data.clone(),
            status: status.to_string(),
            error: self.error.clone(),
        }
    }
}
