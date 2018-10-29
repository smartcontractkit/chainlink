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
        job_run_id: result.job_run_id.clone(),
        amount: result.amount,
        ..Default::default()
    }
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct RunResult {
    pub job_run_id: String,
    pub data: serde_json::Value,
    pub status: String,
    pub error: Option<String>,
    pub amount: Option<u64>,
}

impl RunResult {
    pub fn with_data(&self, data: &serde_json::Value) -> RunResult {
        RunResult {
            job_run_id: self.job_run_id.clone(),
            data: data.clone(),
            status: self.status.clone(),
            error: self.error.clone(),
            amount: self.amount.clone(),
        }
    }

    pub fn with_error(&self, error: &str) -> RunResult {
        RunResult {
            job_run_id: self.job_run_id.clone(),
            data: self.data.clone(),
            status: "errored".to_string(),
            error: Some(error.to_string()),
            amount: self.amount.clone(),
        }
    }

    pub fn with_status(&self, status: &str) -> RunResult {
        RunResult {
            job_run_id: self.job_run_id.clone(),
            data: self.data.clone(),
            status: status.to_string(),
            error: self.error.clone(),
            amount: self.amount.clone(),
        }
    }
}
