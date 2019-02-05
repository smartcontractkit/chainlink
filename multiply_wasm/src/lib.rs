#![feature(link_args)]
#![allow(unused_attributes)] // link_args actually is used
#![link_args = "--import-memory"]

#[cfg(no_std)]

extern crate serde;
extern crate serde_derive;
extern crate serde_json;

use std::str::FromStr;
use std::ffi::CStr;

#[derive(Debug)]
pub enum Error {
    InvalidEncoding,
}

#[no_mangle]
pub extern "C" fn perform(input_ptr: *const i8) {
    let input_str = unsafe { CStr::from_ptr(input_ptr) }.to_str()
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
