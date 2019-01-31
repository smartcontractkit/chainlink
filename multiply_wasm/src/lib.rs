#[cfg(no_std)]

extern crate serde;
extern crate serde_derive;
extern crate serde_json;

use std::slice;
use std::str::FromStr;

#[derive(Debug)]
pub enum Error {
    InvalidEncoding,
}

// string_from_cstr_with_len creates a rust String from a string pointer with a specific length.
fn string_from_cstr_with_len(ptr: *const u8, len: usize) -> Result<String, std::string::FromUtf8Error> {
    let slice = unsafe { slice::from_raw_parts(ptr, len) };
    String::from_utf8(slice.to_vec())
}

#[no_mangle]
pub extern "C" fn perform(input_ptr: *const u8, input_str_len: usize) {
    let input_str = string_from_cstr_with_len(input_ptr, input_str_len)
        .expect("error converting input string");
    let input: serde_json::Value = serde_json::from_str(&input_str)
        .expect("failed to parse input");

    let multiplier_str = match &input.pointer("/adapter/times") {
        Some(serde_json::Value::String(v)) => v,
        _ => panic!("no times value in adapter"),
    };
    let multiplicand_str = match &input.pointer("/input/value/") {
        Some(serde_json::Value::String(v)) => v,
        _ => panic!("no value param in input"),
    };

    let multiplicand = f64::from_str(&multiplicand_str)
        .expect("invalid multiplicand");
    let multiplier = f64::from_str(&multiplier_str)
        .expect("invalid multiplier");

    let _result = multiplicand * multiplier;
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
