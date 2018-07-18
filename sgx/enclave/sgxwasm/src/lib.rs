// Copyright (C) 2017-2018 Baidu, Inc. All Rights Reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  * Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in
//    the documentation and/or other materials provided with the
//    distribution.
//  * Neither the name of Baidu, Inc., nor the names of its
//    contributors may be used to endorse or promote products derived
//    from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

#![crate_name = "sgxwasm"]
#![crate_type = "staticlib"]

#![cfg_attr(not(target_env = "sgx"), no_std)]
#![cfg_attr(target_env = "sgx", feature(rustc_private))]

extern crate sgx_types;
#[cfg(not(target_env = "sgx"))]
#[macro_use]
extern crate sgx_tstd as std;
extern crate wasmi;
extern crate wabt_core;

use std::{i32, i64, u32, u64, f32};
use std::prelude::v1::*;
use std::collections::HashMap;
use wasmi::memory_units::Pages;

pub use wasmi::Error as InterpreterError;
use wasmi::{ModuleInstance,
            ImportsBuilder,
            RuntimeValue,
//          NopExternals,
            MemoryInstance,
            GlobalInstance,
            GlobalRef,
            TableRef,
            MemoryRef,
            TableInstance,
            Trap,
            Externals,
            RuntimeArgs,
            FuncRef,
            Signature,
            FuncInstance,
            ModuleImportResolver,
            TableDescriptor,
            MemoryDescriptor,
            GlobalDescriptor,
            ModuleRef,
            ImportResolver,
            Module,
};

use wabt_core::script;
use wabt_core::script::{Value};

extern crate serde;
#[macro_use]
extern crate serde_derive;
//use serde::{Serialize, Serializer, Deserialize, Deserializer};
#[derive(Debug, Serialize, Deserialize)]
pub enum SgxWasmAction {
    Invoke {
        module: Option<String>,
        field: String,
        args: Vec<BoundaryValue>
    },
    Get {
        module: Option<String>,
        field: String,
    },
    LoadModule {
        name: Option<String>,
        module: Vec<u8>,
    },
    TryLoad {
        module: Vec<u8>,
    },
    Register {
        name: Option<String>,
        as_name: String,
    },
}

#[derive(Debug, Serialize, Deserialize)]
pub enum BoundaryValue {
    I32(i32),
    I64(i64),
    F32(u32),
    F64(u64),
}

pub fn runtime_value_to_boundary_value(rv: RuntimeValue) -> BoundaryValue {
    match rv {
        RuntimeValue::I32(rv) => BoundaryValue::I32(rv),
        RuntimeValue::I64(rv) => BoundaryValue::I64(rv),
        RuntimeValue::F32(rv) => BoundaryValue::F32(rv.to_bits()),
        RuntimeValue::F64(rv) => BoundaryValue::F64(rv.to_bits()),
    }
}

pub fn boundary_value_to_runtime_value(rv: BoundaryValue) -> RuntimeValue {
    match rv {
        BoundaryValue::I32(bv) => RuntimeValue::I32(bv),
        BoundaryValue::I64(bv) => RuntimeValue::I64(bv),
        BoundaryValue::F32(bv) => RuntimeValue::F32(f32::from_bits(bv)),
        BoundaryValue::F64(bv) => RuntimeValue::F64(f64::from_bits(bv)),
    }
}

pub fn result_covert(res : Result<Option<RuntimeValue>, InterpreterError>)
                     -> Result<Option<BoundaryValue>, InterpreterError>
{
    match res {
        Ok(None) => Ok(None),
        Ok(Some(rv)) => Ok(Some(runtime_value_to_boundary_value(rv))),
        Err(x) => Err(x),
    }
}

pub struct SpecModule {
    table: TableRef,
    memory: MemoryRef,
    global_i32: GlobalRef,
    global_f32: GlobalRef,
    global_f64: GlobalRef,
}

