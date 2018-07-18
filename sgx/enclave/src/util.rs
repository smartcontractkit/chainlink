// Util functions, primarily for passing data between Go / Rust / SGX enclaves

use std::ffi::{CString, NulError};
use std::ptr;
use std::slice;
use std::slice::from_raw_parts_mut;
use std::string::{FromUtf8Error, String};

// string_from_cstr_with_len creates a rust String from a string pointer with a specific length.
pub fn string_from_cstr_with_len(ptr: *const u8, len: usize) -> Result<String, FromUtf8Error> {
    let slice = unsafe { slice::from_raw_parts(ptr, len) };
    String::from_utf8(slice.to_vec())
}

#[derive(Debug)]
pub enum OutputCStrError {
    NulError,
    CapacityError,
}

impl From<NulError> for OutputCStrError {
    fn from(_: NulError) -> Self {
        OutputCStrError::NulError
    }
}

// copy_string_to_cstr_ptr takes in input rust String, and copies it to an output pointer
// with a specific capacity, storing the input string length in output_len.
pub fn copy_string_to_cstr_ptr(
    input: &str,
    output_ptr: *mut u8,
    output_capacity: usize,
    output_len: *mut usize,
) -> Result<(), OutputCStrError> {

    let input_cstring = CString::new(input)?;
    let input_slice = input_cstring.to_bytes();
    let input_size = input_slice.len();

    if output_capacity < input_size {
        return Err(OutputCStrError::CapacityError);
    }

    unsafe {
        from_raw_parts_mut(output_ptr, input_size).copy_from_slice(input_slice);
        ptr::copy(&input_size, output_len, 1);
    }

    Ok(())
}
