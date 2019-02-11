use libc;
use base64;

use std::ffi::CStr;
use std::cell::RefCell;

use wasmi::{self, ModuleInstance, MemoryInstance};
use wasmi::{memory_units, RuntimeValue, Externals, Error as WasmError, ValueType, MemoryRef};
use wasmi::{ModuleImportResolver, MemoryDescriptor, Trap};
use wasmi::memory_units::{Bytes, Pages, RoundUpTo};
use wasmi::LINEAR_MEMORY_PAGE_SIZE;

mod ids {
    pub const FUNC_DEBUG: usize = 1;
}

const MAX_MEM: u32 = 1024 * 1024 * 1024; // 1 GiB

struct Resolver {
    max_memory: u32, // in pages.
    memory: RefCell<Option<MemoryRef>>,
}

impl ModuleImportResolver for Resolver {
    fn resolve_func(
        &self,
        field_name: &str,
        signature: &wasmi::Signature
    ) -> Result<wasmi::FuncRef, WasmError> {
        match field_name {
            "debug" => {
                let index = ids::FUNC_DEBUG;
                let (params, ret_ty): (&[ValueType], Option<ValueType>) = (&[ValueType::I32], None);

                if signature.params() != params && signature.return_type() != ret_ty {
                    Err(WasmError::Instantiation(
                        format!("Export {} has a bad signature", field_name)
                    ))
                } else {
                    Ok(wasmi::FuncInstance::alloc_host(
                        wasmi::Signature::new(&params[..], ret_ty),
                        index,
                    ))
                }
            }
            _ => {
                Err(WasmError::Instantiation(
                    format!("Export {} not found", field_name),
                ))
            }
        }

    }

    fn resolve_memory(
        &self,
        field_name: &str,
        descriptor: &MemoryDescriptor,
        ) -> Result<MemoryRef, WasmError> {
        if field_name == "memory" {
            let effective_max = descriptor.maximum().unwrap_or(self.max_memory);
            if descriptor.initial() > self.max_memory || effective_max > self.max_memory {
                Err(WasmError::Instantiation("Module requested too much memory".to_owned()))
            } else {
                let mem = MemoryInstance::alloc(
                    memory_units::Pages(descriptor.initial() as usize),
                    descriptor.maximum().map(|x| memory_units::Pages(x as usize)),
                )?;
                *self.memory.borrow_mut() = Some(mem.clone());
                Ok(mem)
            }
        } else {
            Err(WasmError::Instantiation("Memory imported under unknown name".to_owned()))
        }
    }
}

struct ValidationExternals<'a> {
    memory: &'a MemoryRef,
}

impl<'a> ValidationExternals<'a> {
    fn debug(&mut self, args: ::wasmi::RuntimeArgs) -> Result<(), Trap> {
        let target: u32 = args.nth_checked(0)?;
        
        println!("debug({:?})", target);
        Ok(())
        //let data_ptr: u32 = args.nth_checked(1)?;
        //let data_len: u32 = args.nth_checked(2)?;

        //let (data_ptr, data_len) = (data_ptr as usize, data_len as usize);

        //self.memory.with_direct_access(|mem| {
        //if mem.len() < (data_ptr + data_len) {
        //Err(Trap::new(wasmi::TrapKind::MemoryAccessOutOfBounds))
        //} else {
        //let res = self.externalities.post_message(MessageRef {
        //target,
        //data: &mem[data_ptr..][..data_len],
        //});

        //res.map_err(|e| Trap::new(wasmi::TrapKind::Host(
        //Box::new(e) as Box<_>
        //)))
        //}
        //})
    }
}

impl<'a> Externals for ValidationExternals<'a> {
    fn invoke_index(
        &mut self,
        index: usize,
        args: ::wasmi::RuntimeArgs,
        ) -> Result<Option<RuntimeValue>, Trap> {
        match index {
            ids::FUNC_DEBUG => self.debug(args).map(|_| None),
            _ => panic!("no externality at given index"),
        }
    }
}

#[no_mangle]
pub extern "C" fn wasm(
    adapter_ptr: *const libc::c_char,
    input_ptr: *const libc::c_char,
    result_ptr: *mut libc::c_char,
    result_capacity: usize,
    result_len: *mut usize,
) {
    let adapter_str = unsafe { CStr::from_ptr(adapter_ptr) }.to_str()
        .expect("from_ptr failed on adapter_ptr");
    let adapter : serde_json::Value = serde_json::from_str(&adapter_str)
        .expect("serde_json::from_str failed on adapter_str");
    let input_str = unsafe { CStr::from_ptr(input_ptr) }.to_str()
        .expect("from_ptr failed on input_ptr");
    let input : serde_json::Value = serde_json::from_str(&input_str)
        .expect("serde_json::from_str failed on adapter_str");

    let encoded_program = &adapter.pointer("/wasm")
        .expect("no wasm in data")
        .as_str().expect("input not string");
    let data = base64::decode(&encoded_program).expect("base64::decode failed");

    let module_resolver = Resolver {
        max_memory: MAX_MEM / LINEAR_MEMORY_PAGE_SIZE.0 as u32,
        memory: RefCell::new(None),
    };

    let module = wasmi::Module::from_buffer(data)
        .expect("from_buffer failed");
    
    let module_ref = ModuleInstance::new(
        &module,
        &wasmi::ImportsBuilder::new().with_resolver("env", &module_resolver),
    ).expect("ModuleInstance::new failed");

    let memory = module_resolver.memory.borrow()
        .as_ref()
        .expect("no imported memory instance")
        .clone();

    let mut externals = ValidationExternals {memory: &memory};

    let instance = module_ref.run_start(&mut externals)
        .expect("module_ref.run_start failed");

    let input = json!({"input": input, "adapter": {
        "times": "2",
    }}).to_string();

    println!("input {:?}", input);

    let call_data_pages: Pages = Bytes(input.len()).round_up_to();
    let allocated_mem_start: Bytes = memory.grow(call_data_pages).unwrap().into();
    memory.set(allocated_mem_start.0 as u32, input.as_bytes())
        .expect("not enough memory allocated for input");

    let arguments = vec![RuntimeValue::I32(allocated_mem_start.0 as i32)];
    let output = match instance.invoke_export("perform", &arguments.as_slice(), &mut externals).expect("instance.invoke_export failed") {
        Some(v) => json!({"result": wasm_as_json(&v)}),
        _ => panic!("empty result error"),
    };

    println!("output: {:?}", output);
}

fn wasm_as_json(input: &wasmi::RuntimeValue) -> serde_json::Value {
    match input {
        // RunResult in chainlink only supports string values
        wasmi::RuntimeValue::I32(v) => serde_json::Value::String(format!("{}", v)),
        _ => panic!("argument type not implemented yet"),
    }
}
