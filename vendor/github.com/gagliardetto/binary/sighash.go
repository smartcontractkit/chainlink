// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bin

import (
	"crypto/sha256"
)

// Sighash creates an anchor sighash for the provided namespace and element.
// An anchor sighash is the first 8 bytes of the sha256 of {namespace}:{name}
// NOTE: you must first convert the name to snake case using `ToSnakeForSighash`.
func Sighash(namespace string, name string) []byte {
	data := namespace + ":" + name
	sum := sha256.Sum256([]byte(data))
	return sum[0:8]
}

func SighashInstruction(name string) []byte {
	// Instruction sighash are the first 8 bytes of the sha256 of
	// {SIGHASH_INSTRUCTION_NAMESPACE}:{snake_case(name)}
	return Sighash(SIGHASH_GLOBAL_NAMESPACE, ToSnakeForSighash(name))
}

func SighashAccount(name string) []byte {
	// Account sighash are the first 8 bytes of the sha256 of
	// {SIGHASH_ACCOUNT_NAMESPACE}:{camelCase(name)}
	return Sighash(SIGHASH_ACCOUNT_NAMESPACE, ToPascalCase(name))
}

// NOTE: no casing conversion is done here, it's up to the caller to
// provide the correct casing.
func SighashTypeID(namespace string, name string) TypeID {
	return TypeIDFromBytes(Sighash(namespace, (name)))
}

// Namespace for calculating state instruction sighash signatures.
const SIGHASH_STATE_NAMESPACE string = "state"

// Namespace for calculating instruction sighash signatures for any instruction
// not affecting program state.
const SIGHASH_GLOBAL_NAMESPACE string = "global"

const SIGHASH_ACCOUNT_NAMESPACE string = "account"

const ACCOUNT_DISCRIMINATOR_SIZE = 8

// https://github.com/project-serum/anchor/pull/64/files
// https://github.com/project-serum/anchor/blob/2f780e0d274f47e442b3f0d107db805a41c6def0/ts/src/coder/common.ts#L109
// https://github.com/project-serum/anchor/blob/6b5ed789fc856408986e8868229887354d6d4073/lang/syn/src/codegen/program/common.rs#L17

// TODO:
// https://github.com/project-serum/anchor/blob/84a2b8200cc3c7cb51d7127918e6cbbd836f0e99/ts/src/error.ts#L48
