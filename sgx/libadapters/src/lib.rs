#![feature(libc)]

extern crate errno;
extern crate libc;
extern crate sgx_types;
extern crate sgx_urts;
extern crate utils;

use sgx_types::*;
use sgx_urts::SgxEnclave;
use std::cell::RefCell;

pub mod multiply;
pub mod wasm;

static ENCLAVE_FILE: &str = "enclave.signed.so";
thread_local!(static ENCLAVE: RefCell<Option<SgxEnclave>> = RefCell::new(None));

fn enclave_init() -> SgxResult<SgxEnclave> {
    let mut launch_token: sgx_launch_token_t = [0; 1024];
    let mut launch_token_updated: i32 = 0;
    let debug = 1;
    let mut misc_attr = sgx_misc_attribute_t {
        secs_attr: sgx_attributes_t { flags: 0, xfrm: 0 },
        misc_select: 0,
    };
    SgxEnclave::create(
        ENCLAVE_FILE,
        debug,
        &mut launch_token,
        &mut launch_token_updated,
        &mut misc_attr,
    )
}

// get_enclave lazily initializes a thread local enclave context and returns its ID for use in
// ECALL/OCALLs
pub fn get_enclave() -> SgxResult<u64> {
    ENCLAVE.with(|e| {
        if let Some(ref e) = *e.borrow() {
            return Ok(e.geteid());
        }

        let enclave = enclave_init()?;
        let enclave_id = enclave.geteid();
        *e.borrow_mut() = Some(enclave);
        Ok(enclave_id)
    })
}