impl SpecModule {
    pub fn new() -> Self {
        SpecModule {
            table: TableInstance::alloc(10, Some(20)).unwrap(),
            memory: MemoryInstance::alloc(Pages(1), Some(Pages(2))).unwrap(),
            global_i32: GlobalInstance::alloc(RuntimeValue::I32(666), false),
            global_f32: GlobalInstance::alloc(RuntimeValue::F32(666.0), false),
            global_f64: GlobalInstance::alloc(RuntimeValue::F64(666.0), false),
        }
    }
}

pub fn spec_to_runtime_value(value: Value) -> RuntimeValue {
    match value {
        Value::I32(v) => RuntimeValue::I32(v),
        Value::I64(v) => RuntimeValue::I64(v),
        Value::F32(v) => RuntimeValue::F32(v),
        Value::F64(v) => RuntimeValue::F64(v),
    }
}

#[derive(Debug)]
pub enum Error {
    Load(String),
    Start(Trap),
    Script(script::Error),
    Interpreter(InterpreterError),
}

impl From<InterpreterError> for Error {
    fn from(e: InterpreterError) -> Error {
        Error::Interpreter(e)
    }
}

impl From<script::Error> for Error {
    fn from(e: script::Error) -> Error {
        Error::Script(e)
    }
}

const PRINT_FUNC_INDEX: usize = 0;

impl Externals for SpecModule {
    fn invoke_index(
        &mut self,
        index: usize,
        args: RuntimeArgs,
    ) -> Result<Option<RuntimeValue>, Trap> {
        match index {
            PRINT_FUNC_INDEX => {
                println!("print: {:?}", args);
                Ok(None)
            }
            _ => panic!("SpecModule doesn't provide function at index {}", index),
        }
    }
}

impl ModuleImportResolver for SpecModule {
    fn resolve_func(
        &self,
        field_name: &str,
        func_type: &Signature,
    ) -> Result<FuncRef, InterpreterError> {
        let index = match field_name {
            "print" => PRINT_FUNC_INDEX,
            "print_i32" => PRINT_FUNC_INDEX,
            "print_i32_f32" => PRINT_FUNC_INDEX,
            "print_f64_f64" => PRINT_FUNC_INDEX,
            "print_f32" => PRINT_FUNC_INDEX,
            "print_f64" => PRINT_FUNC_INDEX,
            _ => {
                return Err(InterpreterError::Instantiation(format!(
                    "Unknown host func import {}",
                    field_name
                )));
            }
        };

        if func_type.return_type().is_some() {
            return Err(InterpreterError::Instantiation(
                "Function `print_` have unit return type".into(),
            ));
        }

        let func = FuncInstance::alloc_host(func_type.clone(), index);
        return Ok(func);
    }
    fn resolve_global(
        &self,
        field_name: &str,
        _global_type: &GlobalDescriptor,
    ) -> Result<GlobalRef, InterpreterError> {
        match field_name {
            "global_i32" => Ok(self.global_i32.clone()),
            "global_f32" => Ok(self.global_f32.clone()),
            "global_f64" => Ok(self.global_f64.clone()),
            _ => Err(InterpreterError::Instantiation(format!(
                "Unknown host global import {}",
                field_name
            )))
        }
    }

    fn resolve_memory(
        &self,
        field_name: &str,
        _memory_type: &MemoryDescriptor,
    ) -> Result<MemoryRef, InterpreterError> {
        if field_name == "memory" {
            return Ok(self.memory.clone());
        }

        Err(InterpreterError::Instantiation(format!(
            "Unknown host memory import {}",
            field_name
        )))
    }

    fn resolve_table(
        &self,
        field_name: &str,
        _table_type: &TableDescriptor,
    ) -> Result<TableRef, InterpreterError> {
        if field_name == "table" {
            return Ok(self.table.clone());
        }

        Err(InterpreterError::Instantiation(format!(
            "Unknown host table import {}",
            field_name
        )))
    }
}

pub struct SpecDriver {
    spec_module: SpecModule,
    instances: HashMap<String, ModuleRef>,
    last_module: Option<ModuleRef>,
}

