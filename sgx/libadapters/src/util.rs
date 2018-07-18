use libc;
use std::ffi::CStr;

pub fn cstr_len(string: *const libc::c_char) -> usize {
    let buffer = unsafe { CStr::from_ptr(string).to_bytes() };
    buffer.len()
}
