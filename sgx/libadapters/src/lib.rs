#![feature(libc)]

extern crate base64;
extern crate errno;
extern crate libc;
extern crate sgx_types;
extern crate sgx_urts;
extern crate serde;
extern crate serde_derive;
#[macro_use]
extern crate serde_json;
extern crate utils;
#[macro_use]
extern crate lazy_static;
extern crate wasmi;

use sgx_types::*;
use sgx_urts::SgxEnclave;
use std::sync::{Arc, Mutex};

pub mod multiply;
pub mod wasm;
pub mod attestation;

static ENCLAVE_FILE: &str = "enclave.signed.so";

lazy_static! {
    static ref ENCLAVE: Arc<Mutex<SgxEnclave>> = {
        Arc::new(Mutex::new(enclave_init().unwrap_or_else(|err| {
            panic!("Failed to initialize the enclave: {}", err);
        })))
    };
}

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

// use_enclave takes a closure and passes in an enclave ID for making enclave calls in a thread
// safe way
pub fn use_enclave<F>(ctx: F)
    where F: Fn(u64) {
    let enclave = ENCLAVE.lock().unwrap();
    ctx(enclave.geteid())
}