impl SpecDriver {
    pub fn new() -> SpecDriver {
        SpecDriver {
            spec_module: SpecModule::new(),
            instances: HashMap::new(),
            last_module: None,
        }
    }

    pub fn spec_module(&mut self) -> &mut SpecModule {
        &mut self.spec_module
    }

    pub fn add_module(&mut self, name: Option<String>, module: ModuleRef) {
        self.last_module = Some(module.clone());
        if let Some(name) = name {
            self.instances.insert(name, module);
        }
    }

    pub fn module(&self, name: &str) -> Result<ModuleRef, InterpreterError> {
        self.instances.get(name).cloned().ok_or_else(|| {
            InterpreterError::Instantiation(format!("Module not registered {}", name))
        })
    }

    pub fn module_or_last(&self, name: Option<&str>) -> Result<ModuleRef, InterpreterError> {
        match name {
            Some(name) => self.module(name),
            None => self.last_module
                .clone()
                .ok_or_else(|| InterpreterError::Instantiation("No modules registered".into())),
        }
    }

    pub fn register(&mut self, name : &Option<String>,
                    as_name : String) -> Result<(), InterpreterError> {
        let module = match self.module_or_last(name.as_ref().map(|x| x.as_ref())) {
            Ok(module) => module,
            Err(_) => return Err(InterpreterError::Instantiation("No such modules registered".into())),
        };
        self.add_module(Some(as_name), module);
        Ok(())
    }
}

impl ImportResolver for SpecDriver {
    fn resolve_func(
        &self,
        module_name: &str,
        field_name: &str,
        func_type: &Signature,
    ) -> Result<FuncRef, InterpreterError> {
        if module_name == "spectest" {
            self.spec_module.resolve_func(field_name, func_type)
        } else {
            self.module(module_name)?
                .resolve_func(field_name, func_type)
        }
    }

    fn resolve_global(
        &self,
        module_name: &str,
        field_name: &str,
        global_type: &GlobalDescriptor,
    ) -> Result<GlobalRef, InterpreterError> {
        if module_name == "spectest" {
            self.spec_module.resolve_global(field_name, global_type)
        } else {
            self.module(module_name)?
                .resolve_global(field_name, global_type)
        }
    }

    fn resolve_memory(
        &self,
        module_name: &str,
        field_name: &str,
        memory_type: &MemoryDescriptor,
    ) -> Result<MemoryRef, InterpreterError> {
        if module_name == "spectest" {
            self.spec_module.resolve_memory(field_name, memory_type)
        } else {
            self.module(module_name)?
                .resolve_memory(field_name, memory_type)
        }
    }

    fn resolve_table(
        &self,
        module_name: &str,
        field_name: &str,
        table_type: &TableDescriptor,
    ) -> Result<TableRef, InterpreterError> {
        if module_name == "spectest" {
            self.spec_module.resolve_table(field_name, table_type)
        } else {
            self.module(module_name)?
                .resolve_table(field_name, table_type)
        }
    }
}

pub fn try_load_module(wasm: &[u8]) -> Result<Module, Error> {
    Module::from_buffer(wasm).map_err(|e| Error::Load(e.to_string()))
}

pub fn try_load(wasm: &[u8], spec_driver: &mut SpecDriver) -> Result<(), Error> {
    let module = try_load_module(wasm)?;
    let instance = ModuleInstance::new(&module, &ImportsBuilder::default())?;
    instance
        .run_start(spec_driver.spec_module())
        .map_err(|trap| Error::Start(trap))?;
    Ok(())
}

pub fn load_module(wasm: &[u8], name: &Option<String>, spec_driver: &mut SpecDriver) -> Result<ModuleRef, Error> {
    let module = try_load_module(wasm)?;
    let instance = ModuleInstance::new(&module, spec_driver)
        .map_err(|e| Error::Load(e.to_string()))?
        .run_start(spec_driver.spec_module())
        .map_err(|trap| Error::Start(trap))?;

    let module_name = name.clone();
    spec_driver.add_module(module_name, instance.clone());

    Ok(instance)
}

