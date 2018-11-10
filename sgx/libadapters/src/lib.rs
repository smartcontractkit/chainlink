#![feature(libc)]

extern crate errno;
#[macro_use]
extern crate lazy_static;
extern crate libc;
extern crate pool;
extern crate sgx_types;
extern crate sgx_urts;
extern crate utils;

use pool::Pool;
use sgx_types::*;
use sgx_urts::SgxEnclave;
use std::cell::RefCell;
use std::sync::{Arc, Mutex};

pub mod multiply;
pub mod wasm;

// ENCLAVE_FILE is the path to the signed enclave
static ENCLAVE_FILE: &str = "enclave.signed.so";

// MAX_ECALLS controls the maximum number of concurrent enclave calls can be made
static MAX_ECALLS: usize = 2;

type LazyEnclave = Arc<Mutex<RefCell<Option<SgxEnclave>>>>;
type LazyEnclavePool = Arc<Mutex<Pool<LazyEnclave>>>;

lazy_static! {
    static ref ENCLAVE_POOL: LazyEnclavePool =
        Arc::new(Mutex::new(Pool::with_capacity(MAX_ECALLS, 0, || Arc::new(Mutex::new(RefCell::new(None))))));
}

fn new_enclave() -> SgxResult<SgxEnclave> {
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

    let lazy_enclave: LazyEnclave = {
        Arc::clone(&ENCLAVE_POOL.lock().unwrap().checkout().unwrap())
    };

    let cell: &RefCell<Option<SgxEnclave>> = &*lazy_enclave.lock().unwrap();
    let enclave_id = || -> u64 {
        if let Some(ref e) = *cell.borrow() {
            return e.geteid();
        }

        let e = new_enclave().unwrap();
        println!("new enclave {:?}", e);
        let eid = e.geteid();
        *cell.borrow_mut() = Some(e);
        eid
    }();

    ctx(enclave_id)
}
