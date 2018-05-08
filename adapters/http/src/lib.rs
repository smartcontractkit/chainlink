extern crate reqwest;
extern crate libc;
use std::ffi::{CStr, CString};
use std::ptr;
use reqwest::header::ContentType;

#[no_mangle]
pub extern "C" fn perform_http_get(url: *const libc::c_char) -> *const libc::c_char {
    let buf_url = unsafe { CStr::from_ptr(url).to_bytes() };
    let str_url = String::from_utf8(buf_url.to_vec()).unwrap();

    let result = http_get(&str_url);
    match result {
        Err(_) => ptr::null(),
        Ok(body) => CString::new(body).unwrap().into_raw(),
    }
}

fn http_get(url: &String) -> Result<String, reqwest::Error> {
    let client = reqwest::Client::new();
    let body = client.get(url)
        .header(ContentType::json())
        .send()?
        .text()?;

    Ok(body.to_string())
}

#[no_mangle]
pub extern "C" fn perform_http_post(url: *const libc::c_char, body: *const libc::c_char) -> *const libc::c_char {
    let buf_url = unsafe { CStr::from_ptr(url).to_bytes() };
    let str_url = String::from_utf8(buf_url.to_vec()).unwrap();

    let buf_body = unsafe { CStr::from_ptr(body).to_bytes() };
    let str_body = String::from_utf8(buf_body.to_vec()).unwrap();

    let result = http_post(&str_url, str_body);
    match result {
        Err(_) => ptr::null(),
        Ok(body) => CString::new(body).unwrap().into_raw(),
    }
}

fn http_post(url: &String, body: String) -> Result<String, reqwest::Error> {
    let client = reqwest::Client::new();
    let body = client.post(url)
        .header(ContentType::json())
        .body(body)
        .send()?
        .text()?;
    
    Ok(body.to_string())
}
