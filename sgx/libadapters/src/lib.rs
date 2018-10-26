#![feature(libc)]
extern crate libc;

extern crate errno;
extern crate sgx_types;
extern crate sgx_urts;
extern crate utils;

#[macro_use]
extern crate lazy_static;

use std::panic;

use errno::{set_errno, Errno};
use sgx_types::*;
use sgx_urts::SgxEnclave;

pub mod multiply;
pub mod wasm;

static ENCLAVE_FILE: &'static str = "enclave.signed.so";

// lazy_static allows us to setup a global value to persist the enclave in, so that we don't have
// to pass the SgxEnclave into go territory and then back again.
lazy_static! {
    static ref ENCLAVE: SgxEnclave = {
        perform_enclave_init().unwrap_or_else(|err| {
            panic!("Failed to initialize the enclave: {}", err);
        })
    };
}

#[no_mangle]
pub extern "C" fn init_enclave() {
    // lazy_statics don't have a way to return an error, so setup a panic handler.
    let result = panic::catch_unwind(|| {
        lazy_static::initialize(&ENCLAVE);
    });
    set_errno(Errno(0));
    if result.is_err() {
        // Go uses the C _errno variable to get errors from C
        set_errno(Errno(1));
    }
}

fn perform_enclave_init() -> SgxResult<SgxEnclave> {
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
