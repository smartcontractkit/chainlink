// Util functions, primarily for passing data between Go / Rust / SGX enclaves
//
// NOTE: This library is intended to be unit tested, in test mode it uses the std rust library as
// this is needed for Rust's testing framework. Take care not to introduce a dependency on anything
// in std which is not also in sgx_tstd.

#![cfg_attr(all(not(test), not(feature = "std")), no_std)]

#[cfg(all(not(test), not(feature = "std")))]
extern crate sgx_tstd as std;

use std::ffi::{CStr, CString, NulError};
use std::ptr;
use std::slice;
use std::string::{FromUtf8Error, String};

// cstr_len gets the length of a cstring pointer, not including the null terminator
pub fn cstr_len(string: *const i8) -> usize {
    let buffer = unsafe { CStr::from_ptr(string).to_bytes() };
    buffer.len()
}

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
        ptr::copy_nonoverlapping(input_slice.as_ptr(), output_ptr, input_size);
        ptr::copy(&input_size, output_len, 1);
    }

    Ok(())
}

#[macro_export]
macro_rules! impl_from_error {
    ($from:path, $to:tt::$ctor:tt) => {
        impl From<$from> for $to {
            fn from(e: $from) -> Self {
                $to::$ctor(e)
            }
        }
    };
}

#[cfg(test)]
mod tests {
    use super::{copy_string_to_cstr_ptr, cstr_len, string_from_cstr_with_len};
    use std::ffi::CString;

    #[test]
    fn test_string_from_cstr_with_len() {
        let cstr = b"hello world!";
        let result = string_from_cstr_with_len(cstr as *const u8, 12);

        assert!(result.is_ok());
        assert_eq!(result.unwrap(), "hello world!");
    }

    #[test]
    fn test_copy_string_to_cstr_ptr() {
        let mut buffer: [u8; 64] = [0; 64];
        let mut size: usize = 0;

        let result = copy_string_to_cstr_ptr(
            "hello world!".into(),
            &mut buffer[0],
            buffer.len(),
            &mut size,
        );

        assert!(result.is_ok());
        assert_eq!(size, 12);
        assert_eq!(
            String::from_utf8_lossy(&buffer[..size]),
            "hello world!".to_string()
        );
    }

    #[test]
    fn test_copy_string_to_cstr_ptr_insufficient_capacity() {
        let mut buffer: [u8; 10] = [0; 10];
        let mut size: usize = 0;

        let result = copy_string_to_cstr_ptr(
            "hello world!".into(),
            &mut buffer[0],
            buffer.len(),
            &mut size,
        );
        assert!(result.is_err());
    }

    #[test]
    fn test_cstr_len() {
        let cstr = CString::new("hello world!").unwrap();
        assert_eq!(cstr_len(cstr.as_ptr()), 12);
    }
}
