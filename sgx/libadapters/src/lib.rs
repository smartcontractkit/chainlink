#![feature(libc)]

extern crate crossbeam;
extern crate errno;
extern crate libc;
extern crate sgx_types;
extern crate sgx_urts;
extern crate utils;

use crossbeam::crossbeam_channel::{unbounded, Sender, Receiver};
use errno::{set_errno, Errno};
use sgx_types::*;
use sgx_urts::SgxEnclave;
use std::cell::RefCell;
use std::ffi::CStr;
use std::sync::{Once, ONCE_INIT};
use std::{mem, thread};
use utils::{copy_string_to_cstr_ptr, cstr_len, string_from_cstr_with_len};

static ENCLAVE_FILE: &str = "enclave.signed.so";

extern "C" {
    fn enclave_adapter_perform(
        eid: sgx_enclave_id_t,
        retval: *mut sgx_status_t,
        adapter_type: *const u8,
        adapter_type_len: usize,
        adapter: *const u8,
        adapter_len: usize,
        input: *const u8,
        input_len: usize,
        result_ptr: *mut u8,
        result_capacity: usize,
        result_len: *mut usize) -> sgx_status_t;
}

#[derive(Clone, Debug)]
pub struct Workload {
    adapter_type: String,
    adapter: String,
    input: String,
}

#[derive(Clone, Debug)]
pub struct Response {
    result: String
}

#[derive(Clone)]
pub struct Context {
    sender: Sender<Workload>,
    receiver: Receiver<Response>,
}

impl Context {
    pub fn send(&self, workload: Workload) -> Response {
        self.sender.send(workload).expect("error sending to worker thread");
        self.receiver.recv().unwrap()
    }
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

pub fn get_context() -> Result<Context, ()> {
    static mut SINGLETON: *const Context = 0 as *const Context;
    static ONCE: Once = ONCE_INIT;

    unsafe {
        ONCE.call_once(|| {
            let (client_sender, worker_receiver): (Sender<Workload>, Receiver<Workload>) = unbounded();
            let (worker_sender, client_receiver): (Sender<Response>, Receiver<Response>) = unbounded();
            let singleton = Context {
                sender: client_sender,
                receiver: client_receiver,
            };

            thread::spawn(move || {
                for workload in worker_receiver.iter() {
                    println!("workload received {:?}", workload);

                    let enclave = enclave_init().unwrap();

                    let output: Vec<u8> = Vec::with_capacity(4096);
                    let mut output_len : usize = 0;
                    let mut retval = sgx_status_t::SGX_SUCCESS;
                    let result = unsafe {
                        enclave_adapter_perform(
                            enclave.geteid(),
                            &mut retval,
                            workload.adapter_type.as_ptr() as *const u8,
                            workload.adapter_type.len(),
                            workload.adapter.as_ptr() as *const u8,
                            workload.adapter.len(),
                            workload.input.as_ptr() as *const u8,
                            workload.input.len(),
                            output.as_ptr() as *mut u8,
                            output.capacity(),
                            &mut output_len as *mut usize)
                    };

                    let output_str = string_from_cstr_with_len(output.as_ptr(), output_len).unwrap();

                    println!("output_str {:?}", output_str);

                    worker_sender.send(Response{result: output_str}).expect("failed to send response");
                }
            });

            // Put the context in the heap so it can outlive this call
            SINGLETON = mem::transmute(Box::new(singleton));
        });

        Ok((*SINGLETON).clone())
    }
}

#[no_mangle]
pub extern "C" fn sgx_adapter_perform(
    adapter_type: *const libc::c_char,
    adapter: *const libc::c_char,
    input: *const libc::c_char,
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    let context = match get_context() {
        Ok(e) => e,
        Err(err) => {
            set_errno(Errno(1));
            return;
        }
    };

    let adapter_type = unsafe { CStr::from_ptr(adapter_type) }.to_str().unwrap().to_owned();
    let adapter = unsafe { CStr::from_ptr(adapter) }.to_str().unwrap().to_owned();
    let input = unsafe { CStr::from_ptr(input) }.to_str().unwrap().to_owned();

    println!("{:?} {{ adapter {:?} input {:?} }}", adapter_type, adapter, input);

    let response = context.send(Workload{
        adapter_type: adapter_type,
        adapter: adapter,
        input: input
    });

    copy_string_to_cstr_ptr(&response.result, result_ptr as *mut u8, result_capacity, result_len).unwrap();

    set_errno(Errno(0));
}
